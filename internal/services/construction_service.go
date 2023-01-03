package services

import (
	"errors"
	"strconv"

	"github.com/shopspring/decimal"
	"github.com/yapi-teklif/internal/forms"
	"github.com/yapi-teklif/internal/models"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
	"github.com/yapi-teklif/internal/repositories"
)

type IConstructionService interface {
	Save(form forms.ConstructionCreateForm, companyID uint) (*models.Construction, error)
	UpdateWithFields(construction *models.Construction, form forms.ConstructionUpdateForm, companyID uint) error
	DeleteByID(constructionID uint, companyID uint) error
	FindCategoryByID(id uint) (*uint, error)
	FindByIDAndCompanyID(constructionID uint, companyID uint) (*models.Construction, error)
	FindAllByCategoryID(categoryID uint) ([]models.Construction, error)
	FetchAuth(authD *models.AccessDetails) (uint, error)
}

type ConstructionService struct {
	Repository repositories.IConstructionRepository
	CacheDB    database.ICacheDB
}

func GetConstructionService(r repositories.IConstructionRepository, c database.ICacheDB) *ConstructionService {
	return &ConstructionService{
		Repository: r,
		CacheDB:    c,
	}
}
func (s *ConstructionService) Save(form forms.ConstructionCreateForm, companyID uint) (*models.Construction, error) {
	costOfProject := decimal.NewFromFloat(form.CostOfProject)
	landArea := decimal.NewFromFloat(form.LandArea)
	constructionZone := decimal.NewFromFloat(form.ConstructionZone)
	construction := &models.Construction{
		CompanyID:        companyID,
		Name:             form.Name,
		CategoryID:       form.ConstructionCategoryId,
		GeographicRegion: form.GeographicRegion,
		Province:         form.Province,
		District:         form.District,
		Stage:            form.Stage,
		Start:            form.Start,
		End:              form.End,
		CostOfProject:    &costOfProject,
		LandArea:         &landArea,
		ConstructionZone: &constructionZone,
	}
	constructionImages := []models.ConstructionImage{}
	if form.ConstructionImages != nil {
		for _, s := range form.ConstructionImages {
			image := models.ConstructionImage{
				RemoteLink: s.RemoteLink,
			}
			constructionImages = append(constructionImages, image)
		}
	}
	if len(constructionImages) > 0 {
		construction.ConstructionImages = constructionImages
	}
	construction.CompanyID = companyID
	features := []models.ConstructionFeature{}
	for _, f := range form.ConstructionFeatures {
		feat := models.ConstructionFeature{
			ConstructionFeatureID: f.FeatureID,
			Value:                 f.Value,
		}
		features = append(features, feat)
	}
	construction.ConstructionFeatures = features

	err := s.Repository.Save(construction)
	if err != nil {
		return nil, err
	}
	return construction, nil
}

func (s *ConstructionService) UpdateWithFields(construction *models.Construction, form forms.ConstructionUpdateForm, companyID uint) error {
	var fields = map[string]interface{}{
		"name": form.Name,
	}
	if form.GeographicRegion != nil {
		fields["geographic_region"] = *form.GeographicRegion
	}
	if form.Province != nil {
		fields["province"] = *form.Province
	}
	if form.District != nil {
		fields["district"] = *form.District
	}
	if form.Stage != nil {
		fields["stage"] = *form.Stage
	}
	if form.Start != nil {
		fields["start"] = *form.Start
	}
	if form.End != nil {
		fields["end"] = *form.End
	}
	if form.CostOfProject != nil {
		fields["cost_of_project"] = *form.CostOfProject
	}
	if form.LandArea != nil {
		fields["land_area"] = *form.LandArea
	}
	if form.ConstructionZone != nil {
		fields["construction_zone"] = *form.ConstructionZone
	}
	if form.ConstructionImages != nil {
		fields["construction_images"] = form.ConstructionImages
	}
	return s.Repository.UpdateWithFields(construction, fields)
}

func (s *ConstructionService) DeleteByID(constructionID uint, companyID uint) error {
	return s.Repository.DeleteByID(constructionID, companyID)
}

func (s *ConstructionService) FindCategoryByID(id uint) (*uint, error) {
	return s.Repository.FindCategoryByID(id)
}

func (s *ConstructionService) FindByIDAndCompanyID(constructionID uint, companyID uint) (*models.Construction, error) {
	return s.Repository.FindByIDAndCompanyID(constructionID, companyID)
}

func (s *ConstructionService) FindAllByCategoryID(categoryID uint) ([]models.Construction, error) {
	return s.Repository.FindAllByCategoryID(categoryID)
}

func (s *ConstructionService) FetchAuth(authD *models.AccessDetails) (uint, error) {
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
