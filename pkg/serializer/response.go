package serializer

import (
	"github.com/gin-gonic/gin"
)

// 响应体
type Response struct {
	Code  int         `json:"code"`            // 响应代码
	Data  interface{} `json:"data,omitempty"`  // 数据
	Msg   string      `json:"msg"`             // 消息
	Error string      `json:"error,omitempty"` // 错误
}

func SuccessResponse(data interface{}, msg string) Response {
	return Response{
		Code: SuccessCode,
		Data: data,
		Msg:  msg,
	}
}

// CreateErrorResponse 创建错误响应体
func CreateErrorResponse(errCode int, msg string, err error) Response {
	// 如果 err 是 ServiceError 类型，则覆盖参数传入的错误内容
	if serviceError, ok := err.(ServiceError); ok {
		errCode = serviceError.Code
		err = serviceError.RawError
		msg = serviceError.Msg
	}

	response := Response{
		// 无 Data，Error 只在非生产环境设置
		Code: errCode,
		Msg:  msg,
	}
	if err != nil && gin.Mode() != gin.ReleaseMode {
		response.Error = err.Error()
	}
	return response
}

func DBErrorResponse(msg string, err error) Response {
	if msg == "" {
		msg = "数据库操作失败"
	}
	return CreateErrorResponse(DBErrorCode, msg, err)
}

func ParamErrorResponse(msg string, err error) Response {
	if msg == "" {
		msg = "参数错误"
	}
	return CreateErrorResponse(DBErrorCode, msg, err)
}

func InternalErrorResponse(msg string, err error) Response {
	if msg == "" {
		msg = "系统内部错误"
	}
	return CreateErrorResponse(InternalErrorCode, msg, err)
}
