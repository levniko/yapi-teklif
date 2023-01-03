package repositories

import (
	"github.com/yapi-teklif/internal/models"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
)

type IConstructionImageRepository interface {
	Save(mdl *models.ConstructionImage) error
	FindByRemoteLink(remoteLink string, constructionID uint) (*models.ConstructionImage, error)
}

type ConstructionImageRepository struct {
	Connection database.IConnection
}

func GetConstructionImageRepository(c database.IConnection) *ConstructionImageRepository {
	return &ConstructionImageRepository{
		Connection: c,
	}
}

func (r *ConstructionImageRepository) Save(mdl *models.ConstructionImage) error {
	return r.Connection.PsqlDB().Save(mdl).Error
}

func (r *ConstructionImageRepository) FindByRemoteLink(remoteLink string, constructionID uint) (*models.ConstructionImage, error) {
	var image *models.ConstructionImage
	err := r.Connection.PsqlDB().
		Where("construction_id=?", constructionID).
		Where("remote_link=?", remoteLink).
		Take(image).Error

	if err != nil {
		return nil, err
	}
	return image, err
}
