package friend

import (
	"encoding/json"

	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
)

type FriendHandler struct {
	service *FriendService
	logger  zerolog.Logger
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
			err := h.service.AddFriend(message.Data)
			if err != nil {
				h.sendError(c, err.Error())
			} else {
				h.sendSuccess(c, "friend added successfully")
			}
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
