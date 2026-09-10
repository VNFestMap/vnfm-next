package repository

import (
	"vnfm-api/internal/club/model"

	"gorm.io/gorm"
)

type ClubRepository struct {
	db *gorm.DB
}

func NewClubRepository(db *gorm.DB) *ClubRepository {
	return &ClubRepository{db: db}
}

func (r *ClubRepository) DB() *gorm.DB { return r.db }

func (r *ClubRepository) List(q, country, province, typ string, limit, offset int) ([]model.Club, int64, error) {
	query := r.db.Model(&model.Club{})
	if q != "" {
		like := "%" + q + "%"
		query = query.Where("name ILIKE ? OR school ILIKE ? OR province ILIKE ? OR prefecture ILIKE ?", like, like, like, like)
	}
	if country != "" && country != "all" {
		query = query.Where("country = ?", country)
	}
	if province != "" {
		query = query.Where("province = ? OR prefecture = ?", province, province)
	}
	if typ != "" && typ != "all" {
		query = query.Where("type = ?", typ)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	var rows []model.Club
	err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&rows).Error
	return rows, total, err
}

type RegionCountRow struct {
	Key   string
	Count int64
}

func (r *ClubRepository) RegionCounts(country string) ([]RegionCountRow, error) {
	col := "province"
	if country == "japan" {
		col = "prefecture"
	}
	var rows []RegionCountRow
	err := r.db.Model(&model.Club{}).
		Select(col+" AS key, COUNT(*) AS count").
		Where("country = ?", country).
		Where(col + " <> ''").
		Group(col).
		Order("count DESC, key ASC").
		Scan(&rows).Error
	return rows, err
}

func (r *ClubRepository) FindByID(id int64) (*model.Club, error) {
	var row model.Club
	if err := r.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *ClubRepository) Create(club *model.Club) error {
	return r.db.Create(club).Error
}

func (r *ClubRepository) Update(club *model.Club) error {
	return r.db.Save(club).Error
}

func (r *ClubRepository) FindMembership(userID, clubID int64) (*model.Membership, error) {
	var row model.Membership
	err := r.db.Where("user_id = ? AND club_id = ?", userID, clubID).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *ClubRepository) CreateMembership(m *model.Membership) error {
	return r.db.Create(m).Error
}

func (r *ClubRepository) SaveMembership(m *model.Membership) error {
	return r.db.Save(m).Error
}

func (r *ClubRepository) ListMembers(clubID int64, status string) ([]model.Membership, error) {
	q := r.db.Where("club_id = ?", clubID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var rows []model.Membership
	err := q.Order("joined_at ASC").Find(&rows).Error
	return rows, err
}

func (r *ClubRepository) ListUserMemberships(userID int64) ([]model.Membership, error) {
	var rows []model.Membership
	err := r.db.Where("user_id = ? AND status IN ?", userID, []string{"active", "pending"}).
		Order("joined_at DESC").Find(&rows).Error
	return rows, err
}

func (r *ClubRepository) ListPendingForManager(userID int64) ([]model.Membership, error) {
	var managed []model.Membership
	if err := r.db.Where("user_id = ? AND status = ? AND role IN ?", userID, "active", []string{"manager", "representative"}).
		Find(&managed).Error; err != nil {
		return nil, err
	}
	if len(managed) == 0 {
		return nil, nil
	}
	ids := make([]int64, 0, len(managed))
	for _, m := range managed {
		ids = append(ids, m.ClubID)
	}
	var rows []model.Membership
	err := r.db.Where("club_id IN ? AND status = ?", ids, "pending").Order("joined_at ASC").Find(&rows).Error
	return rows, err
}

func (r *ClubRepository) FindMembershipByID(id int64) (*model.Membership, error) {
	var row model.Membership
	if err := r.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *ClubRepository) Managers(clubID int64) ([]model.Membership, error) {
	var rows []model.Membership
	err := r.db.Where("club_id = ? AND status = ? AND role IN ?", clubID, "active", []string{"manager", "representative"}).
		Find(&rows).Error
	return rows, err
}

func (r *ClubRepository) CreateCode(c *model.VerificationCode) error {
	return r.db.Create(c).Error
}

func (r *ClubRepository) ListCodes(clubID int64) ([]model.VerificationCode, error) {
	var rows []model.VerificationCode
	err := r.db.Where("club_id = ?", clubID).Order("id DESC").Find(&rows).Error
	return rows, err
}

func (r *ClubRepository) FindCode(code string) (*model.VerificationCode, error) {
	var row model.VerificationCode
	if err := r.db.Where("code = ?", code).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *ClubRepository) SaveCode(c *model.VerificationCode) error {
	return r.db.Save(c).Error
}

func (r *ClubRepository) FindCodeByID(id int64) (*model.VerificationCode, error) {
	var row model.VerificationCode
	if err := r.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *ClubRepository) UserNames(ids []int64) map[int64]string {
	out := map[int64]string{}
	if len(ids) == 0 {
		return out
	}
	type row struct {
		ID   int64
		Name string
	}
	var rows []row
	r.db.Table("users").Select("id, name").Where("id IN ?", ids).Scan(&rows)
	for _, x := range rows {
		out[x.ID] = x.Name
	}
	return out
}

func (r *ClubRepository) ClubNames(ids []int64) map[int64]string {
	out := map[int64]string{}
	if len(ids) == 0 {
		return out
	}
	var clubs []model.Club
	r.db.Select("id, name").Where("id IN ?", ids).Find(&clubs)
	for _, c := range clubs {
		out[c.ID] = c.Name
	}
	return out
}
