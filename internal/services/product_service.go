package services

import (
	"errors"
	"strconv"

	"github.com/yapi-teklif/internal/forms"
	"github.com/yapi-teklif/internal/models"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
	"github.com/yapi-teklif/internal/repositories"
)

type IProductService interface {
	Save(form forms.ProductCreateForm, companyID uint) (*models.Product, error)
	UpdateWithFields(product *models.Product, form forms.ProductUpdateForm, companyID uint) error
	DeleteByID(productID uint, companyID uint) error
	FindCategoryByID(id uint) (*uint, error)
	FindByID(productID uint) (*models.Product, error)
	CountBySpuAndCompanyID(spu string, companyID uint) (int64, error)
	FindAllByCategoryID(categoryID uint) ([]models.Product, error)
	FetchAuth(authD *models.AccessDetails) (uint, error)
}

type ProductService struct {
	Repository repositories.IProductRepository
	CacheDB    database.ICacheDB
}

func GetProductService(r repositories.IProductRepository, c database.ICacheDB) *ProductService {
	return &ProductService{
		Repository: r,
		CacheDB:    c,
	}
}

func (s *ProductService) FindCategoryByID(id uint) (*uint, error) {
	return s.Repository.FindCategoryByID(id)
}

func (s *ProductService) Save(form forms.ProductCreateForm, companyID uint) (*models.Product, error) {
	product := &models.Product{
		Name:              form.Name,
		SPU:               form.SPU,
		CompanyID:         uint(companyID),
		HeroImage:         form.HeroImage,
		Description:       form.Description,
		IsActive:          form.IsActive,
		ProductCategoryID: form.CategoryID,
	}
	productImages := []models.ProductImage{}
	if form.ProductImages != nil {
		for _, s := range form.ProductImages {
			image := models.ProductImage{
				RemoteLink: s.RemoteLink,
			}
			productImages = append(productImages, image)
		}
	}
	if len(productImages) > 0 {
		product.ProductImages = productImages
	}
	err := s.Repository.Save(product)

	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) UpdateWithFields(product *models.Product, form forms.ProductUpdateForm, companyID uint) error {
	var fields = map[string]interface{}{
		"name": form.Name,
		"spu":  form.SPU, //TODO: Check SPU is unique for this company
	}
	if form.IsActive != nil {
		fields["is_active"] = *form.IsActive
	}
	if form.Description != nil {
		fields["description"] = *form.Description
	}
	if form.HeroImage != nil {
		fields["hero_image"] = *form.HeroImage
	}
	return s.Repository.UpdateWithFields(product, fields)
}

func (s *ProductService) DeleteByID(productID uint, companyID uint) error {
	return s.Repository.DeleteByID(productID, companyID)
}

func (s *ProductService) FindByID(productID uint) (*models.Product, error) {
	return s.Repository.FindByID(productID)
}

func (s *ProductService) CountBySpuAndCompanyID(spu string, companyID uint) (int64, error) {
	return s.Repository.CountBySpuAndCompanyID(spu, companyID)
}

func (s *ProductService) FindAllByCategoryID(categoryID uint) ([]models.Product, error) {
	return s.Repository.FindAllByCategoryID(categoryID)
}

func (s *ProductService) FetchAuth(authD *models.AccessDetails) (uint, error) {
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
