package repository

import (
	"banner-system/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	GetUser(login, password string) (models.User, bool, error)
}

type Banner interface {
	GetUserBanner(tagId, featId int, accessAdmin bool) (models.BannerContent, error)
	GetAllBanners(filter models.Filter) ([]models.BannerInfo, error)
	CreateBanner(input models.BannerInfo) (int, error)
	UpdateBanner(banner models.BannerInfo) error
	DeleteBanner(bannerId int) error
	DeleteBannerByTagFeat(tagId int, featId int) error
	GetOldVersionBanner(id int) ([]models.BannerInfo, error)
	BackToVersion(id, version int) error
}

type Repository struct {
	Authorization
	Banner
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Banner:        NewBannerPostgres(db),
	}
}
