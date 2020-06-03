package serializer

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

// 响应体
type Response struct {
	Code  int         `json:"code"`            // 响应代码
	Data  interface{} `json:"data,omitempty"`  // 数据
	Msg   string      `json:"msg"`             // 消息
	Error string      `json:"error,omitempty"` // 错误
}

func CreateSuccessResponse(data interface{}, msg string) Response {
	return Response{
		Code: CodeSuccess,
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

func CreateDBErrorResponse(msg string, err error) Response {
	if msg == "" {
		msg = "数据库操作失败"
	}
	return CreateErrorResponse(CodeDBError, msg, err)
}

func CreateInternalErrorResponse(msg string, err error) Response {
	if msg == "" {
		msg = "系统内部错误"
	}
	return CreateErrorResponse(CodeInternalError, msg, err)
}

func CreateGeneralParamErrorResponse(msg string, err error) Response {
	if msg == "" {
		msg = "参数错误"
	}
	return CreateErrorResponse(CodeParamError, msg, err)
}

// CreateParamErrorMsg 根据Validator返回的错误信息给出错误提示
func CreateParamErrorMsg(filed string, tag string) string {
	// 未通过验证的表单域与中文对应
	fieldMap := map[string]string{
		"UserName": "邮箱",
		"Password": "密码",
		"NickName": "昵称",
	}
	// 未通过的规则与中文对应
	tagMap := map[string]string{
		"required": "不能为空",
		"min":      "太短",
		"max":      "太长",
		"email":    "格式不正确",
	}
	fieldVal, findField := fieldMap[filed]
	tagVal, findTag := tagMap[tag]
	if findField && findTag {
		// 返回拼接出来的错误信息
		return fieldVal + tagVal
	}
	return ""
}

// CreateParamErrorResponse 返回错误消息
func CreateParamErrorResponse(err error) Response {
	// 处理 Validator 产生的错误
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			return CreateGeneralParamErrorResponse(
				CreateParamErrorMsg(e.Field(), e.Tag()),
				err,
			)
		}
	}
	// 处理 JSON 类型错误
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return CreateGeneralParamErrorResponse("JSON类型不匹配", err)
	}
	// 其他参数错误
	return CreateGeneralParamErrorResponse("参数错误", err)
}
