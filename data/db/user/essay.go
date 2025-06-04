package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"writescore/global"
	"writescore/models/dao"
	"writescore/models/dto"
	"writescore/utils"

	"github.com/gin-gonic/gin"
)

func RantingEssay(c *gin.Context, param dto.RatingEssayMap, userId int64) (data dto.RatingResult, err error) {
	tx := global.GetDbConn(c).Begin()
	if tx.Error != nil {
		return dto.RatingResult{}, fmt.Errorf("开启事务失败：%w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("事务回滚: %v", r)
		}
	}()

	var content string
	err = tx.Model(&dao.Essay{}).Where("id = ? and user_id = ?", param.EssayId, userId).Select("content").First(&content).Error
	if err != nil {
		tx.Rollback()
		return dto.RatingResult{}, fmt.Errorf("获取文章内容失败:%w", err)
	}

	//先获取评分标准包括哪些
	var creteria []dao.ScoringCriteria
	if err = tx.Model(&dao.ScoringCriteria{}).Find(&creteria).Error; err != nil {
		tx.Rollback()
		return dto.RatingResult{}, fmt.Errorf("获取评分标准失败:%w", err)
	}

	if len(creteria) == 0 {
		tx.Rollback()
		return dto.RatingResult{}, fmt.Errorf("未找到评分标准")
	}

	//把标题存进essay表
	if param.Title != "" {
		if err = tx.Model(&dao.Essay{}).Where("id = ?", param.EssayId).Updates(map[string]interface{}{
			"title": param.Title,
		}).Error; err != nil {
			tx.Rollback()
			return dto.RatingResult{}, fmt.Errorf("保存用户标题失败:%w", err)
		}
	}

	var originScore sql.NullFloat64 //接收可能为null的值
	if err = tx.Model(&dao.Essay{}).Where("id = ?", param.EssayId).Select("score").First(&originScore).Error; err != nil {
		tx.Rollback()
		return dto.RatingResult{}, err
	}

	perScores := make([]dto.PerScore, 0, len(creteria)) // 预分配容量
	for _, key := range creteria {
		score, feekback, err := evaluateCriterion(content, key.CriteriaName) //依次获取各项指标的结果
		if err != nil {
			tx.Rollback() //失败回滚
			return dto.RatingResult{}, fmt.Errorf("获取各项指标的评分以及反馈结果失败:%w", err)
		}

		scoreDetail := dao.EssayScoringDetails{
			EssayID:    param.EssayId,
			CriteriaID: key.ID,
			Score:      score,
			Feedback:   feekback,
			CreateTime: time.Now(),
		}

		if originScore.Valid { //说明之前是有评分过的
			if err = tx.Model(&dao.EssayScoringDetails{}).Where("essay_id = ? and criteria_id = ?", param.EssayId, scoreDetail.CriteriaID).Updates(map[string]interface{}{
				"score":       scoreDetail.Score,
				"feedback":    scoreDetail.Feedback,
				"create_time": scoreDetail.CreateTime,
			}).Error; err != nil {
				tx.Rollback()
				return dto.RatingResult{}, fmt.Errorf("更新评分信息失败:%w", err)
			}
		} else {
			//将评分结果以及反馈结果保存到数据库
			if err = tx.Model(&dao.EssayScoringDetails{}).Create(&scoreDetail).Error; err != nil {
				tx.Rollback()
				return dto.RatingResult{}, fmt.Errorf("将各项指标保存到数据库失败：%w", err)
			}
		}

		perScores = append(perScores, dto.PerScore{
			CriteriaName:  key.CriteriaName,
			CriteriaScore: score,
			Feekback:      feekback,
			CriteriaId:    key.ID,
		})
	}

	if len(perScores) == 0 {
		tx.Rollback()
		return dto.RatingResult{}, fmt.Errorf("评分结果为空")
	}

	log.Println("perScores:")
	for _, perScore := range perScores {
		log.Printf("CriteriaName: %s, CriteriaScore: %.2f, Feedback: %s, CriteriaId: %d",
			perScore.CriteriaName, perScore.CriteriaScore, perScore.Feekback, perScore.CriteriaId)
	}

	finalScore := calculateFinalScore(perScores, creteria)
	finalFeekback, err := calculateFinalFeekback(content, perScores)
	if err != nil {
		tx.Rollback()
		return dto.RatingResult{}, fmt.Errorf("获取最终结果失败:%w", err)
	}

	if err := tx.Model(&dao.Essay{}).
		Where("id = ?", param.EssayId).
		Updates(map[string]interface{}{
			"score":    finalScore,
			"feedback": finalFeekback,
		}).Error; err != nil {
		tx.Rollback()
		return dto.RatingResult{}, fmt.Errorf("将最终评分以及反馈保存到数据库失败:%w", err)
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return dto.RatingResult{}, fmt.Errorf("提交事务失败:%w", err)
	}

	return dto.RatingResult{
		EssayId:       param.EssayId,
		PerScore:      perScores,
		FinalFeekback: finalFeekback,
		FinalScore:    finalScore,
	}, nil
}

