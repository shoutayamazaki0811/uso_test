package models

import (
	"database/sql/driver"
	"time"
)

type CastRank string

const (
	CastRankStandard CastRank = "standard"
	CastRankPremium  CastRank = "premium"
	CastRankVIP      CastRank = "vip"
)

func (cr CastRank) Value() (driver.Value, error) {
	return string(cr), nil
}

func (cr *CastRank) Scan(value interface{}) error {
	*cr = CastRank(value.(string))
	return nil
}

func (cr CastRank) GetHourlyRate() float64 {
	switch cr {
	case CastRankStandard:
		return 60.0
	case CastRankPremium:
		return 100.0
	case CastRankVIP:
		return 150.0
	default:
		return 60.0
	}
}

type ApprovalStatus string

const (
	ApprovalStatusPending  ApprovalStatus = "pending"
	ApprovalStatusApproved ApprovalStatus = "approved"
	ApprovalStatusRejected ApprovalStatus = "rejected"
)

func (as ApprovalStatus) Value() (driver.Value, error) {
	return string(as), nil
}

func (as *ApprovalStatus) Scan(value interface{}) error {
	*as = ApprovalStatus(value.(string))
	return nil
}

type CastProfile struct {
	ID             int            `json:"id"`
	UserID         int            `json:"user_id"`
	Bio            *string        `json:"bio,omitempty"`
	HourlyRate     float64        `json:"hourly_rate"`
	Rank           CastRank       `json:"rank"`
	ServiceAreas   []string       `json:"service_areas"`
	ApprovalStatus ApprovalStatus `json:"approval_status"`
	ApprovedAt     *time.Time     `json:"approved_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	GalleryImages  []GalleryImage `json:"gallery_images,omitempty"`
}

type GalleryImage struct {
	ID            int       `json:"id"`
	CastProfileID int       `json:"cast_profile_id"`
	ImageURL      string    `json:"image_url"`
	DisplayOrder  int       `json:"display_order"`
	CreatedAt     time.Time `json:"created_at"`
}

type CastProfileCreate struct {
	Bio          *string  `json:"bio"`
	Rank         CastRank `json:"rank" binding:"required,oneof=standard premium vip"`
	ServiceAreas []string `json:"service_areas" binding:"required,min=1"`
}

type CastProfileUpdate struct {
	Bio          *string  `json:"bio"`
	Rank         *CastRank `json:"rank" binding:"omitempty,oneof=standard premium vip"`
	ServiceAreas []string `json:"service_areas"`
}

type CastSearchParams struct {
	Location  string    `form:"location"`
	Date      time.Time `form:"date" time_format:"2006-01-02"`
	StartTime string    `form:"start_time"`
	EndTime   string    `form:"end_time"`
	MinPrice  float64   `form:"min_price"`
	MaxPrice  float64   `form:"max_price"`
	Rank      CastRank  `form:"rank"`
	Page      int       `form:"page,default=1"`
	Limit     int       `form:"limit,default=20"`
}

type ServiceArea struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DisplayOrder int    `json:"display_order"`
}