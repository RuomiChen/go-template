package auth

import (
	"fmt"
	"mvc/internal/admin"
	"mvc/internal/user"
	"mvc/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	adminService *admin.Service
	userService  *user.Service
}

func NewHandler(adminService *admin.Service, userService *user.Service) *Handler {
	return &Handler{adminService: adminService, userService: userService}
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var auth Auth
	if err := c.BodyParser(&auth); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request")
	}
	fmt.Print(auth.Role)

	var loginResp interface{}
	var err error

	switch auth.Role {
	case 1: // 管理员登录
		loginResp, err = h.adminService.Login(c, auth.Username, auth.Password)
	case 0: // 普通用户登录
		loginResp, err = h.userService.Login(c, auth.Username, auth.Password)
	default:
		return response.Error(c, fiber.StatusBadRequest, "invalid role")
	}

	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	return response.Success(c, loginResp)
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var auth Auth

	if err := c.BodyParser(&auth); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid body")
	}

	if auth.Username == "" || auth.Password == "" {
		return response.Error(c, fiber.StatusBadRequest, "username or password empty")
	}

	switch auth.Role {
	case 1:
		// 管理员注册
		err := h.adminService.Register(c, auth.Username, auth.Password)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error())
		}
	case 0:
		// 普通用户注册
		err := h.userService.Register(auth.Username, auth.Password)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error())
		}
	default:
		return response.Error(c, fiber.StatusBadRequest, "invalid role")
	}

	return response.Success(c, nil)
}
