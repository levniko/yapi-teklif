package repositories

import (
	"github.com/yapi-teklif/internal/models"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
)

type IFeatureRepository interface {
	GetFeatureTypeById(id uint) (*models.PFeature, error)
	GetFeaturesByCategoryID(categoryID uint) ([]models.PFeature, error)
}

type FeatureRepository struct {
	Connection database.IConnection
}

func GetFeatureRepository(c database.IConnection) *FeatureRepository {
	return &FeatureRepository{
		Connection: c,
	}
}

func (r *FeatureRepository) GetFeatureTypeById(id uint) (*models.PFeature, error) {
	var mdl models.PFeature
	err := r.Connection.PsqlDB().Where("id = ?", id).Take(&mdl).Error
	if err != nil {
		return nil, err
	}
	return &mdl, nil
}

func (r *FeatureRepository) GetFeaturesByCategoryID(categoryID uint) ([]models.PFeature, error) {
	var features []models.PFeature
	err := r.Connection.PsqlDB().Find(&features).Error
	if err != nil {
		return nil, err
	}
	return features, nil
}
