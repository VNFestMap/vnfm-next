package dto

type ClubSummary struct {
	ID            int64  `json:"id"`
	Country       string `json:"country"`
	Name          string `json:"name"`
	School        string `json:"school"`
	Province      string `json:"province"`
	Prefecture    string `json:"prefecture"`
	City          string `json:"city"`
	Type          string `json:"type"`
	LogoKey       string `json:"logo_key"`
	ContactHidden bool   `json:"contact_hidden"`
	Info          string `json:"info,omitempty"`
	Remark        string `json:"remark,omitempty"`
	ExternalLinks string `json:"external_links,omitempty"`
	CanApply      bool   `json:"can_apply"`
	MyRole        string `json:"my_role,omitempty"`
	MyStatus      string `json:"my_status,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
}

type RegionCount struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

type CreateClubRequest struct {
	Name          string `json:"name"`
	Country       string `json:"country"`
	School        string `json:"school"`
	Province      string `json:"province"`
	Prefecture    string `json:"prefecture"`
	City          string `json:"city"`
	Type          string `json:"type"`
	Info          string `json:"info"`
	Remark        string `json:"remark"`
	ContactHidden *bool  `json:"contact_hidden"`
}

type UpdateClubRequest struct {
	Name          *string `json:"name"`
	School        *string `json:"school"`
	Province      *string `json:"province"`
	Prefecture    *string `json:"prefecture"`
	City          *string `json:"city"`
	Type          *string `json:"type"`
	Info          *string `json:"info"`
	Remark        *string `json:"remark"`
	ContactHidden *bool   `json:"contact_hidden"`
}

type ApplyRequest struct {
	ContactAccount string `json:"contact_account"`
	ApplyReason    string `json:"apply_reason"`
	ApplyRole      string `json:"apply_role"`
}

type MembershipView struct {
	ID             int64  `json:"id"`
	UserID         int64  `json:"user_id"`
	UserName       string `json:"user_name"`
	ClubID         int64  `json:"club_id"`
	ClubName       string `json:"club_name"`
	Country        string `json:"country"`
	Role           string `json:"role"`
	Status         string `json:"status"`
	ContactAccount string `json:"contact_account,omitempty"`
	ApplyReason    string `json:"apply_reason,omitempty"`
	ApplyRole      string `json:"apply_role,omitempty"`
}

type ChangeRoleRequest struct {
	Role string `json:"role"`
}

type CreateCodeRequest struct {
	MaxUses int `json:"max_uses"`
	Hours   int `json:"hours"`
}

type CodeView struct {
	ID        int64  `json:"id"`
	Code      string `json:"code"`
	MaxUses   int    `json:"max_uses"`
	UseCount  int    `json:"use_count"`
	IsActive  bool   `json:"is_active"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

type RedeemRequest struct {
	Code string `json:"code"`
}