// 获得某指标的得分和反馈结果(调用deepseek)

func evaluateCriterion(content string, criteriaName string) (score float64, feedback string, err error) {
	prompt := fmt.Sprintf(`
请严格以JSON格式返回，包含"score"和"feedback"字段，不要包含其他内容，不要使用markdown格式或代码块，直接返回JSON对象。
feedback字段只返回详细的反馈信息，"score"字段只返回数字评分
评分标准:"%s"
作文内容:%s

请提供:
1. 该标准的得分（0-100分）
2. 详细的反馈，解释与该标准相关的优点和不足
`, criteriaName, content)

	response, err := callDeepseekAPI(prompt)
	if err != nil {
		return 0, "", fmt.Errorf("call deepseek失败:%w", err)
	}

	log.Printf("API Response: %s", response)

	// 清理返回的内容，移除 Markdown 代码块标记
	response = strings.TrimSpace(response)

	// 移除开头的 ```json 标记
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
	}

	// 移除结尾的 ``` 标记
	if strings.HasSuffix(response, "```") {
		response = strings.TrimSuffix(response, "```")
	}

	// 移除其他可能的标记
	response = strings.Trim(response, "`")
	response = strings.Trim(response, "\"")
	response = strings.TrimSpace(response)

	var result struct {
		Score    float64 `json:"score"`
		Feedback string  `json:"feedback"`
	}

	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		log.Printf("JSON解析失败: %v", err)
		// Fallback to text parsing if JSON parsing fails
		scoreStr := extractScore(response)
		scoreStr = strings.TrimSpace(scoreStr) // 去掉多余空格
		score, err := parseScore(scoreStr)
		if err != nil {
			log.Printf("文本解析失败: %v", err)
			return 70, extractFeedback(response), nil // Default score if parsing fails
		}
		return score, extractFeedback(response), nil
	}

	log.Printf("解析成功: score=%v, feedback=%v", result.Score, result.Feedback)
	return result.Score, result.Feedback, nil

}

func extractScore(response string) string {
	start := strings.Index(response, `"score":`)
	if start == -1 {
		return ""
	}
	end := strings.Index(response[start:], ",")
	if end == -1 {
		end = len(response)
	} else {
		end += start
	}
	return response[start+8 : end]
}

func parseScore(scoreStr string) (float64, error) {
	scoreStr = strings.TrimSpace(scoreStr) // 去掉多余空格
	score, err := strconv.ParseFloat(scoreStr, 64)
	if err != nil {
		return 0, fmt.Errorf("分数解析失败: %w", err)
	}
	return score, nil
}
func extractFeedback(response string) string {
	start := strings.Index(response, `"feedback":`)
	if start == -1 {
		return ""
	}
	return response[start+12:]
}

//func extractFeedback(text string) string {
//	// Try to extract feedback from text
//	feedbackRegex := regexp.MustCompile(`(?i)feedback\s*[:：]\s*(.+)`)
//	matches := feedbackRegex.FindStringSubmatch(text)
//	if len(matches) > 1 {
//		return matches[1]
//	}
//	return text // Return the whole text if no specific feedback section found
//}

