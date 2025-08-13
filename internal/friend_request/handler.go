package friend_request

import (
	"context"
	"encoding/json"
	"mvc/pkg/contextkeys"
	"mvc/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
)

type FriendHandler struct {
	service *FriendService
	logger  zerolog.Logger
}

type AddFriendRequest struct {
	ToUserID uint64 `json:"to_user_id"`
	Message  string `json:"message"`
}

// 新增 AddFriend HTTP 处理函数
func (h *FriendHandler) AddFriend(c *fiber.Ctx) error {
	var req AddFriendRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	// 从中间件解析的上下文中取当前用户 ID
	IDVal := c.Locals("id")
	h.logger.Info().Interface("id", IDVal).Msg("token id")
	ID, ok := IDVal.(int)
	if !ok || ID == 0 {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	// 把请求参数转成 json.RawMessage 传给 Service
	rawData, err := json.Marshal(req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "internal error")
	}

	// 传入带有用户 ID 的 context，方便 Service 取用户 ID
	ctx := context.WithValue(c.UserContext(), contextkeys.IDKey, uint(ID))

	err = h.service.AddFriend(ctx, rawData)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, nil)
}

func NewFriendHandler(service *FriendService, logger zerolog.Logger) *FriendHandler {
	return &FriendHandler{service: service, logger: logger}
}

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (h *FriendHandler) Handle(c *websocket.Conn) {
	defer c.Close()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			h.logger.Info().Msgf("ws read error: %v", err)
			break
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			h.sendError(c, "invalid message format")
			continue
		}

		switch message.Type {
		case "add_friend":
			err := h.service.AddFriend(context.Background(), message.Data)
			if err != nil {
				h.sendError(c, err.Error())
			} else {
				h.sendSuccess(c, "friend added successfully")
			}
		// case "request_list":
		// 	err := h.service.RequestList(context.Background(), message.Data)
		// 	if err != nil {
		// 		h.sendError(c, err.Error())
		// 	} else {
		// 		h.sendSuccess(c, "friend added successfully")
		// 	}
		default:
			h.sendError(c, "unknown message type: "+message.Type)
		}
	}
}

func (h *FriendHandler) sendError(c *websocket.Conn, msg string) {
	resp := map[string]string{"status": "error", "message": msg}
	b, _ := json.Marshal(resp)
	_ = c.WriteMessage(websocket.TextMessage, b)
}

func (h *FriendHandler) sendSuccess(c *websocket.Conn, msg string) {
	resp := map[string]string{"status": "ok", "message": msg}
	b, _ := json.Marshal(resp)
	_ = c.WriteMessage(websocket.TextMessage, b)
}
