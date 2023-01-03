package managers

import (
	"errors"

	"github.com/yapi-teklif/internal/forms"
	"github.com/yapi-teklif/internal/models"
	"github.com/yapi-teklif/internal/pkg/utils"
	"github.com/yapi-teklif/internal/services"
)

type IVariantManager interface {
	Save(form forms.VariantCreateForm, companyID uint) (*models.Variant, error, int)
	Update(form forms.VariantUpdateForm, variantID uint, companyID uint) (*models.Variant, error, int)
	Delete(variantID uint, companyID uint) (error, int)
	FindByID(id uint) (*models.Variant, error)
	FetchAuth(authD *models.AccessDetails) (uint, error)
}

type VariantManager struct {
	VariantService      services.IVariantService
	FeatureService      services.IFeatureService
	ProductService      services.IProductService
	VariantImageService services.IVariantImageService
}

func GetVariantManager(variantService services.IVariantService,
	featureService services.IFeatureService,
	productService services.IProductService,
	variantImageService services.IVariantImageService) *VariantManager {
	return &VariantManager{
		VariantService:      variantService,
		FeatureService:      featureService,
		ProductService:      productService,
		VariantImageService: variantImageService,
	}
}

func (m *VariantManager) Save(form forms.VariantCreateForm, companyID uint) (*models.Variant, error, int) {
	categoryID, err := m.ProductService.FindCategoryByID(uint(form.ProductId))
	if err != nil {
		return nil, errors.New(utils.ProductCategoryCanNotFound), utils.ProductCategoryCanNotFoundCode
	}
	err = m.FeatureService.CheckFeaturesTypes(form.Features, *categoryID)
	if err != nil {
		return nil, err, utils.VariantFeatureIsWrongTypeCode
	}
	variant, err := m.VariantService.Save(form)
	if err != nil {
		return nil, errors.New(utils.VariantCanNotCreated), utils.VariantCanNotCreatedCode
	}
	return variant, nil, 0
}

func (m *VariantManager) Update(form forms.VariantUpdateForm, variantID uint, companyID uint) (*models.Variant, error, int) {
	variant, err := m.VariantService.FindByID(variantID)
	if err != nil {
		return nil, errors.New(utils.VariantCanNotFound), utils.VariantCanNotFoundCode
	}
	categoryID, err := m.ProductService.FindCategoryByID(variant.ProductID)
	if err != nil {
		return nil, errors.New(utils.ProductCategoryCanNotFound), utils.ProductCategoryCanNotFoundCode
	}
	err = m.FeatureService.CheckFeaturesTypes(form.Features, *categoryID)
	if err != nil {
		return nil, err, utils.VariantFeatureIsWrongTypeCode
	}
	if len(form.VariantImages) > 0 {
		err = m.UpdateVariantImages(variant.ID, form.VariantImages)
		if err != nil {
			return nil, errors.New(utils.VariantImageCanNotUpdated), utils.VariantImageCanNotUpdatedCode
		}
	}
	err = m.VariantService.UpdateWithFields(variant, form, companyID)
	if err != nil {
		return nil, errors.New(utils.VariantCanNotUpdated), utils.VariantCanNotUpdatedCode
	}
	return variant, nil, 0
}

func (s *VariantManager) UpdateVariantImages(productID uint, form []forms.VariantImageForm) error {
	for _, value := range form {
		vImage, err := s.VariantImageService.FindByRemoteLink(value.RemoteLink, productID)
		if err != nil {
			_ = s.VariantImageService.Save(&models.VariantImage{
				VariantID:  productID,
				RemoteLink: value.RemoteLink,
			})
			/*
				if err != nil {
					// TODO:log error here
				}*/
		} else {
			vImage.RemoteLink = value.RemoteLink
			_ = s.VariantImageService.Save(vImage)
			/*
				if err != nil {
					// TODO:log error here
				}*/
		}
	}
	return nil
}

func (m *VariantManager) Delete(variantID uint, companyID uint) (error, int) {
	err := m.VariantService.DeleteByIDAndCompanyID(variantID, companyID)
	if err != nil {
		return errors.New(utils.VariantCanNotDeleted), utils.VariantCanNotDeletedCode
	}
	return nil, 0
}

func (m *VariantManager) FindByID(id uint) (*models.Variant, error) {
	variant, err := m.VariantService.FindByID(id)
	if err != nil {
		return nil, errors.New(utils.VariantCanNotFound)
	}
	return variant, nil
}

func (m *VariantManager) FetchAuth(authD *models.AccessDetails) (uint, error) {
	return m.VariantService.FetchAuth(authD)
}