// 获得最终评分，根据各项指标的得分以及所占的比重
func calculateFinalScore(perscore []dto.PerScore, creteria []dao.ScoringCriteria) (finalScore float64) {
	for _, score := range perscore {
		for _, criterion := range creteria {
			if score.CriteriaId == criterion.ID {
				finalScore += score.CriteriaScore * criterion.Weight
			}
		}
	}
	return
}

// 算出最终反馈结果(调用deepseek)
func calculateFinalFeekback(content string, perScores []dto.PerScore) (finalFeekback string, err error) {
	var criteriaFeedback strings.Builder
	for _, score := range perScores {
		criteriaFeedback.WriteString(fmt.Sprintf("%s: %s\n\n", score.CriteriaName, score.Feekback))
	}

	// Prepare the prompt for Deepseek
	prompt := fmt.Sprintf(`
		根据以下作文和详细评估，提供一个全面、建设性的反馈总结。

		作文:
		%s

		详细评估:
		%s

		请提供一个结构良好、鼓励性的反馈总结，突出优点、需要改进的地方和具体建议。
	`, truncateContent(content, 1000), criteriaFeedback.String())

	// Call Deepseek API
	response, err := callDeepseekAPI(prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}
func truncateContent(content string, maxLength int) string {
	if len(content) <= maxLength {
		return content
	}
	return content[:maxLength] + "..."
}

// 调用deepseek的api

func callDeepseekAPI(prompt string) (string, error) {
	url := "https://api.deepseek.com/v1/chat/completions" // API地址
	payload := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %w", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer sk-dd54e9f4bf6c43f9a103f39965b1e008") // 替换为实际的授权令牌

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("API返回空响应")
	}

	return result.Choices[0].Message.Content, nil
}

//func callDeepseekAPI(prompt string) (string, error) {
//	url := "https://api.deepseek.com/v1/chat/completions" //api地址，用来发送请求获取回答
//	payload := map[string]interface{}{                    //构建请求负载
//		"model": "deepseek-chat", //指定使用的模型,需不需要换？
//		"messages": []map[string]string{ //包含用户的信息，是一个数组
//			{
//				"role":    "user", //表示消息的发送者是用户
//				"content": prompt, //用户输入的prompt
//			},
//		},
//		"temperature": 0.7, //用于控制生成的回答的随机性，值越高获取的回答越随机，值越低生成的回答越确定
//		//0.7表示在随机性和确定性之间取得平衡
//	}
//
//	jsonPayload, err := json.Marshal(payload) //将请求负载转换为json
//	if err != nil {
//		return "", fmt.Errorf("JSON编码失败: %w", err)
//	}
//
//	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
//	//表示创建一个http post请求，strings.newReader()//将json数据作为请求的正文
//	if err != nil {
//		return "", fmt.Errorf("创建请求失败: %w", err)
//	}
//
//	req.Header.Add("Content-Type", "application/json") //告诉服务器发送的数据类型
//
//	req.Header.Add("Authorization", "Bearer sk-dd54e9f4bf6c43f9a103f39965b1e008") // 授权头，用于验证身份
//	//这里需要替换
//
//	client := &http.Client{}
//	resp, err := client.Do(req) //使用client.Do发送请求并获取响应
//	if err != nil {
//		return "", fmt.Errorf("API请求失败: %w", err)
//	}
//	defer resp.Body.Close() //确保在函数返回时关闭响应体，避免资源泄露
//
//	body, err := ioutil.ReadAll(resp.Body) //读取响应体的内容
//	if err != nil {
//		return "", fmt.Errorf("读取响应失败: %w", err)
//	}
//	log.Printf("API Response: %s", body) //打印api返回的原始内容
//	//存储解析后的响应数据
//	var result struct {
//		Choices []struct {
//			Message struct {
//				Content string `json:"content"`
//			} `json:"message"`
//		} `json:"choices"`
//	}
//	//将响应体的json数据解析到result中
//	if err := json.Unmarshal(body, &result); err != nil {
//		return "", fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if len(result.Choices) == 0 {
//		return "", fmt.Errorf("API返回空响应")
//	}
//
//	return result.Choices[0].Message.Content, nil
//}

