package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"vnfm-api/internal/user/oauth"
	"vnfm-api/pkg/errors"
	"vnfm-api/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

type contextKey string

const (
	UserInfoKey         contextKey = "userInfo"
	OAuthAccessTokenKey contextKey = "oauthAccessToken"
	SessionCookieName              = "vnfm_session"
	SessionPrefix                  = "vnfm:session:v1:"
	SessionTTL                     = 90 * 24 * time.Hour
	sessionRenewPrefix             = "vnfm:session-renew:"
)

func SessionKey(token string) string { return SessionPrefix + token }

type UserInfo struct {
	ID     int64  `json:"id"`
	Sub    string `json:"sub"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

type SessionData struct {
	UserInfo
	OAuthAccessToken  string `json:"oauth_access_token"`
	OAuthRefreshToken string `json:"oauth_refresh_token"`
	OAuthExpiresAt    int64  `json:"oauth_expires_at"`
}

func Auth(rdb *redis.Client, oauthClient *oauth.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		if err := hydrateSession(c, rdb, oauthClient, true); err != nil {
			return response.Error(c, err)
		}
		if GetUser(c) == nil {
			return response.Error(c, errors.ErrAuthExpired())
		}
		return c.Next()
	}
}

func OptionalAuth(rdb *redis.Client, oauthClient *oauth.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		_ = hydrateSession(c, rdb, oauthClient, false)
		return c.Next()
	}
}

func hydrateSession(c fiber.Ctx, rdb *redis.Client, oauthClient *oauth.Client, required bool) *errors.AppError {
	token := c.Cookies(SessionCookieName)
	if token == "" {
		if required {
			return errors.ErrAuthExpired()
		}
		return nil
	}
	ctx := c.Context()
	val, err := rdb.Get(ctx, SessionKey(token)).Result()
	if err != nil {
		if required {
			return errors.ErrAuthExpired()
		}
		return nil
	}
	var session SessionData
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		if required {
			return errors.ErrAuthExpired()
		}
		return nil
	}

	const refreshSkew = 30 * time.Second
	if session.OAuthExpiresAt > 0 && time.Now().Add(refreshSkew).Unix() > session.OAuthExpiresAt {
		lockKey := "vnfm:refresh_lock:" + token
		var refreshErr *errors.AppError
		if locked, _ := rdb.SetNX(ctx, lockKey, "1", 15*time.Second).Result(); locked {
			refreshErr = refreshSession(ctx, rdb, oauthClient, token, &session)
			rdb.Del(ctx, lockKey)
		} else {
			time.Sleep(200 * time.Millisecond)
			if val2, gerr := rdb.Get(ctx, SessionKey(token)).Result(); gerr == nil {
				_ = json.Unmarshal([]byte(val2), &session)
			}
		}
		if refreshErr != nil {
			if required {
				return refreshErr
			}
			return nil
		}
	}

	renewSlidingSession(c, rdb, token)
	c.Locals(string(UserInfoKey), &session.UserInfo)
	c.Locals(string(OAuthAccessTokenKey), session.OAuthAccessToken)
	return nil
}

func refreshSession(
	ctx context.Context,
	rdb *redis.Client,
	oauthClient *oauth.Client,
	token string,
	session *SessionData,
) *errors.AppError {
	refreshed, err := oauthClient.RefreshOAuthToken(session.OAuthRefreshToken)
	if err != nil {
		if oauth.IsBanned(err) {
			rdb.Del(ctx, SessionKey(token))
			return errors.ErrAccountBanned()
		}
		if oauth.IsRefreshTokenDead(err) {
			rdb.Del(ctx, SessionKey(token))
			return errors.ErrAuthExpired()
		}
		slog.Warn("oauth refresh failed", "error", err)
		return errors.ErrAuthExpired()
	}
	session.OAuthAccessToken = refreshed.AccessToken
	if refreshed.RefreshToken != "" {
		session.OAuthRefreshToken = refreshed.RefreshToken
	}
	exp := refreshed.ExpiresIn
	if exp <= 0 {
		exp = 900
	}
	session.OAuthExpiresAt = time.Now().Unix() + int64(exp)
	data, _ := json.Marshal(session)
	rdb.Set(ctx, SessionKey(token), data, SessionTTL)
	return nil
}

func renewSlidingSession(c fiber.Ctx, rdb *redis.Client, token string) {
	ctx := c.Context()
	ok, _ := rdb.SetNX(ctx, sessionRenewPrefix+token, "1", SessionTTL/2).Result()
	if !ok {
		return
	}
	rdb.Expire(ctx, SessionKey(token), SessionTTL)
	secure := c.Get("X-Forwarded-Proto") == "https" || c.Protocol() == "https"
	c.Cookie(&fiber.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		MaxAge:   int(SessionTTL.Seconds()),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax",
		Path:     "/",
	})
}

func GetUser(c fiber.Ctx) *UserInfo {
	info, ok := c.Locals(string(UserInfoKey)).(*UserInfo)
	if !ok {
		return nil
	}
	return info
}

func MustGetUser(c fiber.Ctx) (*UserInfo, *errors.AppError) {
	info := GetUser(c)
	if info == nil {
		return nil, errors.ErrAuthExpired()
	}
	return info, nil
}
