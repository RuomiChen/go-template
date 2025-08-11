package response

type WSResponse struct {
	Code    int         `json:"code"` // 业务码
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 构造成功响应
func WSSuccess(data interface{}) WSResponse {
	return WSResponse{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	}
}

// 构造错误响应
func WSError(code int, message string) WSResponse {
	return WSResponse{
		Code:    code,
		Message: message,
	}
}
