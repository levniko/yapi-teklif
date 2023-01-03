package services

import (
	"github.com/yapi-teklif/internal/models"
	"github.com/yapi-teklif/internal/repositories"
)

type IProductImageService interface {
	Save(mdl *models.ProductImage) error
	FindByRemoteLink(remoteLink string, productID uint) (*models.ProductImage, error)
}

type ProductImageService struct {
	Repository repositories.IProductImageRepository
}

func GetProductImageService(r repositories.IProductImageRepository) *ProductImageService {
	return &ProductImageService{Repository: r}
}

func (s *ProductImageService) FindByRemoteLink(remoteLink string, productID uint) (*models.ProductImage, error) {
	return s.Repository.FindByRemoteLink(remoteLink, productID)
}

func (s *ProductImageService) Save(mdl *models.ProductImage) error {
	return s.Repository.Save(mdl)
}