func GetEssay(c *gin.Context, userId int64, param dto.GetEssayMap) (data []dto.AllEssays, total int64, err error) {
	db := global.GetDbConn(c).Model(&dao.Essay{}).Where("user_id = ?", userId)
	if param.Title != "" {
		db = db.Where("title like ?", "%"+param.Title+"%")
	}

	// 处理开始时间
	if param.StartTime != "" {
		// 如果只包含日期，自动添加开始时间 00:00:00
		if len(param.StartTime) == 10 { // 格式为 YYYY-MM-DD
			param.StartTime = param.StartTime + " 00:00:00"
		}
		startTime, err := utils.StringToTime(param.StartTime)
		if err != nil {
			return make([]dto.AllEssays, 0), 0, err
		}
		db = db.Where("submit_time >= ?", startTime)
	} else { // 默认从最早时间开始
		db = db.Where("submit_time >= ?", "1999-09-09 00:00:00")
	}

	// 处理结束时间
	if param.EndTime != "" {
		// 如果只包含日期，自动添加结束时间 23:59:59
		if len(param.EndTime) == 10 { // 格式为 YYYY-MM-DD
			param.EndTime = param.EndTime + " 23:59:59"
		}
		endTime, err := utils.StringToTime(param.EndTime)
		if err != nil {
			return make([]dto.AllEssays, 0), 0, err
		}
		db = db.Where("submit_time <= ?", endTime)
	} else { // 默认到最晚时间
		db = db.Where("submit_time <= ?", "2999-09-09 23:59:59")
	}

	if param.MinScore != 0 {
		db = db.Where("score >= ?", param.MinScore)
	} else {
		db = db.Where("score >= ? ", 0)
	}

	if param.MaxScore != 0 { //没传默认是100
		db = db.Where("score <= ?", param.MaxScore)
	} else {
		db = db.Where("score <= ? ", 100)
	}
	if err = db.Select("submit_time,id,score,title").Count(&total).Order("submit_time desc").
		Offset((param.PageIndex - 1) * param.PageSize).Limit(param.PageSize).
		Find(&data).Error; err != nil {
		return make([]dto.AllEssays, 0), 0, err
	}
	for i := range data {
		data[i].SubmitTimeMar = utils.MarshalTime(data[i].SubmitTime)
	}
	return data, total, nil
}

func GetEssayDetails(c *gin.Context, userId int64, id int) (data dto.EssayDetail, err error) {
	var essay dao.Essay
	var perScore []dto.PerScore
	err = global.GetDbConn(c).Model(&dao.Essay{}).Where("essay.user_id = ? and essay.id = ?", userId, id).
		//Joins("joins essay_scoring_details on essay.id = essay_scoring_details.essay_id").
		Select("essay.content,essay.submit_time,essay.score,essay.feedback,essay.title").
		Find(&essay).Error
	if err != nil {
		return dto.EssayDetail{}, err
	}
	err = global.GetDbConn(c).Model(&dao.EssayScoringDetails{}).Where("essay_scoring_details.essay_id = ?", id).
		Joins("join scoring_criteria sc on essay_scoring_details.criteria_id = sc.id").
		Select("sc.criteria_name,sc.id as criteria_id,essay_scoring_details.score as criteria_score,essay_scoring_details.feedback").
		Find(&perScore).Error

	if err != nil {
		return dto.EssayDetail{}, fmt.Errorf("获取评分详情失败: %w", err)
	}
	data = dto.EssayDetail{
		PerScore: perScore,
		Content:  essay.Content,
		//SubmitTime:    essay.SubmitTime, //可以删除
		SubmitTimeMar: utils.MarshalTime(essay.SubmitTime),
		Score:         essay.Score,
		Feedback:      essay.Feedback,
		Title:         essay.Title,
	}
	return data, nil
}

//https://platform.deepseek.com/usage  deepseek开放平台

