package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"vnfm-api/internal/middleware"
	"vnfm-api/internal/user/dto"
	"vnfm-api/internal/user/oauth"
	"vnfm-api/internal/user/repository"
	"vnfm-api/pkg/config"
	"vnfm-api/pkg/errors"

	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	users       *repository.UserRepository
	rdb         *redis.Client
	oauthClient *oauth.Client
	oidc        config.OIDCConfig
}

func NewAuthService(
	users *repository.UserRepository,
	rdb *redis.Client,
	oauthClient *oauth.Client,
	oidc config.OIDCConfig,
) *AuthService {
	return &AuthService{users: users, rdb: rdb, oauthClient: oauthClient, oidc: oidc}
}

func (s *AuthService) OAuthCallback(ctx context.Context, req *dto.OAuthCallbackRequest) (*dto.SessionResponse, *errors.AppError) {
	if req.Code == "" || req.CodeVerifier == "" {
		return nil, errors.ErrBadRequest("缺少授权码")
	}
	tokenResp, err := s.oauthClient.ExchangeCode(req.Code, req.CodeVerifier)
	if err != nil {
		if oauth.IsBanned(err) {
			return nil, errors.ErrAccountBanned()
		}
		return nil, errors.ErrBadRequest(fmt.Sprintf("OAuth 授权码交换失败: %v", err))
	}
	oauthUser, err := s.oauthClient.FetchUserInfo(tokenResp.AccessToken)
	if err != nil {
		if oauth.IsBanned(err) {
			return nil, errors.ErrAccountBanned()
		}
		return nil, errors.ErrBadRequest(fmt.Sprintf("获取 OAuth 用户信息失败: %v", err))
	}
	if oauthUser.ID <= 0 {
		return nil, errors.ErrInternal("OAuth /oauth/userinfo 未返回用户 id")
	}

	if err := s.users.Ensure(int64(oauthUser.ID), oauthUser.Name, oauthUser.Picture); err != nil {
		return nil, errors.ErrInternal("初始化用户失败")
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, errors.ErrInternal("生成会话令牌失败")
	}

	exp := tokenResp.ExpiresIn
	if exp <= 0 {
		exp = 900
	}
	sessionData := middleware.SessionData{
		UserInfo: middleware.UserInfo{
			ID:     int64(oauthUser.ID),
			Sub:    oauthUser.Sub,
			Name:   oauthUser.Name,
			Email:  oauthUser.Email,
			Avatar: oauthUser.Picture,
		},
		OAuthAccessToken:  tokenResp.AccessToken,
		OAuthRefreshToken: tokenResp.RefreshToken,
		OAuthExpiresAt:    time.Now().Unix() + int64(exp),
	}
	data, err := json.Marshal(sessionData)
	if err != nil {
		return nil, errors.ErrInternal("序列化会话数据失败")
	}
	if err := s.rdb.Set(ctx, middleware.SessionKey(sessionToken), data, middleware.SessionTTL).Err(); err != nil {
		return nil, errors.ErrInternal("写入会话失败")
	}

	profile, appErr := s.profileFromLocal(int64(oauthUser.ID), oauthUser.Sub, oauthUser.Email)
	if appErr != nil {
		return nil, appErr
	}
	return &dto.SessionResponse{Token: sessionToken, User: profile}, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionToken string) error {
	val, err := s.rdb.Get(ctx, middleware.SessionKey(sessionToken)).Result()
	if err == nil {
		var session middleware.SessionData
		if json.Unmarshal([]byte(val), &session) == nil && session.OAuthRefreshToken != "" {
			_ = s.oauthClient.RevokeToken(session.OAuthRefreshToken)
		}
	}
	return s.rdb.Del(ctx, middleware.SessionKey(sessionToken)).Err()
}

func (s *AuthService) GetProfile(id int64, sub, email string) (*dto.UserProfile, *errors.AppError) {
	return s.profileFromLocal(id, sub, email)
}

func (s *AuthService) UpdatePrefs(id int64, req *dto.UpdatePrefsRequest) (*dto.UserProfile, *errors.AppError) {
	if req.LanguagePreference != nil {
		lang := *req.LanguagePreference
		if lang != "zh" && lang != "ja" {
			return nil, errors.ErrBadRequest("语言仅支持 zh 或 ja")
		}
	}
	if req.ThemePreference != nil {
		theme := *req.ThemePreference
		if theme != "light" && theme != "dark" && theme != "system" {
			return nil, errors.ErrBadRequest("主题仅支持 light、dark 或 system")
		}
	}
	if err := s.users.UpdatePrefs(id, req.LanguagePreference, req.ThemePreference, req.DisplayMembershipID, req.ClearDisplayClub); err != nil {
		return nil, errors.ErrInternal("保存偏好失败")
	}
	return s.profileFromLocal(id, "", "")
}

func (s *AuthService) profileFromLocal(id int64, sub, email string) (*dto.UserProfile, *errors.AppError) {
	row, err := s.users.FindByID(id)
	if err != nil {
		return nil, errors.ErrNotFound("用户不存在")
	}
	return &dto.UserProfile{
		ID:                  row.ID,
		Sub:                 sub,
		Name:                row.Name,
		Avatar:              row.AvatarURL,
		Email:               email,
		LanguagePreference:  row.LanguagePreference,
		ThemePreference:     row.ThemePreference,
		DisplayMembershipID: row.DisplayMembershipID,
		AccountCenterURL:    s.oidc.FrontendURL + "/profile",
	}, nil
}

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
