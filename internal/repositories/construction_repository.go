package repositories

import (
	"github.com/yapi-teklif/internal/models"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
	"gorm.io/gorm/clause"
)

type IConstructionRepository interface {
	Save(mdl *models.Construction) error
	UpdateWithFields(construction *models.Construction, fields interface{}) error
	DeleteByID(constructionID uint, companyID uint) error
	FindCategoryByID(id uint) (*uint, error)
	FindByIDAndCompanyID(constructionID uint, companyID uint) (*models.Construction, error)
	FindAllByCategoryID(categoryID uint) ([]models.Construction, error)
}

type ConstructionRepository struct {
	Connection database.IConnection
}

func GetConstructionRepository(c database.IConnection) *ConstructionRepository {
	return &ConstructionRepository{
		Connection: c,
	}
}

func (r *ConstructionRepository) Save(mdl *models.Construction) error {
	return r.Connection.PsqlDB().Save(mdl).Error
}

func (r ConstructionRepository) UpdateWithFields(construction *models.Construction, fields interface{}) error {
	return r.Connection.PsqlDB().Model(construction).Updates(fields).Error
}

func (r ConstructionRepository) DeleteByID(constructionID uint, companyID uint) error {
	return r.Connection.PsqlDB().
		Select(clause.Associations).
		Where("id =?", constructionID).
		Delete(&models.Construction{}).Error
}

func (r *ConstructionRepository) FindCategoryByID(id uint) (*uint, error) {
	var categoryID uint
	err := r.Connection.PsqlDB().Model(&models.Construction{}).Select("category_id").Where("id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &categoryID, nil
}

func (r ConstructionRepository) FindByIDAndCompanyID(constructionID uint, companyID uint) (*models.Construction, error) {
	var construction models.Construction
	err := r.Connection.PsqlDB().
		Where("id =?", constructionID).
		Where("company_id = ?", companyID).
		First(&construction).Error
	if err != nil {
		return nil, err
	}
	return &construction, nil
}

func (r *ConstructionRepository) FindAllByCategoryID(categoryID uint) ([]models.Construction, error) {
	var constructions []models.Construction
	err := r.Connection.PsqlDB().
		Where("construction_category_id IN (?)", r.Connection.PsqlDB().
			Table("construction_categories").
			Select("id").
			Where("parent_id IN (?)", r.Connection.PsqlDB().
				Table("construction_categories").
				Select("id").
				Where("parent_id = ?", categoryID))).
		Find(&constructions).Error

	if err != nil {
		return nil, err
	}

	return constructions, nil
}
