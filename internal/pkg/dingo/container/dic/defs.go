package dic

import (
	"errors"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"

	handlers "github.com/yapi-teklif/internal/handlers"
	managers "github.com/yapi-teklif/internal/managers"
	connection "github.com/yapi-teklif/internal/pkg/database/connection"
	repositories "github.com/yapi-teklif/internal/repositories"
	services "github.com/yapi-teklif/internal/services"
)

func getDiDefs(provider dingo.Provider) []di.Def {
	return []di.Def{
		{
			Name:  "cache-connections",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("cache-connections")
				if err != nil {
					var eo connection.ICacheDB
					return eo, err
				}
				b, ok := d.Build.(func() (connection.ICacheDB, error))
				if !ok {
					var eo connection.ICacheDB
					return eo, errors.New("could not cast build function to func() (connection.ICacheDB, error)")
				}
				return b()
			},
			Unshared: false,
		},
		{
			Name:  "company-handler",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("company-handler")
				if err != nil {
					var eo handlers.CompaniesHandler
					return eo, err
				}
				pi0, err := ctn.SafeGet("company-manager")
				if err != nil {
					var eo handlers.CompaniesHandler
					return eo, err
				}
				p0, ok := pi0.(managers.ICompanyManager)
				if !ok {
					var eo handlers.CompaniesHandler
					return eo, errors.New("could not cast parameter 0 to managers.ICompanyManager")
				}
				b, ok := d.Build.(func(managers.ICompanyManager) (handlers.CompaniesHandler, error))
				if !ok {
					var eo handlers.CompaniesHandler
					return eo, errors.New("could not cast build function to func(managers.ICompanyManager) (handlers.CompaniesHandler, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "company-manager",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("company-manager")
				if err != nil {
					var eo managers.ICompanyManager
					return eo, err
				}
				pi0, err := ctn.SafeGet("company-service")
				if err != nil {
					var eo managers.ICompanyManager
					return eo, err
				}
				p0, ok := pi0.(services.ICompanyService)
				if !ok {
					var eo managers.ICompanyManager
					return eo, errors.New("could not cast parameter 0 to services.ICompanyService")
				}
				b, ok := d.Build.(func(services.ICompanyService) (managers.ICompanyManager, error))
				if !ok {
					var eo managers.ICompanyManager
					return eo, errors.New("could not cast build function to func(services.ICompanyService) (managers.ICompanyManager, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "company-repository",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("company-repository")
				if err != nil {
					var eo repositories.ICompanyRepository
					return eo, err
				}
				pi0, err := ctn.SafeGet("external-connections")
				if err != nil {
					var eo repositories.ICompanyRepository
					return eo, err
				}
				p0, ok := pi0.(connection.IConnection)
				if !ok {
					var eo repositories.ICompanyRepository
					return eo, errors.New("could not cast parameter 0 to connection.IConnection")
				}
				b, ok := d.Build.(func(connection.IConnection) (repositories.ICompanyRepository, error))
				if !ok {
					var eo repositories.ICompanyRepository
					return eo, errors.New("could not cast build function to func(connection.IConnection) (repositories.ICompanyRepository, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "company-service",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("company-service")
				if err != nil {
					var eo services.ICompanyService
					return eo, err
				}
				pi0, err := ctn.SafeGet("company-repository")
				if err != nil {
					var eo services.ICompanyService
					return eo, err
				}
				p0, ok := pi0.(repositories.ICompanyRepository)
				if !ok {
					var eo services.ICompanyService
					return eo, errors.New("could not cast parameter 0 to repositories.ICompanyRepository")
				}
				pi1, err := ctn.SafeGet("cache-connections")
				if err != nil {
					var eo services.ICompanyService
					return eo, err
				}
				p1, ok := pi1.(connection.ICacheDB)
				if !ok {
					var eo services.ICompanyService
					return eo, errors.New("could not cast parameter 1 to connection.ICacheDB")
				}
				b, ok := d.Build.(func(repositories.ICompanyRepository, connection.ICacheDB) (services.ICompanyService, error))
				if !ok {
					var eo services.ICompanyService
					return eo, errors.New("could not cast build function to func(repositories.ICompanyRepository, connection.ICacheDB) (services.ICompanyService, error)")
				}
				return b(p0, p1)
			},
			Unshared: false,
		},
		{
			Name:  "construction-feature-repository",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("construction-feature-repository")
				if err != nil {
					var eo repositories.IConstructionFeatureRepository
					return eo, err
				}
				pi0, err := ctn.SafeGet("external-connections")
				if err != nil {
					var eo repositories.IConstructionFeatureRepository
					return eo, err
				}
				p0, ok := pi0.(connection.IConnection)
				if !ok {
					var eo repositories.IConstructionFeatureRepository
					return eo, errors.New("could not cast parameter 0 to connection.IConnection")
				}
				b, ok := d.Build.(func(connection.IConnection) (repositories.IConstructionFeatureRepository, error))
				if !ok {
					var eo repositories.IConstructionFeatureRepository
					return eo, errors.New("could not cast build function to func(connection.IConnection) (repositories.IConstructionFeatureRepository, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "construction-feature-service",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("construction-feature-service")
				if err != nil {
					var eo services.IConstructionFeatureService
					return eo, err
				}
				pi0, err := ctn.SafeGet("construction-feature-repository")
				if err != nil {
					var eo services.IConstructionFeatureService
					return eo, err
				}
				p0, ok := pi0.(repositories.IConstructionFeatureRepository)
				if !ok {
					var eo services.IConstructionFeatureService
					return eo, errors.New("could not cast parameter 0 to repositories.IConstructionFeatureRepository")
				}
				b, ok := d.Build.(func(repositories.IConstructionFeatureRepository) (services.IConstructionFeatureService, error))
				if !ok {
					var eo services.IConstructionFeatureService
					return eo, errors.New("could not cast build function to func(repositories.IConstructionFeatureRepository) (services.IConstructionFeatureService, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "construction-handler",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("construction-handler")
				if err != nil {
					var eo handlers.ConstructionsHandler
					return eo, err
				}
				pi0, err := ctn.SafeGet("construction-manager")
				if err != nil {
					var eo handlers.ConstructionsHandler
					return eo, err
				}
				p0, ok := pi0.(managers.IConstructionManager)
				if !ok {
					var eo handlers.ConstructionsHandler
					return eo, errors.New("could not cast parameter 0 to managers.IConstructionManager")
				}
				b, ok := d.Build.(func(managers.IConstructionManager) (handlers.ConstructionsHandler, error))
				if !ok {
					var eo handlers.ConstructionsHandler
					return eo, errors.New("could not cast build function to func(managers.IConstructionManager) (handlers.ConstructionsHandler, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "construction-image-repository",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("construction-image-repository")
				if err != nil {
					var eo repositories.IConstructionImageRepository
					return eo, err
				}
				pi0, err := ctn.SafeGet("external-connections")
				if err != nil {
					var eo repositories.IConstructionImageRepository
					return eo, err
				}
				p0, ok := pi0.(connection.IConnection)
				if !ok {
					var eo repositories.IConstructionImageRepository
					return eo, errors.New("could not cast parameter 0 to connection.IConnection")
				}
				b, ok := d.Build.(func(connection.IConnection) (repositories.IConstructionImageRepository, error))
				if !ok {
					var eo repositories.IConstructionImageRepository
					return eo, errors.New("could not cast build function to func(connection.IConnection) (repositories.IConstructionImageRepository, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "construction-image-service",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("construction-image-service")
				if err != nil {
					var eo services.IConstructionImageService
					return eo, err
				}
				pi0, err := ctn.SafeGet("construction-image-repository")
				if err != nil {
					var eo services.IConstructionImageService
					return eo, err
				}
				p0, ok := pi0.(repositories.IConstructionImageRepository)
				if !ok {
					var eo services.IConstructionImageService
					return eo, errors.New("could not cast parameter 0 to repositories.IConstructionImageRepository")
				}
				b, ok := d.Build.(func(repositories.IConstructionImageRepository) (services.IConstructionImageService, error))
				if !ok {
					var eo services.IConstructionImageService
					return eo, errors.New("could not cast build function to func(repositories.IConstructionImageRepository) (services.IConstructionImageService, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "construction-manager",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("construction-manager")
				if err != nil {
					var eo managers.IConstructionManager
					return eo, err
				}
				pi0, err := ctn.SafeGet("construction-service")
				if err != nil {
					var eo managers.IConstructionManager
					return eo, err
				}
				p0, ok := pi0.(services.IConstructionService)
				if !ok {
					var eo managers.IConstructionManager
					return eo, errors.New("could not cast parameter 0 to services.IConstructionService")
				}
				pi1, err := ctn.SafeGet("construction-feature-service")
				if err != nil {
					var eo managers.IConstructionManager
					return eo, err
				}
				p1, ok := pi1.(services.IConstructionFeatureService)
				if !ok {
					var eo managers.IConstructionManager
					return eo, errors.New("could not cast parameter 1 to services.IConstructionFeatureService")
				}
				pi2, err := ctn.SafeGet("construction-image-service")
				if err != nil {
					var eo managers.IConstructionManager
					return eo, err
				}
				p2, ok := pi2.(services.IConstructionImageService)
				if !ok {
					var eo managers.IConstructionManager
					return eo, errors.New("could not cast parameter 2 to services.IConstructionImageService")
				}
				b, ok := d.Build.(func(services.IConstructionService, services.IConstructionFeatureService, services.IConstructionImageService) (managers.IConstructionManager, error))
				if !ok {
					var eo managers.IConstructionManager
					return eo, errors.New("could not cast build function to func(services.IConstructionService, services.IConstructionFeatureService, services.IConstructionImageService) (managers.IConstructionManager, error)")
				}
				return b(p0, p1, p2)
			},
			Unshared: false,
		},
		{
			Name:  "construction-repository",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("construction-repository")
				if err != nil {
					var eo repositories.IConstructionRepository
					return eo, err
				}
				pi0, err := ctn.SafeGet("external-connections")
				if err != nil {
					var eo repositories.IConstructionRepository
					return eo, err
				}
				p0, ok := pi0.(connection.IConnection)
				if !ok {
					var eo repositories.IConstructionRepository
					return eo, errors.New("could not cast parameter 0 to connection.IConnection")
				}
				b, ok := d.Build.(func(connection.IConnection) (repositories.IConstructionRepository, error))
				if !ok {
					var eo repositories.IConstructionRepository
					return eo, errors.New("could not cast build function to func(connection.IConnection) (repositories.IConstructionRepository, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "construction-service",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("construction-service")
				if err != nil {
					var eo services.IConstructionService
					return eo, err
				}
				pi0, err := ctn.SafeGet("construction-repository")
				if err != nil {
					var eo services.IConstructionService
					return eo, err
				}
				p0, ok := pi0.(repositories.IConstructionRepository)
				if !ok {
					var eo services.IConstructionService
					return eo, errors.New("could not cast parameter 0 to repositories.IConstructionRepository")
				}
				pi1, err := ctn.SafeGet("cache-connections")
				if err != nil {
					var eo services.IConstructionService
					return eo, err
				}
				p1, ok := pi1.(connection.ICacheDB)
				if !ok {
					var eo services.IConstructionService
					return eo, errors.New("could not cast parameter 1 to connection.ICacheDB")
				}
				b, ok := d.Build.(func(repositories.IConstructionRepository, connection.ICacheDB) (services.IConstructionService, error))
				if !ok {
					var eo services.IConstructionService
					return eo, errors.New("could not cast build function to func(repositories.IConstructionRepository, connection.ICacheDB) (services.IConstructionService, error)")
				}
				return b(p0, p1)
			},
			Unshared: false,
		},
		{
			Name:  "external-connections",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("external-connections")
				if err != nil {
					var eo connection.IConnection
					return eo, err
				}
				b, ok := d.Build.(func() (connection.IConnection, error))
				if !ok {
					var eo connection.IConnection
					return eo, errors.New("could not cast build function to func() (connection.IConnection, error)")
				}
				return b()
			},
			Unshared: false,
		},
		{
			Name:  "feature-handler",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("feature-handler")
				if err != nil {
					var eo handlers.FeaturesHandler
					return eo, err
				}
				pi0, err := ctn.SafeGet("feature-manager")
				if err != nil {
					var eo handlers.FeaturesHandler
					return eo, err
				}
				p0, ok := pi0.(managers.IFeatureManager)
				if !ok {
					var eo handlers.FeaturesHandler
					return eo, errors.New("could not cast parameter 0 to managers.IFeatureManager")
				}
				b, ok := d.Build.(func(managers.IFeatureManager) (handlers.FeaturesHandler, error))
				if !ok {
					var eo handlers.FeaturesHandler
					return eo, errors.New("could not cast build function to func(managers.IFeatureManager) (handlers.FeaturesHandler, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "feature-manager",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("feature-manager")
				if err != nil {
					var eo managers.IFeatureManager
					return eo, err
				}
				pi0, err := ctn.SafeGet("feature-service")
				if err != nil {
					var eo managers.IFeatureManager
					return eo, err
				}
				p0, ok := pi0.(services.IFeatureService)
				if !ok {
					var eo managers.IFeatureManager
					return eo, errors.New("could not cast parameter 0 to services.IFeatureService")
				}
				b, ok := d.Build.(func(services.IFeatureService) (managers.IFeatureManager, error))
				if !ok {
					var eo managers.IFeatureManager
					return eo, errors.New("could not cast build function to func(services.IFeatureService) (managers.IFeatureManager, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "feature-repository",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("feature-repository")
				if err != nil {
					var eo repositories.IFeatureRepository
					return eo, err
				}
				pi0, err := ctn.SafeGet("external-connections")
				if err != nil {
					var eo repositories.IFeatureRepository
					return eo, err
				}
				p0, ok := pi0.(connection.IConnection)
				if !ok {
					var eo repositories.IFeatureRepository
					return eo, errors.New("could not cast parameter 0 to connection.IConnection")
				}
				b, ok := d.Build.(func(connection.IConnection) (repositories.IFeatureRepository, error))
				if !ok {
					var eo repositories.IFeatureRepository
					return eo, errors.New("could not cast build function to func(connection.IConnection) (repositories.IFeatureRepository, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "feature-service",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("feature-service")
				if err != nil {
					var eo services.IFeatureService
					return eo, err
				}
				pi0, err := ctn.SafeGet("feature-repository")
				if err != nil {
					var eo services.IFeatureService
					return eo, err
				}
				p0, ok := pi0.(repositories.IFeatureRepository)
				if !ok {
					var eo services.IFeatureService
					return eo, errors.New("could not cast parameter 0 to repositories.IFeatureRepository")
				}
				b, ok := d.Build.(func(repositories.IFeatureRepository) (services.IFeatureService, error))
				if !ok {
					var eo services.IFeatureService
					return eo, errors.New("could not cast build function to func(repositories.IFeatureRepository) (services.IFeatureService, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "product-handler",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("product-handler")
				if err != nil {
					var eo handlers.ProductsHandler
					return eo, err
				}
				pi0, err := ctn.SafeGet("product-manager")
				if err != nil {
					var eo handlers.ProductsHandler
					return eo, err
				}
				p0, ok := pi0.(managers.IProductManager)
				if !ok {
					var eo handlers.ProductsHandler
					return eo, errors.New("could not cast parameter 0 to managers.IProductManager")
				}
				b, ok := d.Build.(func(managers.IProductManager) (handlers.ProductsHandler, error))
				if !ok {
					var eo handlers.ProductsHandler
					return eo, errors.New("could not cast build function to func(managers.IProductManager) (handlers.ProductsHandler, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "product-image-repository",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("product-image-repository")
				if err != nil {
					var eo repositories.IProductImageRepository
					return eo, err
				}
				pi0, err := ctn.SafeGet("external-connections")
				if err != nil {
					var eo repositories.IProductImageRepository
					return eo, err
				}
				p0, ok := pi0.(connection.IConnection)
				if !ok {
					var eo repositories.IProductImageRepository
					return eo, errors.New("could not cast parameter 0 to connection.IConnection")
				}
				b, ok := d.Build.(func(connection.IConnection) (repositories.IProductImageRepository, error))
				if !ok {
					var eo repositories.IProductImageRepository
					return eo, errors.New("could not cast build function to func(connection.IConnection) (repositories.IProductImageRepository, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "product-image-service",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("product-image-service")
				if err != nil {
					var eo services.IProductImageService
					return eo, err
				}
				pi0, err := ctn.SafeGet("product-image-repository")
				if err != nil {
					var eo services.IProductImageService
					return eo, err
				}
				p0, ok := pi0.(repositories.IProductImageRepository)
				if !ok {
					var eo services.IProductImageService
					return eo, errors.New("could not cast parameter 0 to repositories.IProductImageRepository")
				}
				b, ok := d.Build.(func(repositories.IProductImageRepository) (services.IProductImageService, error))
				if !ok {
					var eo services.IProductImageService
					return eo, errors.New("could not cast build function to func(repositories.IProductImageRepository) (services.IProductImageService, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "product-manager",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("product-manager")
				if err != nil {
					var eo managers.IProductManager
					return eo, err
				}
				pi0, err := ctn.SafeGet("product-service")
				if err != nil {
					var eo managers.IProductManager
					return eo, err
				}
				p0, ok := pi0.(services.IProductService)
				if !ok {
					var eo managers.IProductManager
					return eo, errors.New("could not cast parameter 0 to services.IProductService")
				}
				pi1, err := ctn.SafeGet("product-image-service")
				if err != nil {
					var eo managers.IProductManager
					return eo, err
				}
				p1, ok := pi1.(services.IProductImageService)
				if !ok {
					var eo managers.IProductManager
					return eo, errors.New("could not cast parameter 1 to services.IProductImageService")
				}
				b, ok := d.Build.(func(services.IProductService, services.IProductImageService) (managers.IProductManager, error))
				if !ok {
					var eo managers.IProductManager
					return eo, errors.New("could not cast build function to func(services.IProductService, services.IProductImageService) (managers.IProductManager, error)")
				}
				return b(p0, p1)
			},
			Unshared: false,
		},
		{
			Name:  "product-repository",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("product-repository")
				if err != nil {
					var eo repositories.IProductRepository
					return eo, err
				}
				pi0, err := ctn.SafeGet("external-connections")
				if err != nil {
					var eo repositories.IProductRepository
					return eo, err
				}
				p0, ok := pi0.(connection.IConnection)
				if !ok {
					var eo repositories.IProductRepository
					return eo, errors.New("could not cast parameter 0 to connection.IConnection")
				}
				b, ok := d.Build.(func(connection.IConnection) (repositories.IProductRepository, error))
				if !ok {
					var eo repositories.IProductRepository
					return eo, errors.New("could not cast build function to func(connection.IConnection) (repositories.IProductRepository, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "product-service",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("product-service")
				if err != nil {
					var eo services.IProductService
					return eo, err
				}
				pi0, err := ctn.SafeGet("product-repository")
				if err != nil {
					var eo services.IProductService
					return eo, err
				}
				p0, ok := pi0.(repositories.IProductRepository)
				if !ok {
					var eo services.IProductService
					return eo, errors.New("could not cast parameter 0 to repositories.IProductRepository")
				}
				pi1, err := ctn.SafeGet("cache-connections")
				if err != nil {
					var eo services.IProductService
					return eo, err
				}
				p1, ok := pi1.(connection.ICacheDB)
				if !ok {
					var eo services.IProductService
					return eo, errors.New("could not cast parameter 1 to connection.ICacheDB")
				}
				b, ok := d.Build.(func(repositories.IProductRepository, connection.ICacheDB) (services.IProductService, error))
				if !ok {
					var eo services.IProductService
					return eo, errors.New("could not cast build function to func(repositories.IProductRepository, connection.ICacheDB) (services.IProductService, error)")
				}
				return b(p0, p1)
			},
			Unshared: false,
		},
		{
			Name:  "variant-handler",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("variant-handler")
				if err != nil {
					var eo handlers.VariantsHandler
					return eo, err
				}
				pi0, err := ctn.SafeGet("variant-manager")
				if err != nil {
					var eo handlers.VariantsHandler
					return eo, err
				}
				p0, ok := pi0.(managers.IVariantManager)
				if !ok {
					var eo handlers.VariantsHandler
					return eo, errors.New("could not cast parameter 0 to managers.IVariantManager")
				}
				b, ok := d.Build.(func(managers.IVariantManager) (handlers.VariantsHandler, error))
				if !ok {
					var eo handlers.VariantsHandler
					return eo, errors.New("could not cast build function to func(managers.IVariantManager) (handlers.VariantsHandler, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "variant-image-repository",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("variant-image-repository")
				if err != nil {
					var eo repositories.IVariantImageRepository
					return eo, err
				}
				pi0, err := ctn.SafeGet("external-connections")
				if err != nil {
					var eo repositories.IVariantImageRepository
					return eo, err
				}
				p0, ok := pi0.(connection.IConnection)
				if !ok {
					var eo repositories.IVariantImageRepository
					return eo, errors.New("could not cast parameter 0 to connection.IConnection")
				}
				b, ok := d.Build.(func(connection.IConnection) (repositories.IVariantImageRepository, error))
				if !ok {
					var eo repositories.IVariantImageRepository
					return eo, errors.New("could not cast build function to func(connection.IConnection) (repositories.IVariantImageRepository, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "variant-image-service",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("variant-image-service")
				if err != nil {
					var eo services.IVariantImageService
					return eo, err
				}
				pi0, err := ctn.SafeGet("variant-image-repository")
				if err != nil {
					var eo services.IVariantImageService
					return eo, err
				}
				p0, ok := pi0.(repositories.IVariantImageRepository)
				if !ok {
					var eo services.IVariantImageService
					return eo, errors.New("could not cast parameter 0 to repositories.IVariantImageRepository")
				}
				b, ok := d.Build.(func(repositories.IVariantImageRepository) (services.IVariantImageService, error))
				if !ok {
					var eo services.IVariantImageService
					return eo, errors.New("could not cast build function to func(repositories.IVariantImageRepository) (services.IVariantImageService, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "variant-manager",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("variant-manager")
				if err != nil {
					var eo managers.IVariantManager
					return eo, err
				}
				pi0, err := ctn.SafeGet("variant-service")
				if err != nil {
					var eo managers.IVariantManager
					return eo, err
				}
				p0, ok := pi0.(services.IVariantService)
				if !ok {
					var eo managers.IVariantManager
					return eo, errors.New("could not cast parameter 0 to services.IVariantService")
				}
				pi1, err := ctn.SafeGet("feature-service")
				if err != nil {
					var eo managers.IVariantManager
					return eo, err
				}
				p1, ok := pi1.(services.IFeatureService)
				if !ok {
					var eo managers.IVariantManager
					return eo, errors.New("could not cast parameter 1 to services.IFeatureService")
				}
				pi2, err := ctn.SafeGet("product-service")
				if err != nil {
					var eo managers.IVariantManager
					return eo, err
				}
				p2, ok := pi2.(services.IProductService)
				if !ok {
					var eo managers.IVariantManager
					return eo, errors.New("could not cast parameter 2 to services.IProductService")
				}
				pi3, err := ctn.SafeGet("variant-image-service")
				if err != nil {
					var eo managers.IVariantManager
					return eo, err
				}
				p3, ok := pi3.(services.IVariantImageService)
				if !ok {
					var eo managers.IVariantManager
					return eo, errors.New("could not cast parameter 3 to services.IVariantImageService")
				}
				b, ok := d.Build.(func(services.IVariantService, services.IFeatureService, services.IProductService, services.IVariantImageService) (managers.IVariantManager, error))
				if !ok {
					var eo managers.IVariantManager
					return eo, errors.New("could not cast build function to func(services.IVariantService, services.IFeatureService, services.IProductService, services.IVariantImageService) (managers.IVariantManager, error)")
				}
				return b(p0, p1, p2, p3)
			},
			Unshared: false,
		},
		{
			Name:  "variant-repository",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("variant-repository")
				if err != nil {
					var eo repositories.IVariantRepository
					return eo, err
				}
				pi0, err := ctn.SafeGet("external-connections")
				if err != nil {
					var eo repositories.IVariantRepository
					return eo, err
				}
				p0, ok := pi0.(connection.IConnection)
				if !ok {
					var eo repositories.IVariantRepository
					return eo, errors.New("could not cast parameter 0 to connection.IConnection")
				}
				b, ok := d.Build.(func(connection.IConnection) (repositories.IVariantRepository, error))
				if !ok {
					var eo repositories.IVariantRepository
					return eo, errors.New("could not cast build function to func(connection.IConnection) (repositories.IVariantRepository, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "variant-service",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("variant-service")
				if err != nil {
					var eo services.IVariantService
					return eo, err
				}
				pi0, err := ctn.SafeGet("variant-repository")
				if err != nil {
					var eo services.IVariantService
					return eo, err
				}
				p0, ok := pi0.(repositories.IVariantRepository)
				if !ok {
					var eo services.IVariantService
					return eo, errors.New("could not cast parameter 0 to repositories.IVariantRepository")
				}
				pi1, err := ctn.SafeGet("cache-connections")
				if err != nil {
					var eo services.IVariantService
					return eo, err
				}
				p1, ok := pi1.(connection.ICacheDB)
				if !ok {
					var eo services.IVariantService
					return eo, errors.New("could not cast parameter 1 to connection.ICacheDB")
				}
				b, ok := d.Build.(func(repositories.IVariantRepository, connection.ICacheDB) (services.IVariantService, error))
				if !ok {
					var eo services.IVariantService
					return eo, errors.New("could not cast build function to func(repositories.IVariantRepository, connection.ICacheDB) (services.IVariantService, error)")
				}
				return b(p0, p1)
			},
			Unshared: false,
		},
	}
}
