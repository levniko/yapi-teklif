package services

import (
	"errors"
	"strconv"

	"github.com/yapi-teklif/internal/forms"
	"github.com/yapi-teklif/internal/models"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
	"github.com/yapi-teklif/internal/repositories"
)

type IVariantService interface {
	Save(form forms.VariantCreateForm) (*models.Variant, error)
	UpdateWithFields(variant *models.Variant, form forms.VariantUpdateForm, companyID uint) error
	DeleteByIDAndCompanyID(variantID uint, companyID uint) error
	FindByID(id uint) (*models.Variant, error)
	FetchAuth(authD *models.AccessDetails) (uint, error)
}

type VariantService struct {
	Repository repositories.IVariantRepository
	CacheDB    database.ICacheDB
}

func GetVariantService(r repositories.IVariantRepository, c database.ICacheDB) *VariantService {
	return &VariantService{
		Repository: r,
		CacheDB:    c,
	}
}

func (s *VariantService) Save(form forms.VariantCreateForm) (*models.Variant, error) {
	features := []models.ProductFeature{}
	for _, f := range form.Features {
		feat := models.ProductFeature{
			ProductFeatureID: f.FeatureID,
			Value:            f.Value,
		}
		features = append(features, feat)
	}
	variantImages := []models.VariantImage{}
	if form.VariantImages != nil {
		for _, s := range form.VariantImages {
			image := models.VariantImage{
				RemoteLink: s.RemoteLink,
			}
			variantImages = append(variantImages, image)
		}
	}
	variant := &models.Variant{
		ProductID: uint(form.ProductId),
		SKU:       form.SKU,
		Features:  features,
	}
	if len(variantImages) > 0 {
		variant.VariantImages = variantImages
	}
	err := s.Repository.Save(variant)
	if err != nil {
		return nil, err
	}
	return variant, nil
}

func (s VariantService) UpdateWithFields(variant *models.Variant, form forms.VariantUpdateForm, companyID uint) error {
	var fields = map[string]interface{}{
		"name": form.Name,
		"sku":  form.SKU, //TODO: check SKU with companyID
	}
	if form.Description != nil {
		fields["description"] = *form.Description
	}
	if form.IsActive != nil {
		fields["is_active"] = *form.IsActive
	}
	return s.Repository.UpdateWithFields(variant, fields)
}

func (s *VariantService) DeleteByIDAndCompanyID(variantID uint, companyID uint) error {
	return s.Repository.DeleteByIDAndCompanyID(variantID, companyID)
}

func (s *VariantService) FindByID(id uint) (*models.Variant, error) {
	return s.Repository.FindByID(id)
}

func (s *VariantService) FetchAuth(authD *models.AccessDetails) (uint, error) {
	cID, err := s.CacheDB.Get(authD.AccessUuid)
	if err != nil {
		return 0, err
	}
	companyID, _ := strconv.ParseUint(cID, 10, 64)
	if authD.CompanyID != uint(companyID) {
		return 0, errors.New("unauthorized")
	}
	return uint(companyID), nil
}
