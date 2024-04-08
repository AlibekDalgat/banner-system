package repository

import (
	"banner-system/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	GetUser(login, password string) (models.User, error)
}

type Banner interface {
	GetUserBanner(tagId, featId int) (models.BannerContent, error)
	GetAllBanners(tagId, featId int) ([]models.BannerInfo, error)
	CreateBanner(tagId []int, featId int, input models.BannerContent) (int, error)
	UpdateBanner(banner models.BannerInfo) error
	DeleteBanner(bannerId int) error
}

type Repository struct {
	Authorization
	Banner
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
