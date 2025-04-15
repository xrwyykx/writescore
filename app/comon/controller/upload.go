package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"net/http"
	"writescore/app"
	"writescore/data/db/comon"
	"writescore/global"
	"writescore/models/co"
	"writescore/models/dto"
)

func RecognizeText(c *gin.Context) {
	var data dto.Shangchuan
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("参数绑定失败"))
		return
	}
	text, err := comon.OcrHandwritingWithBaidu(data.ImageURL)
	if err != nil {
		fmt.Println("识别失败:", err)
		return
	}
	fmt.Println("识别结果:\n", text)
}

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

func RestoreImageInfo(c *gin.Context) {
	userId := app.GetUserId(c)
	var param dto.RestoreImageInfoMap
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("参数绑定失败"))
		return
	}
	data, err := comon.RestoreImageInfo(c, userId, param)
	//获取到照片的信息返回保存到数据库,并将识别出来的文字返回，这个地方应该要启动事务，因为
	//同时需要进行两个操作，避免有一个事件没有顺利完成而另一个事务完成了导致数据不同步问题
	if err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("解析图片信息失败,"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, co.Success("解析图片信息成功", data))
}
