package comon

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"writescore/global"
	"writescore/models/dao"
	"writescore/models/dto"
	"writescore/utils"
)

// GetAccessToken 获取新的 access_token
func GetAccessToken(apiKey, secretKey string) (string, error) {
	tokenURL := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", apiKey, secretKey)
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

// OcrHandwritingWithBaidu 使用百度OCR识别手写文字
func OcrHandwritingWithBaidu(imgURL string) (string, error) {
	apiKey := "aLnqDRQEwLa7yX7XhJJKa1X7"
	secretKey := "6IkAaJhaZAGPKmaJADy2EPAdkkUauyVD"

	// 1. 获取Access Token
	accessToken, err := GetAccessToken(apiKey, secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %v", err)
	}

	// 2. 调用手写OCR接口
	ocrURL := fmt.Sprintf("https://aip.baidubce.com/rest/2.0/ocr/v1/handwriting?access_token=%s", accessToken)
	formData := url.Values{"url": {imgURL}, "recognize_granularity": {"big"}} // 大颗粒度适合手写

	resp, err := http.Post(ocrURL, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to call OCR API: %v", err)
	}
	defer resp.Body.Close()

	// 3. 解析结果（百度返回JSON结构较复杂，需多层解析）
	ocrBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read OCR response body: %v", err)
	}

	var ocrResult map[string]interface{}
	err = json.Unmarshal(ocrBody, &ocrResult)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal OCR response: %v", err)
	}

	// 检查是否包含错误信息
	if _, ok := ocrResult["error_code"]; ok {
		errorMsg, ok := ocrResult["error_msg"].(string)
		if ok {
			return "", fmt.Errorf("OCR failed: %s", errorMsg)
		}
		return "", fmt.Errorf("OCR failed with unknown error")
	}

	wordsResult, ok := ocrResult["words_result"].([]interface{})
	if !ok {
		return "", fmt.Errorf("words_result not found in OCR response")
	}

	var text strings.Builder
	for _, word := range wordsResult {
		wordMap, ok := word.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("invalid word format in OCR response")
		}
		text.WriteString(wordMap["words"].(string) + "\n")
	}

	return text.String(), nil
}

func RestoreImageInfo(c *gin.Context, userId int64, param dto.RestoreImageInfoMap) (data dto.ImageToEssay, err error) {
	//第一个事务是将照片信息保存在照片表，然后接着根据前面的照片信息将照片解析成文字，不是同时进行的
	//如果前面的事务没有完成后面的事务也不能顺利进行，应该不需要启动事务？
	//照片信息不需要关联用户id，如果解析照片文字成功就可以把essay_id更新到image_info表，这个过程中需不需要开启事务
	//1.保存图片信息
	nowTime := time.Now()
	imageInfo := dao.ImageInfo{
		ImageURL:   param.ImageURL,
		ImageName:  param.ImageName,
		UploadTime: nowTime,
	}
	if err := global.GetDbConn(c).Model(&dao.ImageInfo{}).Create(&imageInfo).Error; err != nil {
		return dto.ImageToEssay{}, err
	} //保存图片信息失败时
	//2.解析图片文字,并将这篇文章的信息更新到数据库表
	essay, err := OcrHandwritingWithBaidu(param.ImageURL)
	if err != nil {
		return dto.ImageToEssay{}, err
	}
	essayInfo := dao.Essay{
		UserID:     userId,
		Content:    essay,
		SubmitTime: nowTime,
		Title:      param.Title,
		//WordCount: ,//计算字符串长度
		//UploadMethod: 0,            //图片上传
		ImageId: imageInfo.ID, //实现与image_info表的关联
	}
	if err := global.GetDbConn(c).Model(&dao.Essay{}).Create(&essayInfo).Error; err != nil {
		return dto.ImageToEssay{}, err
	}

	//返回信息
	data.EssayId = essayInfo.ID
	data.SubmitTime = nowTime
	data.SubmitTimeMar = utils.MarshalTime(data.SubmitTime)
	data.Content = essay
	return data, nil
}
