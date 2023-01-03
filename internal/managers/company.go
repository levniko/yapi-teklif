package managers

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/yapi-teklif/internal/forms"
	"github.com/yapi-teklif/internal/models"
	"github.com/yapi-teklif/internal/pkg/utils"
	"github.com/yapi-teklif/internal/pkg/utils/argon2"
	"github.com/yapi-teklif/internal/services"
)

type ICompanyManager interface {
	Create(form forms.CompanyCreateForm) (*models.Company, error)
	Login(form forms.LoginForm) (*models.TokenDetails, error)
	CreateToken(companyID uint, companyClaims jwt.MapClaims) (*models.TokenDetails, error)
	CreateAuth(companyID uint, td *models.TokenDetails) error
	DeleteAuth(uuid string) error
	DeleteTokens(authD *models.AccessDetails) error
}

type CompanyManager struct {
	CompanyService services.ICompanyService
}

func GetCompanyManager(companyService services.ICompanyService) *CompanyManager {
	return &CompanyManager{
		CompanyService: companyService,
	}
}

func (m *CompanyManager) Create(form forms.CompanyCreateForm) (*models.Company, error) {
	/*	_, err := m.CompanyService.FindByEmail(strings.ToLower(form.Email))
		if err == nil {
			return nil, errors.New(utils.EmailAlreadyExist)
		}

			_, err = m.CompanyService.FindByPhoneNumber(form.PhoneNumber)
			if err == nil {
				return nil, errors.New(utils.PhoneNumberAlreadyExist)
			}
	*/
	if form.Password != form.PasswordAgain {
		return nil, errors.New(utils.PasswordsAreNotSame)
	}
	password_config := argon2.PasswordConfig{
		Time:    1,
		Memory:  64 * 1024,
		Threads: 4,
		KeyLen:  32,
	}
	var err error
	form.Password, err = argon2.GeneratePassword(&password_config, form.Password)
	if err != nil {
		return nil, errors.New(utils.CompanySaveError)
	}
	model, err := m.CompanyService.Save(form)
	if err != nil {
		return nil, errors.New(utils.CompanySaveError)
	}

	return model, nil
}

func (m *CompanyManager) Login(form forms.LoginForm) (*models.TokenDetails, error) {
	company, err := m.CompanyService.FindByEmail(form.Email)
	if err != nil || company == nil {
		return nil, errors.New(utils.CompanyRecordNotFound)
	}

	is_match, err := argon2.ComparePassword(form.Password, company.PasswordHash)
	if err != nil || !is_match {
		return nil, errors.New(utils.EmailOrPasswordIncorrect)
	}

	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 30).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + strconv.FormatUint(uint64(company.ID), 10)

	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["company_id"] = company.ID
	accessTokenClaims["company_authorized_name"] = company.CompanyAuthorizedName
	accessTokenClaims["company_authorized_surname"] = company.CompanyAuthorizedSurname
	accessTokenClaims["email"] = company.Email
	accessTokenClaims["authorized"] = true
	accessTokenClaims["is_supplier"] = company.IsSupplier
	accessTokenClaims["is_constructor"] = company.IsConstructor
	accessTokenClaims["access_uuid"] = td.AccessUuid
	accessTokenClaims["expiration"] = td.AtExpires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	td.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["company_id"] = company.ID
	refreshTokenClaims["company_authorized_name"] = company.CompanyAuthorizedName
	refreshTokenClaims["company_authorized_surname"] = company.CompanyAuthorizedSurname
	refreshTokenClaims["email"] = company.Email
	refreshTokenClaims["is_supplier"] = company.IsSupplier
	refreshTokenClaims["is_constructor"] = company.IsConstructor
	refreshTokenClaims["refresh_uuid"] = td.RefreshUuid
	refreshTokenClaims["expiration"] = td.RtExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	td.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	err = m.CreateAuth(company.ID, td)
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (m *CompanyManager) CreateToken(companyID uint, companyClaims jwt.MapClaims) (*models.TokenDetails, error) {
	var err error
	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 30).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + strconv.FormatUint(uint64(companyID), 10)

	authorizedName, _ := companyClaims["company_authorized_name"].(string)
	authorizedSurname, _ := companyClaims["company_authorized_surname"].(string)
	email, _ := companyClaims["email"].(string)
	is_supplier, _ := companyClaims["is_supplier"].(bool)
	is_constructor, _ := companyClaims["is_constructor"].(bool)
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["company_id"] = companyID
	accessTokenClaims["company_authorized_name"] = authorizedName
	accessTokenClaims["company_authorized_surname"] = authorizedSurname
	accessTokenClaims["email"] = email
	accessTokenClaims["authorized"] = true
	accessTokenClaims["is_supplier"] = is_supplier
	accessTokenClaims["is_constructor"] = is_constructor
	accessTokenClaims["access_uuid"] = td.AccessUuid
	accessTokenClaims["expiration"] = td.AtExpires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	td.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["company_id"] = companyID
	refreshTokenClaims["company_authorized_name"] = authorizedName
	refreshTokenClaims["company_authorized_surname"] = authorizedSurname
	refreshTokenClaims["email"] = email
	refreshTokenClaims["refresh_uuid"] = td.RefreshUuid
	refreshTokenClaims["expiration"] = td.RtExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	td.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (m *CompanyManager) CreateAuth(companyID uint, td *models.TokenDetails) error {
	return m.CompanyService.CreateAuth(companyID, td)
}

func (m *CompanyManager) DeleteAuth(uuid string) error {
	return m.CompanyService.DeleteAuth(uuid)
}

func (m *CompanyManager) DeleteTokens(authD *models.AccessDetails) error {
	return m.CompanyService.DeleteTokens(authD)
}
