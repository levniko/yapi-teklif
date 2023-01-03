package servicedefs

import (
	"github.com/sarulabs/dingo/v4"
	"github.com/yapi-teklif/internal/managers"
	"github.com/yapi-teklif/internal/services"
)

var ManagersDefs = []dingo.Def{
	{
		Name: "product-manager",
		Build: func(productService services.IProductService, productImageService services.IProductImageService) (m managers.IProductManager, err error) {
			return managers.GetProductManager(productService, productImageService), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("product-service"),
			"1": dingo.Service("product-image-service"),
		},
	},
	{
		Name: "variant-manager",
		Build: func(variantService services.IVariantService,
			featureService services.IFeatureService,
			productService services.IProductService,
			variantImageService services.IVariantImageService,
		) (m managers.IVariantManager, err error) {
			return managers.GetVariantManager(variantService, featureService, productService, variantImageService), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("variant-service"),
			"1": dingo.Service("feature-service"),
			"2": dingo.Service("product-service"),
			"3": dingo.Service("variant-image-service"),
		},
	},
	{
		Name: "feature-manager",
		Build: func(featureService services.IFeatureService) (m managers.IFeatureManager, err error) {
			return managers.GetFeatureManager(featureService), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("feature-service"),
		},
	},
	{
		Name: "company-manager",
		Build: func(companyService services.ICompanyService) (m managers.ICompanyManager, err error) {
			return managers.GetCompanyManager(companyService), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("company-service"),
		},
	},
	{
		Name: "construction-manager",
		Build: func(constructionService services.IConstructionService,
			constructionFeatureService services.IConstructionFeatureService,
			constructionImageService services.IConstructionImageService) (m managers.IConstructionManager, err error) {
			return managers.GetConstructionManager(constructionService, constructionFeatureService, constructionImageService), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("construction-service"),
			"1": dingo.Service("construction-feature-service"),
			"2": dingo.Service("construction-image-service"),
		},
	},
}
