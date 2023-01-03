package managers

import (
	"errors"

	"github.com/yapi-teklif/internal/forms"
	"github.com/yapi-teklif/internal/models"
	"github.com/yapi-teklif/internal/pkg/utils"
	"github.com/yapi-teklif/internal/services"
)

type IConstructionManager interface {
	Save(form forms.ConstructionCreateForm, companyID uint) (*models.Construction, error, int)
	Update(form forms.ConstructionUpdateForm, constructionID uint, companyID uint) (*models.Construction, error, int)
	Delete(constructionID uint, companyID uint) (error, int)
	FindByIDAndCompanyID(constructionID uint, companyID uint) (*models.Construction, error, int)
	// FindAllByCategoryID(categoryID uint) ([]models.Construction, error, int)
	FetchAuth(authD *models.AccessDetails) (uint, error)
}

type ConstructionManager struct {
	ConstructionService        services.IConstructionService
	ConstructionFeatureService services.IConstructionFeatureService
	ConstructionImageService   services.IConstructionImageService
}

func GetConstructionManager(constructionService services.IConstructionService,
	constructionFeatureService services.IConstructionFeatureService,
	constructionImageService services.IConstructionImageService) *ConstructionManager {
	return &ConstructionManager{
		ConstructionService:        constructionService,
		ConstructionFeatureService: constructionFeatureService,
		ConstructionImageService:   constructionImageService,
	}
}
func (m *ConstructionManager) Save(form forms.ConstructionCreateForm, companyID uint) (*models.Construction, error, int) {
	model, err := m.ConstructionService.Save(form, companyID)
	if err != nil {
		return nil, errors.New(utils.ConstructionCanNotCreated), utils.ConstructionCanNotCreatedCode
	}
	return model, nil, 0
}

func (m *ConstructionManager) Update(form forms.ConstructionUpdateForm, constructionID uint, companyID uint) (*models.Construction, error, int) {
	construction, err := m.ConstructionService.FindByIDAndCompanyID(constructionID, companyID)
	if err != nil {
		return nil, errors.New(utils.ConstructionCanNotFound), utils.ConstructionCanNotFoundCode
	}
	if len(form.ConstructionImages) > 0 {
		err = m.UpdateConstructionImages(construction.ID, form.ConstructionImages)
		if err != nil {
			return nil, errors.New(utils.ConstructionImageCanNotFound), utils.ConstructionImageCanNotFoundCode
		}
	}
	categoryID, err := m.ConstructionService.FindCategoryByID(construction.CategoryID)
	if err != nil {
		return nil, errors.New(utils.ProductCategoryCanNotFound), utils.ProductCategoryCanNotFoundCode
	}
	err = m.ConstructionFeatureService.CheckFeaturesTypes(form.ConstructionFeatures, *categoryID)
	if err != nil {
		return nil, err, utils.VariantFeatureIsWrongTypeCode
	}
	err = m.ConstructionService.UpdateWithFields(construction, form, companyID)
	if err != nil {
		return nil, errors.New(utils.ProductCanNotUpdated), utils.ProductCanNotUpdatedCode
	}
	return construction, nil, 0
}

func (s *ConstructionManager) UpdateConstructionImages(constructionID uint, form []forms.ConstructionImageForm) error {
	for _, value := range form {
		vImage, err := s.ConstructionImageService.FindByRemoteLink(value.RemoteLink, constructionID)
		if err != nil {
			_ = s.ConstructionImageService.Save(&models.ConstructionImage{
				ConstructionID: constructionID,
				RemoteLink:     value.RemoteLink,
			})
			/*
				if err != nil {
					// TODO:log error here
				}*/
		} else {
			vImage.RemoteLink = value.RemoteLink
			_ = s.ConstructionImageService.Save(vImage)
			/*
				if err != nil {
					// TODO:log error here
				}*/
		}
	}
	return nil
}

func (m *ConstructionManager) Delete(constructionID uint, companyID uint) (error, int) {
	err := m.ConstructionService.DeleteByID(constructionID, companyID)
	if err != nil {
		return errors.New(utils.ConstructionImageCanNotDeleted), utils.ConstructionImageCanNotDeletedCode
	}
	return nil, 0
}

func (m *ConstructionManager) FindByIDAndCompanyID(constructionID uint, companyID uint) (*models.Construction, error, int) {
	construction, err := m.ConstructionService.FindByIDAndCompanyID(constructionID, companyID)
	if err != nil {
		return nil, errors.New(utils.ConstructionCanNotFound), utils.ConstructionCanNotFoundCode
	}
	return construction, nil, 0
}

func (m *ConstructionManager) FetchAuth(authD *models.AccessDetails) (uint, error) {
	return m.ConstructionService.FetchAuth(authD)
}
