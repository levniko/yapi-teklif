package servicedefs

import (
	"github.com/sarulabs/dingo/v4"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
	"github.com/yapi-teklif/internal/repositories"
	"github.com/yapi-teklif/internal/services"
)

var ServicesDefs = []dingo.Def{
	{
		Name: "product-service",
		Build: func(repository repositories.IProductRepository, cacheDB database.ICacheDB) (s services.IProductService, err error) {
			return services.GetProductService(repository, cacheDB), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("product-repository"),
			"1": dingo.Service("cache-connections"),
		},
	},
	{
		Name: "variant-service",
		Build: func(repository repositories.IVariantRepository, cacheDB database.ICacheDB) (s services.IVariantService, err error) {
			return services.GetVariantService(repository, cacheDB), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("variant-repository"),
			"1": dingo.Service("cache-connections"),
		},
	},
	{
		Name: "feature-service",
		Build: func(repository repositories.IFeatureRepository) (s services.IFeatureService, err error) {
			return services.GetFeatureService(repository), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("feature-repository"),
		},
	},
	{
		Name: "company-service",
		Build: func(repository repositories.ICompanyRepository, cacheDB database.ICacheDB) (s services.ICompanyService, err error) {
			return services.GetCompanyService(repository, cacheDB), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("company-repository"),
			"1": dingo.Service("cache-connections"),
		},
	},
	{
		Name: "product-image-service",
		Build: func(repository repositories.IProductImageRepository) (s services.IProductImageService, err error) {
			return services.GetProductImageService(repository), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("product-image-repository"),
		},
	},
	{
		Name: "construction-service",
		Build: func(repository repositories.IConstructionRepository, cacheDB database.ICacheDB) (s services.IConstructionService, err error) {
			return services.GetConstructionService(repository, cacheDB), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("construction-repository"),
			"1": dingo.Service("cache-connections"),
		},
	},
	{
		Name: "variant-image-service",
		Build: func(repository repositories.IVariantImageRepository) (s services.IVariantImageService, err error) {
			return services.GetVariantImageService(repository), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("variant-image-repository"),
		},
	},
	{
		Name: "construction-image-service",
		Build: func(repository repositories.IConstructionImageRepository) (s services.IConstructionImageService, err error) {
			return services.GetConstructionImageService(repository), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("construction-image-repository"),
		},
	},
	{
		Name: "construction-feature-service",
		Build: func(repository repositories.IConstructionFeatureRepository) (s services.IConstructionFeatureService, err error) {
			return services.GetConstructionFeatureService(repository), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("construction-feature-repository"),
		},
	},
}
