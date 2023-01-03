package services

import (
	"github.com/yapi-teklif/internal/models"
	"github.com/yapi-teklif/internal/repositories"
)

type IVariantImageService interface {
	Save(mdl *models.VariantImage) error
	FindByRemoteLink(remoteLink string, variantID uint) (*models.VariantImage, error)
}

type VariantImageService struct {
	Repository repositories.IVariantImageRepository
}

func GetVariantImageService(r repositories.IVariantImageRepository) *VariantImageService {
	return &VariantImageService{Repository: r}
}

func (s *VariantImageService) FindByRemoteLink(remoteLink string, variantID uint) (*models.VariantImage, error) {
	return s.Repository.FindByRemoteLink(remoteLink, variantID)
}

func (s *VariantImageService) Save(mdl *models.VariantImage) error {
	return s.Repository.Save(mdl)
}
