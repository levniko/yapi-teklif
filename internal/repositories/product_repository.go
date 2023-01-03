package repositories

import (
	"github.com/yapi-teklif/internal/models"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
	"gorm.io/gorm/clause"
)

type IProductRepository interface {
	Save(mdl *models.Product) error
	UpdateWithFields(product *models.Product, fields interface{}) error
	DeleteByID(productID uint, companyID uint) error
	FindCategoryByID(id uint) (*uint, error)
	FindByID(productID uint) (*models.Product, error)
	CountBySpuAndCompanyID(spu string, companyID uint) (int64, error)
	FindAllByCategoryID(categoryID uint) ([]models.Product, error)
}

type ProductRepository struct {
	Connection database.IConnection
}

func GetProductRepository(c database.IConnection) *ProductRepository {
	return &ProductRepository{
		Connection: c,
	}
}

func (r *ProductRepository) Save(mdl *models.Product) error {
	return r.Connection.PsqlDB().Save(mdl).Error
}

func (r ProductRepository) UpdateWithFields(product *models.Product, fields interface{}) error {
	return r.Connection.PsqlDB().Model(product).Omit(clause.Associations).Updates(fields).Error
}

func (r ProductRepository) DeleteByID(productID uint, companyID uint) error {
	return r.Connection.PsqlDB().
		Select(clause.Associations).
		Where("id =?", productID).
		Delete(&models.Product{}).Error
}

func (r *ProductRepository) FindCategoryByID(id uint) (*uint, error) {
	var categoryID uint
	err := r.Connection.PsqlDB().Model(&models.Product{}).Select("category_id").Where("id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &categoryID, nil
}

func (r *ProductRepository) FindByID(productID uint) (*models.Product, error) {
	var product models.Product
	err := r.Connection.PsqlDB().
		Where("id = ?", productID).
		Preload("ProductImages").
		Preload("ProductCategory").
		Take(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) CountBySpuAndCompanyID(spu string, companyID uint) (int64, error) {
	var count int64
	err := r.Connection.PsqlDB().
		Model(&models.Product{}).
		Where("spu =?", spu).
		Where("company_id =?", companyID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ProductRepository) FindAllByCategoryID(categoryID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.Connection.PsqlDB().
		Where("product_category_id IN (?)", r.Connection.PsqlDB().
			Table("product_categories").
			Select("id").
			Where("parent_id IN (?)", r.Connection.PsqlDB().
				Table("product_categories").
				Select("id").
				Where("parent_id = ?", categoryID))).
		Find(&products).Error

	if err != nil {
		return nil, err
	}
	
	return products, nil
}
