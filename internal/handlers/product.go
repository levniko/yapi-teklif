package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yapi-teklif/internal/forms"
	"github.com/yapi-teklif/internal/managers"
	"github.com/yapi-teklif/internal/models"
	"github.com/yapi-teklif/internal/pkg/utils"
)

type ProductsHandler struct {
	ProductManager managers.IProductManager
}

func GetProductsHandler(productManager managers.IProductManager) ProductsHandler {
	return ProductsHandler{
		ProductManager: productManager,
	}
}

func (handler *ProductsHandler) Create(ctx *fiber.Ctx) (err error) {
	var form forms.ProductCreateForm
	if err = ctx.BodyParser(&form); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message: err.Error(),
			Error:   true,
			Data:    nil,
		})
	}

	if err = utils.GetCustomValidator().Validate(form); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message: err.Error(),
			Error:   true,
			Data:    utils.CustomValidatorErr(err),
		})
	}
	//Extract the access token metadata
	metadata, err := ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	companyID, err := handler.ProductManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	product, err, errCode := handler.ProductManager.Save(form, companyID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: errCode,
			Data:      nil,
		})

	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.ProductSuccessfullyCreated,
		Error:     false,
		ErrorCode: utils.ProductSuccessfullyCreatedCode,
		Data:      product,
	})
}

func (handler *ProductsHandler) Update(ctx *fiber.Ctx) (err error) {
	var form forms.ProductUpdateForm
	if err = ctx.BodyParser(&form); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message: err.Error(),
			Error:   true,
			Data:    nil,
		})
	}

	if err = utils.GetCustomValidator().Validate(form); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message: err.Error(),
			Error:   true,
			Data:    utils.CustomValidatorErr(err),
		})
	}
	productID, _ := ctx.ParamsInt("id")
	//Extract the access token metadata
	metadata, err := ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	companyID, err := handler.ProductManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}

	product, err, errCode := handler.ProductManager.Update(form, uint(productID), companyID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: errCode,
			Data:      nil,
		})
	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.ProductSuccessfullyUpdated,
		Error:     false,
		ErrorCode: utils.ProductSuccessfullyUpdatedCode,
		Data:      product,
	})
}

func (handler *ProductsHandler) Delete(ctx *fiber.Ctx) (err error) {
	productID, _ := ctx.ParamsInt("id")
	//Extract the access token metadata
	metadata, err := ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	companyID, err := handler.ProductManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}

	err, errCode := handler.ProductManager.Delete(uint(productID), uint(companyID))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: errCode,
			Data:      nil,
		})
	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.ProductSuccessfullyDeleted,
		Error:     false,
		ErrorCode: utils.ProductSuccessfullyDeletedCode,
		Data:      nil,
	})
}

func (handler *ProductsHandler) Get(ctx *fiber.Ctx) (err error) {
	id, _ := ctx.ParamsInt("id")
	//Extract the access token metadata
	metadata, err := ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	companyID, err := handler.ProductManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	product, err, errCode := handler.ProductManager.FindByID(uint(id), uint(companyID))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: errCode,
			Data:      nil,
		})
	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.ProductSuccessFullyFound,
		Error:     false,
		ErrorCode: utils.ProductSuccessFullyFoundCode,
		Data:      product,
	})
}

func (handler *ProductsHandler) GetAllByCategory(ctx *fiber.Ctx) (err error) {
	categoryID, _ := ctx.ParamsInt("id")
	//Extract the access token metadata
	metadata, err := ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	_, err = handler.ProductManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	products, err, errCode := handler.ProductManager.FindAllByCategoryID(uint(categoryID))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: errCode,
			Data:      nil,
		})
	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.VariantSuccessFullyFound,
		Error:     false,
		ErrorCode: utils.VariantSuccessFullyFoundCode,
		Data:      products,
	})
}
