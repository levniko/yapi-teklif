package managers

import (
	"errors"

	"github.com/yapi-teklif/internal/models"
	"github.com/yapi-teklif/internal/pkg/utils"
	"github.com/yapi-teklif/internal/services"
)

type IFeatureManager interface {
	GetFeaturesByCategoryID(categoryID uint) ([]models.PFeature, error, int)
}

type FeatureManager struct {
	FeatureService services.IFeatureService
}

func GetFeatureManager(featureService services.IFeatureService) *FeatureManager {
	return &FeatureManager{
		FeatureService: featureService,
	}
}

func (m *FeatureManager) GetFeaturesByCategoryID(categoryID uint) ([]models.PFeature, error, int) {
	features, err := m.FeatureService.GetFeaturesByCategoryID(categoryID)
	if err != nil {
		return nil, errors.New(utils.FeatureNotFound), utils.FeatureNotFoundCode
	}
	return features, nil, 0
}
