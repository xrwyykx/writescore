package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"net/http"
	"writescore/data/db/comon"
	"writescore/global"
)

func RecognizeText(c *gin.Context, imageUrl string) {
	text, err := comon.OcrHandwritingWithBaidu(imageUrl)
	if err != nil {
		fmt.Println("识别失败:", err)
		return
	}
	fmt.Println("识别结果:\n", text)
}

//func GetUploadToken(c *gin.Context) {
//	mac := auth.New(global.QiniuyunAK, global.QiniuyunSK)
//	putPolicy := storage.PutPolicy{
//		Scope: global.Bucket,
//	}
//	uptoken := putPolicy.UploadToken(mac)
//	c.JSON(http.StatusOK, co.Success("上传凭证获取成功", uptoken))
//}

func GetUploadToken(c *gin.Context) {
	mac := auth.New(global.QiniuyunAK, global.QiniuyunSK)
	putPolicy := storage.PutPolicy{
		Scope: global.Bucket,
	}
	uptoken := putPolicy.UploadToken(mac)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    uptoken,
		"message": "上传凭证获取成功",
	})
}
