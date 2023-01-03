package managers

import (
	"errors"

	"github.com/yapi-teklif/internal/forms"
	"github.com/yapi-teklif/internal/models"
	"github.com/yapi-teklif/internal/pkg/utils"
	"github.com/yapi-teklif/internal/services"
)

type IProductManager interface {
	Save(form forms.ProductCreateForm, companyID uint) (*models.Product, error, int)
	Update(form forms.ProductUpdateForm, productID uint, companyID uint) (*models.Product, error, int)
	Delete(productID uint, companyID uint) (error, int)
	FindByID(productID uint, companyID uint) (*models.Product, error, int)
	FindAllByCategoryID(categoryID uint) ([]models.Product, error, int)
	UpdateProductImages(productID uint, form []forms.ProductImageForm) error
	FetchAuth(authD *models.AccessDetails) (uint, error)
}

type ProductManager struct {
	ProductService      services.IProductService
	ProductImageService services.IProductImageService
}

func GetProductManager(
	productService services.IProductService,
	productImageService services.IProductImageService) *ProductManager {
	return &ProductManager{
		ProductService:      productService,
		ProductImageService: productImageService,
	}
}

func (m *ProductManager) Save(form forms.ProductCreateForm, companyID uint) (*models.Product, error, int) {
	count, err := m.ProductService.CountBySpuAndCompanyID(form.SPU, companyID)
	if err != nil || count > 0 {
		return nil, errors.New(utils.ProductSpuMustBeUnique), utils.ProductSpuMustBeUniqueCode
	}
	model, err := m.ProductService.Save(form, companyID)
	if err != nil {
		return nil, errors.New(utils.ProductCanNotCreated), utils.ProductCanNotCreatedCode
	}
	return model, nil, 0
}

func (m *ProductManager) Update(form forms.ProductUpdateForm, productID uint, companyID uint) (*models.Product, error, int) {
	product, err := m.ProductService.FindByID(productID)
	if err != nil {
		return nil, errors.New(utils.ProductCanNotFound), utils.ProductCanNotFoundCode
	}
	if form.SPU != product.SPU {
		count, err := m.ProductService.CountBySpuAndCompanyID(form.SPU, companyID)
		if err != nil && count > 0 {
			return nil, errors.New(utils.ProductSpuMustBeUnique), utils.ProductSpuMustBeUniqueCode
		}
	}
	if len(form.ProductImages) > 0 {
		err = m.UpdateProductImages(product.ID, form.ProductImages)
		if err != nil {
			return nil, errors.New(utils.ProductImageCanNotUpdated), utils.ProductImageCanNotUpdatedCode
		}
	}
	err = m.ProductService.UpdateWithFields(product, form, companyID)
	if err != nil {
		return nil, errors.New(utils.ProductCanNotUpdated), utils.ProductCanNotUpdatedCode
	}
	return product, nil, 0
}

func (s *ProductManager) UpdateProductImages(productID uint, form []forms.ProductImageForm) error {
	for _, value := range form {
		pImage, err := s.ProductImageService.FindByRemoteLink(value.RemoteLink, productID)
		if err != nil {
			_ = s.ProductImageService.Save(&models.ProductImage{
				ProductID:  productID,
				RemoteLink: value.RemoteLink,
			})
			/*
				if err != nil {
					// TODO:log error here
				}*/
		} else {
			pImage.RemoteLink = value.RemoteLink
			_ = s.ProductImageService.Save(pImage)
			/*
				if err != nil {
					// TODO:log error here
				}*/
		}
	}
	return nil
}

func (m *ProductManager) Delete(productID uint, companyID uint) (error, int) {
	err := m.ProductService.DeleteByID(productID, companyID)
	if err != nil {
		return errors.New(utils.ProductCanNotDeleted), utils.ProductCanNotDeletedCode
	}
	return nil, 0
}

func (m *ProductManager) FindByID(productID uint, companyID uint) (*models.Product, error, int) {
	product, err := m.ProductService.FindByID(productID)
	if err != nil {
		return nil, errors.New(utils.ProductCanNotFound), utils.ProductCanNotFoundCode
	}
	return product, nil, 0
}

func (m *ProductManager) FindAllByCategoryID(categoryID uint) ([]models.Product, error, int) {
	products, err := m.ProductService.FindAllByCategoryID(categoryID)
	if err != nil {
		return nil, errors.New(utils.ProductCanNotFound), utils.ProductCanNotFoundCode
	}
	return products, nil, 0
}

func (m *ProductManager) FetchAuth(authD *models.AccessDetails) (uint, error) {
	return m.ProductService.FetchAuth(authD)
}
