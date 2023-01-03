package repositories

import (
	"github.com/yapi-teklif/internal/models"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
)

type IConstructionFeatureRepository interface {
	GetFeatureTypeById(id uint) (*models.CFeature, error)
	GetFeaturesByCategoryID(categoryID uint) ([]models.CFeature, error)
}

type ConstructionFeatureRepository struct {
	Connection database.IConnection
}

func GetConstructionFeatureRepository(c database.IConnection) *ConstructionFeatureRepository {
	return &ConstructionFeatureRepository{
		Connection: c,
	}
}

func (r *ConstructionFeatureRepository) GetFeatureTypeById(id uint) (*models.CFeature, error) {
	var mdl models.CFeature
	err := r.Connection.PsqlDB().Where("id = ?", id).Take(&mdl).Error
	if err != nil {
		return nil, err
	}
	return &mdl, nil
}

func (r *ConstructionFeatureRepository) GetFeaturesByCategoryID(categoryID uint) ([]models.CFeature, error) {
	var features []models.CFeature
	err := r.Connection.PsqlDB().Find(&features).Error
	if err != nil {
		return nil, err
	}
	return features, nil
}
