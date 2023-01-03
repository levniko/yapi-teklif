package servicedefs

import (
	"github.com/sarulabs/dingo/v4"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
	"github.com/yapi-teklif/internal/repositories"
)

var RepositoriesDefs = []dingo.Def{
	{
		Name: "product-repository",
		Build: func(connection database.IConnection) (repositories.IProductRepository, error) {
			return repositories.GetProductRepository(connection), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("external-connections"),
		},
	},
	{
		Name: "variant-repository",
		Build: func(connection database.IConnection) (repositories.IVariantRepository, error) {
			return repositories.GetVariantRepository(connection), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("external-connections"),
		},
	},
	{
		Name: "feature-repository",
		Build: func(connection database.IConnection) (repositories.IFeatureRepository, error) {
			return repositories.GetFeatureRepository(connection), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("external-connections"),
		},
	},
	{
		Name: "company-repository",
		Build: func(connection database.IConnection) (repositories.ICompanyRepository, error) {
			return repositories.GetCompanyRepository(connection), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("external-connections"),
		},
	},
	{
		Name: "product-image-repository",
		Build: func(connection database.IConnection) (repositories.IProductImageRepository, error) {
			return repositories.GetProductImageRepository(connection), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("external-connections"),
		},
	},
	{
		Name: "construction-repository",
		Build: func(connection database.IConnection) (repositories.IConstructionRepository, error) {
			return repositories.GetConstructionRepository(connection), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("external-connections"),
		},
	},
	{
		Name: "variant-image-repository",
		Build: func(connection database.IConnection) (repositories.IVariantImageRepository, error) {
			return repositories.GetVariantImageRepository(connection), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("external-connections"),
		},
	},
	{
		Name: "construction-image-repository",
		Build: func(connection database.IConnection) (repositories.IConstructionImageRepository, error) {
			return repositories.GetConstructionImageRepository(connection), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("external-connections"),
		},
	},
	{
		Name: "construction-feature-repository",
		Build: func(connection database.IConnection) (repositories.IConstructionFeatureRepository, error) {
			return repositories.GetConstructionFeatureRepository(connection), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("external-connections"),
		},
	},
}
