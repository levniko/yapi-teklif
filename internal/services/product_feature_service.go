package services

import (
	"github.com/yapi-teklif/internal/repositories"
)

type IProductFeatureService interface {
}

type ProductFeatureService struct {
	Repository repositories.IProductFeatureRepository
}

func GetProductFeatureService(r repositories.IProductFeatureRepository) *ProductFeatureService {
	return &ProductFeatureService{Repository: r}
}
