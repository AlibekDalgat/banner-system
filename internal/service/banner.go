package service

import (
	"banner-system/internal/models"
	"banner-system/internal/repository"
	"errors"
	"time"
)

type BannerService struct {
	repo repository.Banner
}

func NewBannerService(repo repository.Banner) *BannerService {
	return &BannerService{repo: repo}
}

func (s *BannerService) GetUserBanner(tagId, featId int, accessAdmin bool) (models.BannerContent, error) {
	return s.repo.GetUserBanner(tagId, featId, accessAdmin)
}

func (s *BannerService) GetAllBanners(filter models.Filter) ([]models.BannerInfo, error) {
	return s.repo.GetAllBanners(filter)
}

func (s *BannerService) CreateBanner(input models.BannerInfo) (int, error) {
	if ok := input.Validate(); !ok {
		return 0, errors.New("Неполная информация о баннере")
	}
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()
	return s.repo.CreateBanner(input)
}

func (s *BannerService) UpdateBanner(input models.BannerInfo) error {
	input.UpdatedAt = time.Now()
	return s.repo.UpdateBanner(input)
}

func (s *BannerService) DeleteBanner(id int) error {
	return s.repo.DeleteBanner(id)
}

func (s *BannerService) DeleteBannerByTagFeat(tagId, featureId int) error {
	return s.repo.DeleteBannerByTagFeat(tagId, featureId)
}

func (s *BannerService) GetOldVersionBanner(id int) ([]models.BannerInfo, error) {
	return s.repo.GetOldVersionBanner(id)
}

func (s *BannerService) BackToVersion(id, version int) error {
	if version < 1 || version > 3 {
		return errors.New("Неверный номер версии")
	}
	return s.repo.BackToVersion(id, version)
}
