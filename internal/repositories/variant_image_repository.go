package repositories

import (
	"github.com/yapi-teklif/internal/models"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
)

type IVariantImageRepository interface {
	Save(mdl *models.VariantImage) error
	FindByRemoteLink(remoteLink string, variantID uint) (*models.VariantImage, error)
}

type VariantImageRepository struct {
	Connection database.IConnection
}

func GetVariantImageRepository(c database.IConnection) *VariantImageRepository {
	return &VariantImageRepository{
		Connection: c,
	}
}

func (r *VariantImageRepository) Save(mdl *models.VariantImage) error {
	return r.Connection.PsqlDB().Save(mdl).Error
}

func (r *VariantImageRepository) FindByRemoteLink(remoteLink string, variantID uint) (*models.VariantImage, error) {
	var image *models.VariantImage
	err := r.Connection.PsqlDB().
		Where("variant_id=?", variantID).
		Where("remote_link=?", remoteLink).
		Take(image).Error

	if err != nil {
		return nil, err
	}
	return image, err
}
