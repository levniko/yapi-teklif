package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yapi-teklif/internal/forms"
	"github.com/yapi-teklif/internal/models"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
	"github.com/yapi-teklif/internal/repositories"
)

type ICompanyService interface {
	Save(form forms.CompanyCreateForm) (*models.Company, error)
	FindByEmail(email string) (*models.Company, error)
	CreateAuth(companyID uint, td *models.TokenDetails) error
	DeleteAuth(uuid string) error
	DeleteTokens(authD *models.AccessDetails) error
	FetchAuth(authD *models.AccessDetails) (uint, error)
}

type CompanyService struct {
	Repository repositories.ICompanyRepository
	CacheDB    database.ICacheDB
}

func GetCompanyService(r repositories.ICompanyRepository, c database.ICacheDB) *CompanyService {
	return &CompanyService{
		Repository: r,
		CacheDB:    c,
	}
}

func (m *CompanyService) Save(form forms.CompanyCreateForm) (*models.Company, error) {
	model := &models.Company{
		Name:                     form.Name,
		CompanyType:              form.CompanyType,
		Email:                    strings.ToLower(form.Email),
		CompanyAuthorizedName:    form.CompanyAuthorizedName,
		CompanyAuthorizedSurname: form.CompanyAuthorizedSurname,
		Password:                 form.Password,
		PasswordHash:             form.Password,
	}
	if form.WebSite != nil {
		model.WebSite = *form.WebSite
	}

	err := m.Repository.Save(model)

	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *CompanyService) FindByEmail(email string) (*models.Company, error) {
	return s.Repository.FindByEmail(email)
}

func (s *CompanyService) CreateAuth(companyID uint, td *models.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	err := s.CacheDB.Set(td.AccessUuid, strconv.FormatUint(uint64(companyID), 10), at.Sub(now))
	if err != nil {
		return err
	}
	err = s.CacheDB.Set(td.RefreshUuid, strconv.FormatUint(uint64(companyID), 10), rt.Sub(now))
	if err != nil {
		return err
	}
	return nil
}

func (s *CompanyService) DeleteAuth(uuid string) error {
	deletedRt, err := s.CacheDB.Del(uuid)
	if err != nil {
		return err
	}
	if deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

func (s *CompanyService) DeleteTokens(authD *models.AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.AccessUuid, authD.CompanyID)
	//delete access token
	deletedAt, err := s.CacheDB.Del(authD.AccessUuid)
	if err != nil {
		return err
	}
	if deletedAt != 1 {
		return errors.New("something went wrong")
	}
	//delete refresh token
	deletedRt, err := s.CacheDB.Del(refreshUuid)
	if err != nil {
		return err
	}
	if deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

func (s *CompanyService) FetchAuth(authD *models.AccessDetails) (uint, error) {
	cID, err := s.CacheDB.Get(authD.AccessUuid)
	if err != nil {
		return 0, err
	}
	companyID, _ := strconv.ParseUint(cID, 10, 64)
	return uint(companyID), nil
}
