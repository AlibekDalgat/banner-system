package models

import "time"

type BannerContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

type BannerInfo struct {
	BannerContent
	Id        int       `json:"id"`
	TagsId    []int     `json:"tags_id"`
	FeatureId int       `json:"feature_id"`
	CreatedAt time.Time `json:"created_time"`
	UpdatedAt time.Time `json:"updated_time"`
}
