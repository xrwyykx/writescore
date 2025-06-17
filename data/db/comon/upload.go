package comon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"writescore/global"
	"writescore/models/dao"
	"writescore/models/dto"
	"writescore/utils"

	"github.com/gin-gonic/gin"
)

// GetAccessToken 获取新的 access_token
//func GetAccessToken(apiKey, secretKey string) (string, error) {
//	tokenURL := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", apiKey, secretKey)
//	resp, err := http.Get(tokenURL)
//	if err != nil {
//		return "", fmt.Errorf("failed to get access token: %v", err)
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return "", fmt.Errorf("failed to read response body: %v", err)
//	}
//
//	var result map[string]interface{}
//	err = json.Unmarshal(body, &result)
//	if err != nil {
//		return "", fmt.Errorf("failed to unmarshal response: %v", err)
//	}
//
//	accessToken, ok := result["access_token"].(string)
//	if !ok {
//		errorMsg, ok := result["error_msg"].(string)
//		if ok {
//			return "", fmt.Errorf("failed to get access token: %s", errorMsg)
//		}
//		return "", fmt.Errorf("failed to get access token from response")
//	}
//
//	return accessToken, nil
//}
//
//// OcrHandwritingWithBaidu 使用百度OCR识别手写文字
//func OcrHandwritingWithBaidu(imgURL string) (string, error) {
//	apiKey := "aLnqDRQEwLa7yX7XhJJKa1X7"
//	secretKey := "6IkAaJhaZAGPKmaJADy2EPAdkkUauyVD"
//
//	// 1. 获取Access Token
//	accessToken, err := GetAccessToken(apiKey, secretKey)
//	if err != nil {
//		return "", fmt.Errorf("failed to get access token: %v", err)
//	}
//
//	// 2. 调用手写OCR接口
//	ocrURL := fmt.Sprintf("https://aip.baidubce.com/rest/2.0/ocr/v1/handwriting?access_token=%s", accessToken)
//	formData := url.Values{"url": {imgURL}, "recognize_granularity": {"big"}} // 大颗粒度适合手写
//
//	resp, err := http.Post(ocrURL, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
//	if err != nil {
//		return "", fmt.Errorf("failed to call OCR API: %v", err)
//	}
//	defer resp.Body.Close()
//
//	// 3. 解析结果（百度返回JSON结构较复杂，需多层解析）
//	ocrBody, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return "", fmt.Errorf("failed to read OCR response body: %v", err)
//	}
//
//	var ocrResult map[string]interface{}
//	err = json.Unmarshal(ocrBody, &ocrResult)
//	if err != nil {
//		return "", fmt.Errorf("failed to unmarshal OCR response: %v", err)
//	}
//
//	// 检查是否包含错误信息
//	if _, ok := ocrResult["error_code"]; ok {
//		errorMsg, ok := ocrResult["error_msg"].(string)
//		if ok {
//			return "", fmt.Errorf("OCR failed: %s", errorMsg)
//		}
//		return "", fmt.Errorf("OCR failed with unknown error")
//	}
//
//	wordsResult, ok := ocrResult["words_result"].([]interface{})
//	if !ok {
//		return "", fmt.Errorf("words_result not found in OCR response")
//	}
//
//	var text strings.Builder
//	for _, word := range wordsResult {
//		wordMap, ok := word.(map[string]interface{})
//		if !ok {
//			return "", fmt.Errorf("invalid word format in OCR response")
//		}
//		text.WriteString(wordMap["words"].(string) + "\n")
//	}
//
//	return text.String(), nil
//}
//
//func RestoreImageInfo(c *gin.Context, userId int64, param dto.RestoreImageInfoMap) (data dto.ImageToEssay, err error) {
//	//第一个事务是将照片信息保存在照片表，然后接着根据前面的照片信息将照片解析成文字，不是同时进行的
//	//如果前面的事务没有完成后面的事务也不能顺利进行，应该不需要启动事务？
//	//照片信息不需要关联用户id，如果解析照片文字成功就可以把essay_id更新到image_info表，这个过程中需不需要开启事务
//	//1.保存图片信息
//	nowTime := time.Now()
//	imageInfo := dao.ImageInfo{
//		ImageURL:   param.ImageURL,
//		ImageName:  param.ImageName,
//		UploadTime: nowTime,
//	}
//	if err := global.GetDbConn(c).Model(&dao.ImageInfo{}).Create(&imageInfo).Error; err != nil {
//		return dto.ImageToEssay{}, err
//	} //保存图片信息失败时
//	//2.解析图片文字,并将这篇文章的信息更新到数据库表
//	essay, err := OcrHandwritingWithBaidu(param.ImageURL)
//	if err != nil {
//		return dto.ImageToEssay{}, err
//	}
//	essayInfo := dao.Essay{
//		UserID:     userId,
//		Content:    essay,
//		SubmitTime: nowTime,
//		Title:      param.Title,
//		//WordCount: ,//计算字符串长度
//		//UploadMethod: 0,            //图片上传
//		ImageId: imageInfo.ID, //实现与image_info表的关联
//	}
//	if err := global.GetDbConn(c).Model(&dao.Essay{}).Create(&essayInfo).Error; err != nil {
//		return dto.ImageToEssay{}, err
//	}
//
//	//返回信息
//	data.EssayId = essayInfo.ID
//	data.SubmitTime = nowTime
//	data.SubmitTimeMar = utils.MarshalTime(data.SubmitTime)
//	data.Content = essay
//	return data, nil
//}

