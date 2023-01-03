package repositories

import (
	"github.com/yapi-teklif/internal/models"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
)

type IProductImageRepository interface {
	Save(mdl *models.ProductImage) error
	FindByRemoteLink(remoteLink string, productID uint) (*models.ProductImage, error)
}

type ProductImageRepository struct {
	Connection database.IConnection
}

func GetProductImageRepository(c database.IConnection) *ProductImageRepository {
	return &ProductImageRepository{
		Connection: c,
	}
}

func (r *ProductImageRepository) Save(mdl *models.ProductImage) error {
	return r.Connection.PsqlDB().Save(mdl).Error
}

func (r *ProductImageRepository) FindByRemoteLink(remoteLink string, productID uint) (*models.ProductImage, error) {
	var image *models.ProductImage
	err := r.Connection.PsqlDB().
		Where("product_id=?", productID).
		Where("remote_link=?", remoteLink).
		Take(image).Error

	if err != nil {
		return nil, err
	}
	return image, err
}
