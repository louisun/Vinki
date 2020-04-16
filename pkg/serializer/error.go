package serializer

const (
	SuccessCode       = 200
	InternalErrorCode = 1000 // 内部错误
	DBErrorCode       = 1001 // 数据库错误
	AuthErrorCode     = 1002 // 认证错误
	ParamErrorCode    = 1003 // 参数错误
)

type ServiceError struct {
	Code     int
	Msg      string
	RawError error
}

func NewServiceError(code int, msg string, err error) ServiceError {
	return ServiceError{
		Code:     code,
		Msg:      msg,
		RawError: err,
	}
}

func (se *ServiceError) WrapError(err error) ServiceError {
	se.RawError = err
	return *se
}

// 实现 error 接口，返回信息为内部 Msg
func (se ServiceError) Error() string {
	return se.Msg
}