// GetAccessToken 获取新的 access_token
func GetAccessToken(apiKey, secretKey string) (string, error) {
	tokenURL := fmt.Sprintf(
		"https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s",
		apiKey, secretKey)

	resp, err := http.Get(tokenURL)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		errorMsg, ok := result["error_msg"].(string)
		if ok {
			return "", fmt.Errorf("failed to get access token: %s", errorMsg)
		}
		return "", fmt.Errorf("failed to get access token from response")
	}

	return accessToken, nil
}

// OcrHandwritingWithBaidu 使用百度OCR识别手写文字 (修复双重编码问题)
func OcrHandwritingWithBaidu(imgURL string) (string, error) {
	apiKey := "aLnqDRQEwLa7yX7XhJJKa1X7"
	secretKey := "6IkAaJhaZAGPKmaJADy2EPAdkkUauyVD"

	// 关键修复：移除协议头校验（百度API会自动处理编码）
	// 因为表单编码会处理特殊字符，我们不需要额外校验

	// 1. 获取Access Token
	accessToken, err := GetAccessToken(apiKey, secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %v", err)
	}

	ocrURL := fmt.Sprintf(
		"https://aip.baidubce.com/rest/2.0/ocr/v1/handwriting?access_token=%s",
		accessToken)

	formData := url.Values{
		"url":                   {imgURL}, // 使用原始URL
		"recognize_granularity": {"big"},
	}

	resp, err := http.Post(ocrURL, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to call OCR API: %v", err)
	}
	defer resp.Body.Close()

	ocrBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read OCR response body: %v", err)
	}

	var ocrResult map[string]interface{}
	err = json.Unmarshal(ocrBody, &ocrResult)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal OCR response: %v", err)
	}

	// 增强错误处理
	if errorCode, ok := ocrResult["error_code"]; ok {
		errorMsg, _ := ocrResult["error_msg"].(string)
		fmt.Errorf("百度OCR错误: code=%v msg=%s 请求URL=%s",
			errorCode, errorMsg, imgURL)
		return "", fmt.Errorf("OCR失败: %s", errorMsg)
	}

	wordsResult, ok := ocrResult["words_result"].([]interface{})
	if !ok {
		return "", fmt.Errorf("OCR响应中未找到words_result")
	}

	var text strings.Builder
	for _, word := range wordsResult {
		wordMap, ok := word.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("OCR响应中的词条格式无效")
		}
		text.WriteString(wordMap["words"].(string) + "\n")
	}

	return text.String(), nil
}

// RestoreImageInfo 恢复图片信息（修复双重编码问题）
func RestoreImageInfo(c *gin.Context, userId int64, param dto.RestoreImageInfoMap) (data dto.ImageToEssay, err error) {
	// 开启事务
	tx := global.GetDbConn(c).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("事务回滚: %v", r)
		}
	}()

	nowTime := time.Now()

	// 1. 保存图片信息（事务内）
	imageInfo := dao.ImageInfo{
		ImageURL:   param.ImageURL,
		ImageName:  param.ImageName,
		UploadTime: nowTime,
	}
	if err = tx.Create(&imageInfo).Error; err != nil {
		tx.Rollback()
		return dto.ImageToEssay{}, fmt.Errorf("保存图片信息失败: %v", err)
	}

	// 关键修复：直接使用原始URL（不再编码）
	// 因为OcrHandwritingWithBaidu内部会进行表单编码
	essay, err := OcrHandwritingWithBaidu(param.ImageURL)
	if err != nil {
		tx.Rollback()
		return dto.ImageToEssay{}, fmt.Errorf("解析图片文字失败: %v", err)
	}

	// 3. 保存文章信息（事务内）
	essayInfo := dao.Essay{
		UserID:     userId,
		Content:    essay,
		SubmitTime: nowTime,
		Title:      param.Title,
		//WordCount:    len([]rune(essay)),
		//UploadMethod: 0, // 图片上传
		ImageId: imageInfo.ID,
	}
	if err = tx.Create(&essayInfo).Error; err != nil {
		tx.Rollback()
		return dto.ImageToEssay{}, fmt.Errorf("保存文章信息失败: %v", err)
	}

	//// 4. 更新图片记录的文章ID
	//if err = tx.Model(&dao.ImageInfo{}).Where("id = ?", imageInfo.ID).
	//	Update("essay_id", essayInfo.ID).Error; err != nil {
	//	tx.Rollback()
	//	return dto.ImageToEssay{}, fmt.Errorf("更新图片记录失败: %v", err)
	//}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		return dto.ImageToEssay{}, fmt.Errorf("提交事务失败: %v", err)
	}

	// 返回信息
	return dto.ImageToEssay{
		EssayId:       essayInfo.ID,
		SubmitTime:    essayInfo.SubmitTime,
		SubmitTimeMar: utils.MarshalTime(essayInfo.SubmitTime),
		Content:       essay,
	}, nil
}

