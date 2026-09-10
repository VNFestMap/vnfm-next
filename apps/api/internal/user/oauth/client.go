package oauth

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"vnfm-api/pkg/config"
)

const (
	CodeAccountBanned       = 10014
	CodeRefreshTokenExpired = 10003
	CodeInvalidToken        = 10002
	CodeInvalidGrant        = 15005
	CodeInvalidClientSecret = 15008
)

type Error struct {
	Code       int
	HTTPStatus int
	Message    string
}

func (e *Error) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("oauth: code=%d http=%d msg=%q", e.Code, e.HTTPStatus, e.Message)
	}
	return fmt.Sprintf("oauth: http=%d msg=%q", e.HTTPStatus, e.Message)
}

func IsBanned(err error) bool {
	var oe *Error
	if !stderrors.As(err, &oe) {
		return false
	}
	return oe.Code == CodeAccountBanned || oe.HTTPStatus == http.StatusForbidden
}

func IsRefreshTokenDead(err error) bool {
	var oe *Error
	if !stderrors.As(err, &oe) {
		return false
	}
	switch oe.Code {
	case CodeRefreshTokenExpired, CodeInvalidToken, CodeInvalidGrant, CodeInvalidClientSecret:
		return true
	}
	return false
}

type Client struct {
	cfg        config.OIDCConfig
	httpClient *http.Client
}

func NewClient(cfg config.OIDCConfig) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type UserInfo struct {
	ID        int      `json:"id"`
	Sub       string   `json:"sub"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Picture   string   `json:"picture"`
	Roles     []string `json:"roles"`
	SiteRoles []string `json:"site_roles"`
	UpdatedAt int64    `json:"updated_at"`
}

func (c *Client) ExchangeCode(code, codeVerifier string) (*TokenResponse, error) {
	payload := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  c.cfg.RedirectURI,
		"client_id":     c.cfg.ClientID,
		"client_secret": c.cfg.ClientSecret,
		"code_verifier": codeVerifier,
	}
	data, err := c.postProtocol("/oauth/token", payload)
	if err != nil {
		return nil, err
	}
	var tok TokenResponse
	if jerr := json.Unmarshal(data, &tok); jerr != nil {
		return nil, &Error{Message: "解析 token 响应失败: " + jerr.Error()}
	}
	if tok.AccessToken == "" {
		return nil, &Error{Message: "token 响应缺 access_token"}
	}
	return &tok, nil
}

func (c *Client) FetchUserInfo(accessToken string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", c.cfg.ServerURL+"/oauth/userinfo", nil)
	if err != nil {
		return nil, &Error{Message: "创建 userinfo 请求失败: " + err.Error()}
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &Error{Message: "请求 userinfo 失败: " + err.Error()}
	}
	defer resp.Body.Close()
	data, derr := decodeProtocol(resp)
	if derr != nil {
		return nil, derr
	}
	var info UserInfo
	if jerr := json.Unmarshal(data, &info); jerr != nil {
		return nil, &Error{Message: "解析 userinfo 响应失败: " + jerr.Error()}
	}
	return &info, nil
}

func (c *Client) RefreshOAuthToken(refreshToken string) (*TokenResponse, error) {
	payload := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     c.cfg.ClientID,
		"client_secret": c.cfg.ClientSecret,
	}
	data, err := c.postProtocol("/oauth/token", payload)
	if err != nil {
		return nil, err
	}
	var tok TokenResponse
	if jerr := json.Unmarshal(data, &tok); jerr != nil {
		return nil, &Error{Message: "解析刷新响应失败: " + jerr.Error()}
	}
	if tok.AccessToken == "" {
		return nil, &Error{Message: "刷新响应缺 access_token"}
	}
	return &tok, nil
}

func (c *Client) RevokeToken(refreshToken string) error {
	payload, err := json.Marshal(map[string]string{"token": refreshToken})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.cfg.ServerURL+"/oauth/revoke", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) postProtocol(path string, payload any) (json.RawMessage, error) {
	body, jerr := json.Marshal(payload)
	if jerr != nil {
		return nil, &Error{Message: "序列化请求失败: " + jerr.Error()}
	}
	req, rerr := http.NewRequest("POST", c.cfg.ServerURL+path, bytes.NewReader(body))
	if rerr != nil {
		return nil, &Error{Message: "创建请求失败: " + rerr.Error()}
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &Error{Message: "请求失败: " + err.Error()}
	}
	defer resp.Body.Close()
	return decodeProtocol(resp)
}

func decodeProtocol(resp *http.Response) (json.RawMessage, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &Error{HTTPStatus: resp.StatusCode, Message: "读取响应体失败: " + err.Error()}
	}
	if resp.StatusCode == http.StatusOK {
		return json.RawMessage(body), nil
	}
	var p struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	if jerr := json.Unmarshal(body, &p); jerr != nil {
		return nil, &Error{
			HTTPStatus: resp.StatusCode,
			Message:    fmt.Sprintf("解析响应失败: %v, body=%s", jerr, truncateBody(body)),
		}
	}
	msg := p.ErrorDescription
	if msg == "" {
		msg = p.Error
	}
	return nil, &Error{Code: oauthErrToCode(p.Error), HTTPStatus: resp.StatusCode, Message: msg}
}

func oauthErrToCode(errStr string) int {
	switch errStr {
	case "invalid_token":
		return CodeInvalidToken
	case "invalid_grant", "unauthorized_client":
		return CodeInvalidGrant
	case "invalid_client":
		return CodeInvalidClientSecret
	default:
		return 0
	}
}

func truncateBody(b []byte) string {
	const max = 256
	if len(b) <= max {
		return string(b)
	}
	return string(b[:max]) + "...(truncated)"
}
