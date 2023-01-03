package services

import (
	"github.com/yapi-teklif/internal/models"
	"github.com/yapi-teklif/internal/repositories"
)

type IConstructionImageService interface {
	Save(mdl *models.ConstructionImage) error
	FindByRemoteLink(remoteLink string, constructionID uint) (*models.ConstructionImage, error)
}

type ConstructionImageService struct {
	Repository repositories.IConstructionImageRepository
}

func GetConstructionImageService(r repositories.IConstructionImageRepository) *ConstructionImageService {
	return &ConstructionImageService{Repository: r}
}

func (s *ConstructionImageService) FindByRemoteLink(remoteLink string, constructionID uint) (*models.ConstructionImage, error) {
	return s.Repository.FindByRemoteLink(remoteLink, constructionID)
}

func (s *ConstructionImageService) Save(mdl *models.ConstructionImage) error {
	return s.Repository.Save(mdl)
}
