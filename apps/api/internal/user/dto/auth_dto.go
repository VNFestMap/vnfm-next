package dto

type OAuthCallbackRequest struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
}

type UserProfile struct {
	ID                  int64  `json:"id"`
	Sub                 string `json:"sub"`
	Name                string `json:"name"`
	Avatar              string `json:"avatar"`
	Email               string `json:"email,omitempty"`
	LanguagePreference  string `json:"language_preference"`
	ThemePreference     string `json:"theme_preference"`
	DisplayMembershipID *int64 `json:"display_membership_id"`
	AccountCenterURL    string `json:"account_center_url"`
}

type SessionResponse struct {
	Token string       `json:"-"`
	User  *UserProfile `json:"user"`
}

type UpdatePrefsRequest struct {
	LanguagePreference  *string `json:"language_preference"`
	ThemePreference     *string `json:"theme_preference"`
	DisplayMembershipID *int64  `json:"display_membership_id"`
	ClearDisplayClub    bool    `json:"clear_display_club"`
}