// StreamRatingEssay 流式评分接口
func StreamRatingEssay(c *gin.Context, param dto.RatingEssayMap, userId int64) error {
	// 设置SSE头部
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	tx := global.GetDbConn(c).Begin()
	if tx.Error != nil {
		return fmt.Errorf("开启事务失败：%w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var content string
	err := tx.Model(&dao.Essay{}).Where("id = ? and user_id = ?", param.EssayId, userId).Select("content").First(&content).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("获取文章内容失败:%w", err)
	}

	// 获取评分标准
	var creteria []dao.ScoringCriteria
	if err = tx.Model(&dao.ScoringCriteria{}).Find(&creteria).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("获取评分标准失败:%w", err)
	}

	if len(creteria) == 0 {
		tx.Rollback()
		return fmt.Errorf("未找到评分标准")
	}

	// 保存标题
	if param.Title != "" {
		if err = tx.Model(&dao.Essay{}).Where("id = ?", param.EssayId).Updates(map[string]interface{}{
			"title": param.Title,
		}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("保存用户标题失败:%w", err)
		}
	}

	var originScore sql.NullFloat64
	if err = tx.Model(&dao.Essay{}).Where("id = ?", param.EssayId).Select("score").First(&originScore).Error; err != nil {
		tx.Rollback()
		return err
	}

	perScores := make([]dto.PerScore, 0, len(creteria))

	// 发送开始评分的消息
	c.SSEvent("message", map[string]interface{}{
		"type": "start",
		"data": "开始评分...",
	})
	c.Writer.Flush()

	for _, key := range creteria {
		// 发送正在评分的标准
		c.SSEvent("message", map[string]interface{}{
			"type": "progress",
			"data": fmt.Sprintf("正在评分: %s", key.CriteriaName),
		})
		c.Writer.Flush()

		score, feedback, err := evaluateCriterion(content, key.CriteriaName)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("获取各项指标的评分以及反馈结果失败:%w", err)
		}

		scoreDetail := dao.EssayScoringDetails{
			EssayID:    param.EssayId,
			CriteriaID: key.ID,
			Score:      score,
			Feedback:   feedback,
			CreateTime: time.Now(),
		}

		if originScore.Valid {
			if err = tx.Model(&dao.EssayScoringDetails{}).Where("essay_id = ? and criteria_id = ?", param.EssayId, scoreDetail.CriteriaID).Updates(map[string]interface{}{
				"score":       scoreDetail.Score,
				"feedback":    scoreDetail.Feedback,
				"create_time": scoreDetail.CreateTime,
			}).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("更新评分信息失败:%w", err)
			}
		} else {
			if err = tx.Model(&dao.EssayScoringDetails{}).Create(&scoreDetail).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("将各项指标保存到数据库失败：%w", err)
			}
		}

		perScore := dto.PerScore{
			CriteriaName:  key.CriteriaName,
			CriteriaScore: score,
			Feekback:      feedback,
			CriteriaId:    key.ID,
		}
		perScores = append(perScores, perScore)

		// 发送单个标准的评分结果
		c.SSEvent("message", map[string]interface{}{
			"type": "criterion_result",
			"data": perScore,
		})
		c.Writer.Flush()
	}

	if len(perScores) == 0 {
		tx.Rollback()
		return fmt.Errorf("评分结果为空")
	}

	// 计算最终分数
	finalScore := calculateFinalScore(perScores, creteria)
	finalFeedback, err := calculateFinalFeekback(content, perScores)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("获取最终结果失败:%w", err)
	}

	// 更新最终分数和反馈
	if err := tx.Model(&dao.Essay{}).
		Where("id = ?", param.EssayId).
		Updates(map[string]interface{}{
			"score":    finalScore,
			"feedback": finalFeedback,
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("将最终评分以及反馈保存到数据库失败:%w", err)
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("提交事务失败:%w", err)
	}

	// 发送最终结果
	c.SSEvent("message", map[string]interface{}{
		"type": "final_result",
		"data": dto.RatingResult{
			EssayId:       param.EssayId,
			PerScore:      perScores,
			FinalFeekback: finalFeedback,
			FinalScore:    finalScore,
		},
	})
	c.Writer.Flush()

	// 发送完成消息
	c.SSEvent("message", map[string]interface{}{
		"type": "complete",
		"data": "评分完成",
	})
	c.Writer.Flush()

	return nil
}
