package common

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Response struct {
	Status  string      `json:"status"`            // "ok" or "error"
	Message string      `json:"message,omitempty"` // 提示信息
	Data    interface{} `json:"data,omitempty"`    // 返回数据
}

// 统一发送 WS 响应
func SendWSResponse(c *websocket.Conn, status, msg string, data interface{}) error {
	resp := Response{
		Status:  status,
		Message: msg,
		Data:    data,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return c.WriteMessage(websocket.TextMessage, b)
}

// 统一发送 HTTP JSON 响应
func SendHTTPResponse(c *fiber.Ctx, status, msg string, data interface{}) error {
	return c.JSON(Response{
		Status:  status,
		Message: msg,
		Data:    data,
	})
}
