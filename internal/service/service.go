package service

import (
	"banner-system/internal/models"
	"banner-system/internal/repository"
)

type Authorization interface {
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (bool, error)
}

type Banner interface {
	GetUserBanner(tagId, featId int, accessAdmin bool) (models.BannerContent, error)
	GetAllBanners(filter models.Filter) ([]models.BannerInfo, error)
	CreateBanner(input models.BannerInfo) (int, error)
	UpdateBanner(banner models.BannerInfo) error
	DeleteBanner(bannerId int) (models.CachKey, error)
}

type Service struct {
	Authorization
	Banner
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo),
		Banner:        NewBannerService(repo),
	}
}
