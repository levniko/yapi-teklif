package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yapi-teklif/internal/managers"
	"github.com/yapi-teklif/internal/models"
	"github.com/yapi-teklif/internal/pkg/utils"
)

type FeaturesHandler struct {
	FeatureManager managers.IFeatureManager
}

func GetFeaturesHandler(featureManager managers.IFeatureManager) FeaturesHandler {
	return FeaturesHandler{
		FeatureManager: featureManager,
	}
}

func (handler *FeaturesHandler) Get(ctx *fiber.Ctx) (err error) {
	categoryID, _ := ctx.ParamsInt("id")

	features, err, errorCode := handler.FeatureManager.GetFeaturesByCategoryID(uint(categoryID))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ResponseModel{
			Message:   err.Error(),
			Error:     true,
			ErrorCode: errorCode,
			Data:      nil,
		})
	}
	return ctx.Status(http.StatusOK).JSON(models.ResponseModel{
		Message:   utils.VariantSuccessFullyFound,
		Error:     false,
		ErrorCode: utils.VariantSuccessFullyFoundCode,
		Data:      features,
	})
}
