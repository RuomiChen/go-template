package friend_relation

import (
	"context"
	"encoding/json"

	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
)

type FriendHandler struct {
	service *FriendService
	logger  zerolog.Logger
}

func NewHandler(service *FriendService, logger zerolog.Logger) *FriendHandler {
	return &FriendHandler{service: service, logger: logger}
}

type wsRequest struct {
	Action string          `json:"action"` // add_friend / remove_friend / list_friends
	Data   json.RawMessage `json:"data"`
}

type addFriendRequest struct {
	UserA uint `json:"user_a"`
	UserB uint `json:"user_b"`
}

type removeFriendRequest struct {
	UserA uint `json:"user_a"`
	UserB uint `json:"user_b"`
}

type listFriendsRequest struct {
	UserID uint `json:"user_id"`
}

func (h *FriendHandler) Handle(c *websocket.Conn) {
	defer c.Close()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		var req wsRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			h.sendError(c, "invalid request format")
			continue
		}

		switch req.Action {
		case "add_friend":
			var data addFriendRequest
			if err := json.Unmarshal(req.Data, &data); err != nil {
				h.sendError(c, "invalid add_friend data")
				continue
			}
			err := h.service.AddFriend(context.Background(), data.UserA, data.UserB)
			if err != nil {
				h.sendError(c, err.Error())
			} else {
				h.sendSuccess(c, "friend added")
			}

		case "remove_friend":
			var data removeFriendRequest
			if err := json.Unmarshal(req.Data, &data); err != nil {
				h.sendError(c, "invalid remove_friend data")
				continue
			}
			err := h.service.RemoveFriend(context.Background(), data.UserA, data.UserB)
			if err != nil {
				h.sendError(c, err.Error())
			} else {
				h.sendSuccess(c, "friend removed")
			}

		case "list_friends":
			var data listFriendsRequest
			if err := json.Unmarshal(req.Data, &data); err != nil {
				h.sendError(c, "invalid list_friends data")
				continue
			}
			list, err := h.service.ListFriends(context.Background(), data.UserID)
			if err != nil {
				h.sendError(c, err.Error())
			} else {
				h.sendData(c, list)
			}

		default:
			h.sendError(c, "unknown action: "+req.Action)
		}
	}
}

func (h *FriendHandler) sendError(c *websocket.Conn, msg string) {
	resp := map[string]interface{}{
		"status":  "error",
		"message": msg,
	}
	_ = c.WriteJSON(resp)
}

func (h *FriendHandler) sendSuccess(c *websocket.Conn, msg string) {
	resp := map[string]interface{}{
		"status":  "ok",
		"message": msg,
	}
	_ = c.WriteJSON(resp)
}

func (h *FriendHandler) sendData(c *websocket.Conn, data interface{}) {
	resp := map[string]interface{}{
		"status": "ok",
		"data":   data,
	}
	_ = c.WriteJSON(resp)
}