// SaveEssay 保存用户输入的文章
func SaveEssay(c *gin.Context, userId int64, param dto.SaveEssayMap) (data dto.ImageToEssay, err error) {
	// 开启事务
	tx := global.GetDbConn(c).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("事务回滚: %v", r)
		}
	}()

	nowTime := time.Now()

	// 保存文章信息
	essayInfo := dao.Essay{
		UserID:     userId,
		Content:    param.Content,
		SubmitTime: nowTime,
		// Title:      param.Title,
	}

	if err = tx.Create(&essayInfo).Error; err != nil {
		tx.Rollback()
		return dto.ImageToEssay{}, fmt.Errorf("保存文章信息失败: %v", err)
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		return dto.ImageToEssay{}, fmt.Errorf("提交事务失败: %v", err)
	}

	// 返回信息
	return dto.ImageToEssay{
		EssayId:       essayInfo.ID,
		SubmitTime:    essayInfo.SubmitTime,
		SubmitTimeMar: utils.MarshalTime(essayInfo.SubmitTime),
		Content:       param.Content,
	}, nil
}

// UpdateEssayContent 修改文章内容
func UpdateEssayContent(c *gin.Context, userId int64, param dto.UpdateEssayContentMap) (data dto.ImageToEssay, err error) {
	// 开启事务
	tx := global.GetDbConn(c).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("事务回滚: %v", r)
		}
	}()

	// 检查文章是否存在且属于当前用户
	var essay dao.Essay
	if err = tx.Where("id = ? AND user_id = ?", param.EssayId, userId).First(&essay).Error; err != nil {
		tx.Rollback()
		return dto.ImageToEssay{}, fmt.Errorf("文章不存在或无权限修改: %w", err)
	}

	// 更新文章内容
	if err = tx.Model(&essay).Update("content", param.Content).Error; err != nil {
		tx.Rollback()
		return dto.ImageToEssay{}, fmt.Errorf("更新文章内容失败: %w", err)
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		return dto.ImageToEssay{}, fmt.Errorf("提交事务失败: %w", err)
	}

	// 返回更新后的文章信息
	return dto.ImageToEssay{
		EssayId:       essay.ID,
		SubmitTime:    essay.SubmitTime,
		SubmitTimeMar: utils.MarshalTime(essay.SubmitTime),
		Content:       param.Content,
	}, nil
}

// RestoreMultiImageInfo 处理多页图片信息
func RestoreMultiImageInfo(c *gin.Context, userId int64, param dto.RestoreMultiImageInfoMap) (data dto.ImageToEssay, err error) {
	tx := global.GetDbConn(c).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("事务回滚: %v", r)
		}
	}()

	nowTime := time.Now()
	var allText strings.Builder
	var imageIds []int

	for i, page := range param.Pages {
		imageInfo := dao.ImageInfo{
			ImageURL:   page.ImageURL,
			ImageName:  page.ImageName,
			UploadTime: nowTime,
		}
		if err = tx.Create(&imageInfo).Error; err != nil {
			tx.Rollback()
			return dto.ImageToEssay{}, fmt.Errorf("保存第%d页图片信息失败: %v", i+1, err)
		}
		imageIds = append(imageIds, imageInfo.ID)

		pageText, err := OcrHandwritingWithBaidu(page.ImageURL)
		if err != nil {
			tx.Rollback()
			return dto.ImageToEssay{}, fmt.Errorf("识别第%d页图片文字失败: %v", i+1, err)
		}

		if i > 0 {
			allText.WriteString("\n\n")
		}
		allText.WriteString(pageText)
	}

	essayInfo := dao.Essay{
		UserID:     userId,
		Content:    allText.String(),
		SubmitTime: nowTime,
		Title:      param.Title,
		ImageId:    imageIds[0],
	}
	if err = tx.Create(&essayInfo).Error; err != nil {
		tx.Rollback()
		return dto.ImageToEssay{}, fmt.Errorf("保存文章信息失败: %v", err)
	}

	if err = tx.Commit().Error; err != nil {
		return dto.ImageToEssay{}, fmt.Errorf("提交事务失败: %v", err)
	}

	return dto.ImageToEssay{
		EssayId:       essayInfo.ID,
		SubmitTime:    essayInfo.SubmitTime,
		SubmitTimeMar: utils.MarshalTime(essayInfo.SubmitTime),
		Content:       allText.String(),
	}, nil
}
