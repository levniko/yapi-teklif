package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yapi-teklif/internal/forms"
	"github.com/yapi-teklif/internal/managers"
	"github.com/yapi-teklif/internal/models"
	"github.com/yapi-teklif/internal/pkg/utils"
)

type VariantsHandler struct {
	VariantManager managers.IVariantManager
}

func GetVariantsHandler(variantManager managers.IVariantManager) VariantsHandler {
	return VariantsHandler{
		VariantManager: variantManager,
	}
}

func (handler *VariantsHandler) Create(ctx *fiber.Ctx) (err error) {
	var form forms.VariantCreateForm
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
	metadata, err := ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	companyID, err := handler.VariantManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	variant, err, errCode := handler.VariantManager.Save(form, companyID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: errCode,
			Data:      nil,
		})

	}

	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.VariantSuccessFullyCreated,
		Error:     false,
		ErrorCode: utils.VariantSuccessFullyCreatedCode,
		Data:      variant,
	})
}

func (handler *VariantsHandler) Update(ctx *fiber.Ctx) (err error) {
	variantID, _ := ctx.ParamsInt("id")
	var form forms.VariantUpdateForm
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
	metadata, err := ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	companyID, err := handler.VariantManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	variant, err, errCode := handler.VariantManager.Update(form, uint(variantID), companyID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: errCode,
			Data:      nil,
		})

	}

	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.VariantSuccessFullyCreated,
		Error:     false,
		ErrorCode: utils.VariantSuccessFullyCreatedCode,
		Data:      variant,
	})
}

func (handler *VariantsHandler) Delete(ctx *fiber.Ctx) (err error) {
	variantID, _ := ctx.ParamsInt("id")
	metadata, err := ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	companyID, err := handler.VariantManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	err, errCode := handler.VariantManager.Delete(uint(variantID), companyID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: errCode,
			Data:      nil,
		})
	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.VariantSuccessFullyDeleted,
		Error:     false,
		ErrorCode: utils.VariantSuccessFullyDeletedCode,
		Data:      nil,
	})
}

func (handler *VariantsHandler) Get(ctx *fiber.Ctx) (err error) {
	id, _ := ctx.ParamsInt("id")
	metadata, err := ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	_, err = handler.VariantManager.FetchAuth(metadata)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: http.StatusUnauthorized,
			Data:      nil,
		})
	}
	variant, err := handler.VariantManager.FindByID(uint(id))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: utils.VariantCanNotFoundCode,
			Data:      nil,
		})
	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.VariantSuccessFullyFound,
		Error:     false,
		ErrorCode: utils.VariantSuccessFullyFoundCode,
		Data:      variant,
	})
}
