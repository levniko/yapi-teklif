package repositories

import (
	"github.com/yapi-teklif/internal/models"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
	"gorm.io/gorm/clause"
)

type IVariantRepository interface {
	Save(mdl *models.Variant) error
	UpdateWithFields(variant *models.Variant, fields map[string]interface{}) error
	DeleteByIDAndCompanyID(variantID uint, companyID uint) error
	FindByID(id uint) (*models.Variant, error)
}

type VariantRepository struct {
	Connection database.IConnection
}

func GetVariantRepository(c database.IConnection) *VariantRepository {
	return &VariantRepository{
		Connection: c,
	}
}

func (r *VariantRepository) Save(mdl *models.Variant) error {
	return r.Connection.PsqlDB().Save(mdl).Error
}

func (r *VariantRepository) UpdateWithFields(variant *models.Variant, fields map[string]interface{}) error {
	return r.Connection.PsqlDB().Model(variant).Updates(fields).Error
}

func (r *VariantRepository) DeleteByIDAndCompanyID(variantID uint, companyID uint) error {
	return r.Connection.PsqlDB().Where("id =?", variantID).
		Select(clause.Associations).
		Where("company_id =?", companyID).
		Delete(&models.Variant{}).Error
}

func (r *VariantRepository) FindByID(id uint) (*models.Variant, error) {
	var variant models.Variant
	err := r.Connection.PsqlDB().Preload("Features").Preload("Features.ProductFeature").Take(&variant).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}
