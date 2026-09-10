package dto

type EventView struct {
	ID            int64  `json:"id"`
	ClubID        int64  `json:"club_id"`
	ClubName      string `json:"club_name"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Location      string `json:"location"`
	CoverKey      string `json:"cover_key"`
	StartsAt      string `json:"starts_at"`
	EndsAt        string `json:"ends_at,omitempty"`
	RegisterUntil string `json:"register_until,omitempty"`
	Registered    bool   `json:"registered"`
	Locked        bool   `json:"locked"`
	CanManage     bool   `json:"can_manage"`
}

type UpsertEventRequest struct {
	ClubID        int64  `json:"club_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Location      string `json:"location"`
	StartsAt      string `json:"starts_at"`
	EndsAt        string `json:"ends_at"`
	RegisterUntil string `json:"register_until"`
}
