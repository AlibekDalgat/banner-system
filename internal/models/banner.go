package models

import (
	"database/sql/driver"
	"github.com/lib/pq"
	"time"
)

type BannerContent struct {
	Title *string `json:"title" db:"title"`
	Text  *string `json:"text" db:"text"`
	URL   *string `json:"url" db:"url"`
}

type IntArray []int

func (a *IntArray) Scan(value interface{}) error {
	return pq.Array(a).Scan(value)
}

func (a IntArray) Value() (driver.Value, error) {
	return pq.Array(a).Value()
}

type BannerInfo struct {
	Id        int            `json:"id" db:"id"`
	Content   *BannerContent `json:"content" db:"content"`
	TagIds    *IntArray      `json:"tag_ids" db:"tag_ids"`
	FeatureId *int           `json:"feature_id" db:"feature_id"`
	IsActive  *bool          `json:"is_active" db:"is_active"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

func (b BannerInfo) Validate() bool {
	if b.FeatureId == nil {
		return false
	}
	if b.TagIds == nil {
		return false
	} else if len(*b.TagIds) == 0 {
		return false
	}
	if b.IsActive == nil {
		return false
	}
	if b.Content == nil {
		return false
	} else if b.Content.Title == nil || b.Content.Text == nil || b.Content.URL == nil {
		return false
	}
	return true
}

type Filter struct {
	TagId     int
	FeatureId int
	Limit     int
	Offset    int
}

type Version struct {
	Number int `json:"version"`
}
