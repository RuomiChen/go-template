package response

import "github.com/gofiber/fiber/v2"

type HTTPResponse struct {
	Code    int         `json:"code"` // 业务码
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(c *fiber.Ctx, httpStatus, code int, message string, data interface{}) error {
	return c.Status(httpStatus).JSON(HTTPResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Success(c *fiber.Ctx, data interface{}) error {
	return JSON(c, fiber.StatusOK, CodeSuccess, "success", data)
}

func Error(c *fiber.Ctx, httpStatus, code int, message string) error {
	return JSON(c, httpStatus, code, message, nil)
}
