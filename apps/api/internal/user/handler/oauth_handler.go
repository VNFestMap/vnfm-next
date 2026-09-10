package handler

import (
	"vnfm-api/internal/middleware"
	"vnfm-api/internal/user/dto"
	"vnfm-api/internal/user/service"
	"vnfm-api/pkg/response"
	"vnfm-api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

type OAuthHandler struct {
	auth   *service.AuthService
	secure bool
}

func NewOAuthHandler(auth *service.AuthService, isProd bool) *OAuthHandler {
	return &OAuthHandler{auth: auth, secure: isProd}
}

func (h *OAuthHandler) Callback(c fiber.Ctx) error {
	var req dto.OAuthCallbackRequest
	if err := utils.BindJSON(c, &req); err != nil {
		return response.Error(c, err)
	}
	session, appErr := h.auth.OAuthCallback(c.Context(), &req)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	c.Cookie(&fiber.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    session.Token,
		MaxAge:   int(middleware.SessionTTL.Seconds()),
		HTTPOnly: true,
		Secure:   h.secure,
		SameSite: "Lax",
		Path:     "/",
	})
	return response.OK(c, session.User)
}

func (h *OAuthHandler) Logout(c fiber.Ctx) error {
	token := c.Cookies(middleware.SessionCookieName)
	if token != "" {
		_ = h.auth.Logout(c.Context(), token)
	}
	c.Cookie(&fiber.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   h.secure,
		SameSite: "Lax",
		Path:     "/",
	})
	return response.OKMessage(c, "已登出")
}

func (h *OAuthHandler) Me(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	profile, appErr := h.auth.GetProfile(user.ID, user.Sub, user.Email)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	if profile.Sub == "" {
		profile.Sub = user.Sub
	}
	if profile.Email == "" {
		profile.Email = user.Email
	}
	if profile.Avatar == "" {
		profile.Avatar = user.Avatar
	}
	if profile.Name == "" {
		profile.Name = user.Name
	}
	return response.OK(c, profile)
}

func (h *OAuthHandler) UpdateMe(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	var req dto.UpdatePrefsRequest
	if bindErr := utils.BindJSON(c, &req); bindErr != nil {
		return response.Error(c, bindErr)
	}
	profile, appErr := h.auth.UpdatePrefs(user.ID, &req)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, profile)
}
