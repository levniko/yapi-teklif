package repositories

import (
	"github.com/yapi-teklif/internal/models"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
	"gorm.io/gorm/clause"
)

type ICompanyRepository interface {
	Save(mdl *models.Company) error
	FindByEmail(email string) (*models.Company, error)
}

type CompanyRepository struct {
	Connection database.IConnection
}

func GetCompanyRepository(c database.IConnection) *CompanyRepository {
	return &CompanyRepository{
		Connection: c,
	}
}

func (r *CompanyRepository) Save(mdl *models.Company) error {
	return r.Connection.PsqlDB().Omit(clause.Associations).Save(mdl).Error
}

func (r *CompanyRepository) FindByEmail(email string) (*models.Company, error) {
	var company models.Company
	err := r.Connection.PsqlDB().Where("email = ?", email).Take(&company).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}
