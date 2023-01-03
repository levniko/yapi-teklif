package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yapi-teklif/internal/forms"
	"github.com/yapi-teklif/internal/managers"
	"github.com/yapi-teklif/internal/models"
	"github.com/yapi-teklif/internal/pkg/utils"
)

type ConstructionsHandler struct {
	ConstructionManager managers.IConstructionManager
}

func GetConstructionsHandler(constructionManager managers.IConstructionManager) ConstructionsHandler {
	return ConstructionsHandler{
		ConstructionManager: constructionManager,
	}
}

func (handler *ConstructionsHandler) Create(ctx *fiber.Ctx) (err error) {
	var form forms.ConstructionCreateForm
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
	companyID, err := handler.ConstructionManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	product, err, errCode := handler.ConstructionManager.Save(form, companyID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: errCode,
			Data:      nil,
		})

	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.ConstructionSuccessfullyCreated,
		Error:     false,
		ErrorCode: utils.ConstructionSuccessfullyCreatedCode,
		Data:      product,
	})
}

func (handler *ConstructionsHandler) Update(ctx *fiber.Ctx) (err error) {
	var form forms.ConstructionUpdateForm
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
	companyID, err := handler.ConstructionManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}

	construction, err, errCode := handler.ConstructionManager.Update(form, uint(productID), companyID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: errCode,
			Data:      nil,
		})
	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.ConstructionSuccessfullyUpdated,
		Error:     false,
		ErrorCode: utils.ConstructionSuccessfullyUpdatedCode,
		Data:      construction,
	})
}

func (handler *ConstructionsHandler) Delete(ctx *fiber.Ctx) (err error) {
	constructionID, _ := ctx.ParamsInt("id")
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
	companyID, err := handler.ConstructionManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}

	err, errCode := handler.ConstructionManager.Delete(uint(constructionID), uint(companyID))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: errCode,
			Data:      nil,
		})
	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.ConstructionSuccessfullyDeleted,
		Error:     false,
		ErrorCode: utils.ConstructionSuccessfullyDeletedCode,
		Data:      nil,
	})
}

func (handler *ConstructionsHandler) Get(ctx *fiber.Ctx) (err error) {
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
	companyID, err := handler.ConstructionManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	construction, err, errCode := handler.ConstructionManager.FindByIDAndCompanyID(uint(id), uint(companyID))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: errCode,
			Data:      nil,
		})
	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.ConstructionSuccessFullyFound,
		Error:     false,
		ErrorCode: utils.ConstructionSuccessFullyFoundCode,
		Data:      construction,
	})
}

func (handler *ConstructionsHandler) GetAllByCategory(ctx *fiber.Ctx) (err error) {
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
	_, err = handler.ConstructionManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	constructions, err, errCode := handler.ConstructionManager.FindAllByCategoryID(uint(categoryID))
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
		Data:      constructions,
	})
}
