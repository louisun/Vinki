package serializer

const (
	CodeSuccess          = 200
	CodeUnauthorized     = 401  // 认证错误
	CodeForbidden        = 403  // 无权限
	CodeInternalError    = 1000 // 内部错误
	CodeDBError          = 1001 // 数据库错误
	CodeParamError       = 1002 // 参数错误
	CodeConditionNotMeet = 1003 // 条件不满足错误
	CodeAdminRequired    = 2000 // 需要为管理员账号
	CodeActiveRequired   = 2001 // 需要为激活账号
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
