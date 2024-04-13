package service

import (
	"banner-system/internal/models"
	"banner-system/internal/repository"
	"time"
)

type BannerService struct {
	repo repository.Banner
}

func (s *BannerService) GetUserBanner(tagId, featId int, accessAdmin bool) (models.BannerContent, error) {
	return s.repo.GetUserBanner(tagId, featId, accessAdmin)
}

func (s *BannerService) GetAllBanners(filter models.Filter) ([]models.BannerInfo, error) {
	return s.repo.GetAllBanners(filter)
}

func (s *BannerService) CreateBanner(input models.BannerInfo) (int, error) {
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()
	return s.repo.CreateBanner(input)
}

func (s *BannerService) UpdateBanner(input models.BannerInfo) error {
	input.UpdatedAt = time.Now()
	return s.repo.UpdateBanner(input)
}

func (s *BannerService) DeleteBanner(id int) (models.CachKey, error) {
	return s.repo.DeleteBanner(id)
}

func NewBannerService(repo repository.Banner) *BannerService {
	return &BannerService{repo: repo}
}
