package response

// 业务错误码（0 表示成功）
// 1000-1999 通用错误
// 2000-2999 用户相关
// 3000-3999 新闻公告相关
const (
	// 通用
	CodeSuccess      = 0
	CodeInvalidParam = 1001
	CodeInternalErr  = 1002
	CodeUnauthorized = 1003

	// 用户
	CodeUserNotFound = 2001
	CodeUserExists   = 2002
	CodePasswordErr  = 2003

	// 新闻公告
	CodeNewsNotFound     = 3001
	CodeNewsExpired      = 3002
	CodeNewsUpdatedError = 3003
)
