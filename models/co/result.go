package co

type ResultModel struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message interface{} `json:"message"`
}

const (
	success              = 0     //请求成功
	redirect             = 30001 //重定向
	badRequest           = 40000 //客户端请求参数错误
	notLoggedIn          = 40100 //未登录
	notAuth              = 40101 //请求成功，但是无权限
	notData              = 40400 //请求数据不存在
	forbiddenAccess      = 40300 //禁止访问
	limitRate            = 40301 //限流
	internal             = 50000 //服务端内部错误
	controlError         = 50010 //操作失败
	callError            = 50010 //接口调用失败
	internalCallError    = 50011 //服务内部之间接口调用失败
	controlDataBaseError = 50020 //数据库操作失败
)

func Success(message interface{}, data interface{}) ResultModel {
	var result ResultModel
	result.Code = success
	result.Data = data
	result.Message = message
	return result
}
func BadRequest(message string) ResultModel {
	var result ResultModel
	result.Code = badRequest
	result.Message = message
	result.Data = -1
	return result
}
