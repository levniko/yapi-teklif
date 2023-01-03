package servicedefs

import (
	"github.com/sarulabs/dingo/v4"
	"github.com/yapi-teklif/internal/handlers"
	"github.com/yapi-teklif/internal/managers"
)

var HandlersDefs = []dingo.Def{
	{
		Name: "product-handler",
		Build: func(productManager managers.IProductManager) (handlers.ProductsHandler, error) {
			return handlers.GetProductsHandler(productManager), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("product-manager"),
		},
	},
	{
		Name: "variant-handler",
		Build: func(variantManager managers.IVariantManager) (handlers.VariantsHandler, error) {
			return handlers.GetVariantsHandler(variantManager), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("variant-manager"),
		},
	},
	{
		Name: "feature-handler",
		Build: func(featureManager managers.IFeatureManager) (handlers.FeaturesHandler, error) {
			return handlers.GetFeaturesHandler(featureManager), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("feature-manager"),
		},
	},
	{
		Name: "company-handler",
		Build: func(companyManager managers.ICompanyManager) (handlers.CompaniesHandler, error) {
			return handlers.GetCompaniesHandler(companyManager), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("company-manager"),
		},
	},
	{
		Name: "construction-handler",
		Build: func(constructionManager managers.IConstructionManager) (handlers.ConstructionsHandler, error) {
			return handlers.GetConstructionsHandler(constructionManager), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("construction-manager"),
		},
	},
}
