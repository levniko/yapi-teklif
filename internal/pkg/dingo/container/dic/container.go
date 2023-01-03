package dic

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"

	providerPkg "github.com/yapi-teklif/internal/pkg/dingo/provider"

	handlers "github.com/yapi-teklif/internal/handlers"
	managers "github.com/yapi-teklif/internal/managers"
	connection "github.com/yapi-teklif/internal/pkg/database/connection"
	repositories "github.com/yapi-teklif/internal/repositories"
	services "github.com/yapi-teklif/internal/services"
)

// C retrieves a Container from an interface.
// The function panics if the Container can not be retrieved.
//
// The interface can be :
//   - a *Container
//   - an *http.Request containing a *Container in its context.Context
//     for the dingo.ContainerKey("dingo") key.
//
// The function can be changed to match the needs of your application.
var C = func(i interface{}) *Container {
	if c, ok := i.(*Container); ok {
		return c
	}
	r, ok := i.(*http.Request)
	if !ok {
		panic("could not get the container with dic.C()")
	}
	c, ok := r.Context().Value(dingo.ContainerKey("dingo")).(*Container)
	if !ok {
		panic("could not get the container from the given *http.Request in dic.C()")
	}
	return c
}

type builder struct {
	builder *di.Builder
}

// NewBuilder creates a builder that can be used to create a Container.
// You probably should use NewContainer to create the container directly.
// But using NewBuilder allows you to redefine some di services.
// This can be used for testing.
// But this behavior is not safe, so be sure to know what you are doing.
func NewBuilder(scopes ...string) (*builder, error) {
	if len(scopes) == 0 {
		scopes = []string{di.App, di.Request, di.SubRequest}
	}
	b, err := di.NewBuilder(scopes...)
	if err != nil {
		return nil, fmt.Errorf("could not create di.Builder: %v", err)
	}
	provider := &providerPkg.Provider{}
	if err := provider.Load(); err != nil {
		return nil, fmt.Errorf("could not load definitions with the Provider (Provider from github.com/yapi-teklif/internal/pkg/dingo/provider): %v", err)
	}
	for _, d := range getDiDefs(provider) {
		if err := b.Add(d); err != nil {
			return nil, fmt.Errorf("could not add di.Def in di.Builder: %v", err)
		}
	}
	return &builder{builder: b}, nil
}

// Add adds one or more definitions in the Builder.
// It returns an error if a definition can not be added.
func (b *builder) Add(defs ...di.Def) error {
	return b.builder.Add(defs...)
}

// Set is a shortcut to add a definition for an already built object.
func (b *builder) Set(name string, obj interface{}) error {
	return b.builder.Set(name, obj)
}

// Build creates a Container in the most generic scope.
func (b *builder) Build() *Container {
	return &Container{ctn: b.builder.Build()}
}

// NewContainer creates a new Container.
// If no scope is provided, di.App, di.Request and di.SubRequest are used.
// The returned Container has the most generic scope (di.App).
// The SubContainer() method should be called to get a Container in a more specific scope.
func NewContainer(scopes ...string) (*Container, error) {
	b, err := NewBuilder(scopes...)
	if err != nil {
		return nil, err
	}
	return b.Build(), nil
}

// Container represents a generated dependency injection container.
// It is a wrapper around a di.Container.
//
// A Container has a scope and may have a parent in a more generic scope
// and children in a more specific scope.
// Objects can be retrieved from the Container.
// If the requested object does not already exist in the Container,
// it is built thanks to the object definition.
// The following attempts to get this object will return the same object.
type Container struct {
	ctn di.Container
}

// Scope returns the Container scope.
func (c *Container) Scope() string {
	return c.ctn.Scope()
}

// Scopes returns the list of available scopes.
func (c *Container) Scopes() []string {
	return c.ctn.Scopes()
}

// ParentScopes returns the list of scopes wider than the Container scope.
func (c *Container) ParentScopes() []string {
	return c.ctn.ParentScopes()
}

// SubScopes returns the list of scopes that are more specific than the Container scope.
func (c *Container) SubScopes() []string {
	return c.ctn.SubScopes()
}

// Parent returns the parent Container.
func (c *Container) Parent() *Container {
	if p := c.ctn.Parent(); p != nil {
		return &Container{ctn: p}
	}
	return nil
}

// SubContainer creates a new Container in the next sub-scope
// that will have this Container as parent.
func (c *Container) SubContainer() (*Container, error) {
	sub, err := c.ctn.SubContainer()
	if err != nil {
		return nil, err
	}
	return &Container{ctn: sub}, nil
}

// SafeGet retrieves an object from the Container.
// The object has to belong to this scope or a more generic one.
// If the object does not already exist, it is created and saved in the Container.
// If the object can not be created, it returns an error.
func (c *Container) SafeGet(name string) (interface{}, error) {
	return c.ctn.SafeGet(name)
}

// Get is similar to SafeGet but it does not return the error.
// Instead it panics.
func (c *Container) Get(name string) interface{} {
	return c.ctn.Get(name)
}

// Fill is similar to SafeGet but it does not return the object.
// Instead it fills the provided object with the value returned by SafeGet.
// The provided object must be a pointer to the value returned by SafeGet.
func (c *Container) Fill(name string, dst interface{}) error {
	return c.ctn.Fill(name, dst)
}

// UnscopedSafeGet retrieves an object from the Container, like SafeGet.
// The difference is that the object can be retrieved
// even if it belongs to a more specific scope.
// To do so, UnscopedSafeGet creates a sub-container.
// When the created object is no longer needed,
// it is important to use the Clean method to delete this sub-container.
func (c *Container) UnscopedSafeGet(name string) (interface{}, error) {
	return c.ctn.UnscopedSafeGet(name)
}

// UnscopedGet is similar to UnscopedSafeGet but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGet(name string) interface{} {
	return c.ctn.UnscopedGet(name)
}

// UnscopedFill is similar to UnscopedSafeGet but copies the object in dst instead of returning it.
func (c *Container) UnscopedFill(name string, dst interface{}) error {
	return c.ctn.UnscopedFill(name, dst)
}

// Clean deletes the sub-container created by UnscopedSafeGet, UnscopedGet or UnscopedFill.
func (c *Container) Clean() error {
	return c.ctn.Clean()
}

// DeleteWithSubContainers takes all the objects saved in this Container
// and calls the Close function of their Definition on them.
// It will also call DeleteWithSubContainers on each child and remove its reference in the parent Container.
// After deletion, the Container can no longer be used.
// The sub-containers are deleted even if they are still used in other goroutines.
// It can cause errors. You may want to use the Delete method instead.
func (c *Container) DeleteWithSubContainers() error {
	return c.ctn.DeleteWithSubContainers()
}

// Delete works like DeleteWithSubContainers if the Container does not have any child.
// But if the Container has sub-containers, it will not be deleted right away.
// The deletion only occurs when all the sub-containers have been deleted manually.
// So you have to call Delete or DeleteWithSubContainers on all the sub-containers.
func (c *Container) Delete() error {
	return c.ctn.Delete()
}

// IsClosed returns true if the Container has been deleted.
func (c *Container) IsClosed() bool {
	return c.ctn.IsClosed()
}

// SafeGetCacheConnections retrieves the "cache-connections" object from the main scope.
//
// ---------------------------------------------
//
//	name: "cache-connections"
//	type: connection.ICacheDB
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetCacheConnections() (connection.ICacheDB, error) {
	i, err := c.ctn.SafeGet("cache-connections")
	if err != nil {
		var eo connection.ICacheDB
		return eo, err
	}
	o, ok := i.(connection.ICacheDB)
	if !ok {
		return o, errors.New("could get 'cache-connections' because the object could not be cast to connection.ICacheDB")
	}
	return o, nil
}

// GetCacheConnections retrieves the "cache-connections" object from the main scope.
//
// ---------------------------------------------
//
//	name: "cache-connections"
//	type: connection.ICacheDB
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetCacheConnections() connection.ICacheDB {
	o, err := c.SafeGetCacheConnections()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCacheConnections retrieves the "cache-connections" object from the main scope.
//
// ---------------------------------------------
//
//	name: "cache-connections"
//	type: connection.ICacheDB
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetCacheConnections() (connection.ICacheDB, error) {
	i, err := c.ctn.UnscopedSafeGet("cache-connections")
	if err != nil {
		var eo connection.ICacheDB
		return eo, err
	}
	o, ok := i.(connection.ICacheDB)
	if !ok {
		return o, errors.New("could get 'cache-connections' because the object could not be cast to connection.ICacheDB")
	}
	return o, nil
}

// UnscopedGetCacheConnections retrieves the "cache-connections" object from the main scope.
//
// ---------------------------------------------
//
//	name: "cache-connections"
//	type: connection.ICacheDB
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetCacheConnections() connection.ICacheDB {
	o, err := c.UnscopedSafeGetCacheConnections()
	if err != nil {
		panic(err)
	}
	return o
}

// CacheConnections retrieves the "cache-connections" object from the main scope.
//
// ---------------------------------------------
//
//	name: "cache-connections"
//	type: connection.ICacheDB
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetCacheConnections method.
// If the container can not be retrieved, it panics.
func CacheConnections(i interface{}) connection.ICacheDB {
	return C(i).GetCacheConnections()
}

// SafeGetCompanyHandler retrieves the "company-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-handler"
//	type: handlers.CompaniesHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.ICompanyManager) ["company-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetCompanyHandler() (handlers.CompaniesHandler, error) {
	i, err := c.ctn.SafeGet("company-handler")
	if err != nil {
		var eo handlers.CompaniesHandler
		return eo, err
	}
	o, ok := i.(handlers.CompaniesHandler)
	if !ok {
		return o, errors.New("could get 'company-handler' because the object could not be cast to handlers.CompaniesHandler")
	}
	return o, nil
}

// GetCompanyHandler retrieves the "company-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-handler"
//	type: handlers.CompaniesHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.ICompanyManager) ["company-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetCompanyHandler() handlers.CompaniesHandler {
	o, err := c.SafeGetCompanyHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCompanyHandler retrieves the "company-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-handler"
//	type: handlers.CompaniesHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.ICompanyManager) ["company-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetCompanyHandler() (handlers.CompaniesHandler, error) {
	i, err := c.ctn.UnscopedSafeGet("company-handler")
	if err != nil {
		var eo handlers.CompaniesHandler
		return eo, err
	}
	o, ok := i.(handlers.CompaniesHandler)
	if !ok {
		return o, errors.New("could get 'company-handler' because the object could not be cast to handlers.CompaniesHandler")
	}
	return o, nil
}

// UnscopedGetCompanyHandler retrieves the "company-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-handler"
//	type: handlers.CompaniesHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.ICompanyManager) ["company-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetCompanyHandler() handlers.CompaniesHandler {
	o, err := c.UnscopedSafeGetCompanyHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// CompanyHandler retrieves the "company-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-handler"
//	type: handlers.CompaniesHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.ICompanyManager) ["company-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetCompanyHandler method.
// If the container can not be retrieved, it panics.
func CompanyHandler(i interface{}) handlers.CompaniesHandler {
	return C(i).GetCompanyHandler()
}

// SafeGetCompanyManager retrieves the "company-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-manager"
//	type: managers.ICompanyManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.ICompanyService) ["company-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetCompanyManager() (managers.ICompanyManager, error) {
	i, err := c.ctn.SafeGet("company-manager")
	if err != nil {
		var eo managers.ICompanyManager
		return eo, err
	}
	o, ok := i.(managers.ICompanyManager)
	if !ok {
		return o, errors.New("could get 'company-manager' because the object could not be cast to managers.ICompanyManager")
	}
	return o, nil
}

// GetCompanyManager retrieves the "company-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-manager"
//	type: managers.ICompanyManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.ICompanyService) ["company-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetCompanyManager() managers.ICompanyManager {
	o, err := c.SafeGetCompanyManager()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCompanyManager retrieves the "company-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-manager"
//	type: managers.ICompanyManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.ICompanyService) ["company-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetCompanyManager() (managers.ICompanyManager, error) {
	i, err := c.ctn.UnscopedSafeGet("company-manager")
	if err != nil {
		var eo managers.ICompanyManager
		return eo, err
	}
	o, ok := i.(managers.ICompanyManager)
	if !ok {
		return o, errors.New("could get 'company-manager' because the object could not be cast to managers.ICompanyManager")
	}
	return o, nil
}

// UnscopedGetCompanyManager retrieves the "company-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-manager"
//	type: managers.ICompanyManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.ICompanyService) ["company-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetCompanyManager() managers.ICompanyManager {
	o, err := c.UnscopedSafeGetCompanyManager()
	if err != nil {
		panic(err)
	}
	return o
}

// CompanyManager retrieves the "company-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-manager"
//	type: managers.ICompanyManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.ICompanyService) ["company-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetCompanyManager method.
// If the container can not be retrieved, it panics.
func CompanyManager(i interface{}) managers.ICompanyManager {
	return C(i).GetCompanyManager()
}

// SafeGetCompanyRepository retrieves the "company-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-repository"
//	type: repositories.ICompanyRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetCompanyRepository() (repositories.ICompanyRepository, error) {
	i, err := c.ctn.SafeGet("company-repository")
	if err != nil {
		var eo repositories.ICompanyRepository
		return eo, err
	}
	o, ok := i.(repositories.ICompanyRepository)
	if !ok {
		return o, errors.New("could get 'company-repository' because the object could not be cast to repositories.ICompanyRepository")
	}
	return o, nil
}

// GetCompanyRepository retrieves the "company-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-repository"
//	type: repositories.ICompanyRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetCompanyRepository() repositories.ICompanyRepository {
	o, err := c.SafeGetCompanyRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCompanyRepository retrieves the "company-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-repository"
//	type: repositories.ICompanyRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetCompanyRepository() (repositories.ICompanyRepository, error) {
	i, err := c.ctn.UnscopedSafeGet("company-repository")
	if err != nil {
		var eo repositories.ICompanyRepository
		return eo, err
	}
	o, ok := i.(repositories.ICompanyRepository)
	if !ok {
		return o, errors.New("could get 'company-repository' because the object could not be cast to repositories.ICompanyRepository")
	}
	return o, nil
}

// UnscopedGetCompanyRepository retrieves the "company-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-repository"
//	type: repositories.ICompanyRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetCompanyRepository() repositories.ICompanyRepository {
	o, err := c.UnscopedSafeGetCompanyRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// CompanyRepository retrieves the "company-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-repository"
//	type: repositories.ICompanyRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetCompanyRepository method.
// If the container can not be retrieved, it panics.
func CompanyRepository(i interface{}) repositories.ICompanyRepository {
	return C(i).GetCompanyRepository()
}

// SafeGetCompanyService retrieves the "company-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-service"
//	type: services.ICompanyService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.ICompanyRepository) ["company-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetCompanyService() (services.ICompanyService, error) {
	i, err := c.ctn.SafeGet("company-service")
	if err != nil {
		var eo services.ICompanyService
		return eo, err
	}
	o, ok := i.(services.ICompanyService)
	if !ok {
		return o, errors.New("could get 'company-service' because the object could not be cast to services.ICompanyService")
	}
	return o, nil
}

// GetCompanyService retrieves the "company-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-service"
//	type: services.ICompanyService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.ICompanyRepository) ["company-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetCompanyService() services.ICompanyService {
	o, err := c.SafeGetCompanyService()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetCompanyService retrieves the "company-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-service"
//	type: services.ICompanyService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.ICompanyRepository) ["company-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetCompanyService() (services.ICompanyService, error) {
	i, err := c.ctn.UnscopedSafeGet("company-service")
	if err != nil {
		var eo services.ICompanyService
		return eo, err
	}
	o, ok := i.(services.ICompanyService)
	if !ok {
		return o, errors.New("could get 'company-service' because the object could not be cast to services.ICompanyService")
	}
	return o, nil
}

// UnscopedGetCompanyService retrieves the "company-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-service"
//	type: services.ICompanyService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.ICompanyRepository) ["company-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetCompanyService() services.ICompanyService {
	o, err := c.UnscopedSafeGetCompanyService()
	if err != nil {
		panic(err)
	}
	return o
}

// CompanyService retrieves the "company-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "company-service"
//	type: services.ICompanyService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.ICompanyRepository) ["company-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetCompanyService method.
// If the container can not be retrieved, it panics.
func CompanyService(i interface{}) services.ICompanyService {
	return C(i).GetCompanyService()
}

// SafeGetConstructionFeatureRepository retrieves the "construction-feature-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-feature-repository"
//	type: repositories.IConstructionFeatureRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetConstructionFeatureRepository() (repositories.IConstructionFeatureRepository, error) {
	i, err := c.ctn.SafeGet("construction-feature-repository")
	if err != nil {
		var eo repositories.IConstructionFeatureRepository
		return eo, err
	}
	o, ok := i.(repositories.IConstructionFeatureRepository)
	if !ok {
		return o, errors.New("could get 'construction-feature-repository' because the object could not be cast to repositories.IConstructionFeatureRepository")
	}
	return o, nil
}

// GetConstructionFeatureRepository retrieves the "construction-feature-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-feature-repository"
//	type: repositories.IConstructionFeatureRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetConstructionFeatureRepository() repositories.IConstructionFeatureRepository {
	o, err := c.SafeGetConstructionFeatureRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetConstructionFeatureRepository retrieves the "construction-feature-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-feature-repository"
//	type: repositories.IConstructionFeatureRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetConstructionFeatureRepository() (repositories.IConstructionFeatureRepository, error) {
	i, err := c.ctn.UnscopedSafeGet("construction-feature-repository")
	if err != nil {
		var eo repositories.IConstructionFeatureRepository
		return eo, err
	}
	o, ok := i.(repositories.IConstructionFeatureRepository)
	if !ok {
		return o, errors.New("could get 'construction-feature-repository' because the object could not be cast to repositories.IConstructionFeatureRepository")
	}
	return o, nil
}

// UnscopedGetConstructionFeatureRepository retrieves the "construction-feature-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-feature-repository"
//	type: repositories.IConstructionFeatureRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetConstructionFeatureRepository() repositories.IConstructionFeatureRepository {
	o, err := c.UnscopedSafeGetConstructionFeatureRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// ConstructionFeatureRepository retrieves the "construction-feature-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-feature-repository"
//	type: repositories.IConstructionFeatureRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetConstructionFeatureRepository method.
// If the container can not be retrieved, it panics.
func ConstructionFeatureRepository(i interface{}) repositories.IConstructionFeatureRepository {
	return C(i).GetConstructionFeatureRepository()
}

// SafeGetConstructionFeatureService retrieves the "construction-feature-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-feature-service"
//	type: services.IConstructionFeatureService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionFeatureRepository) ["construction-feature-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetConstructionFeatureService() (services.IConstructionFeatureService, error) {
	i, err := c.ctn.SafeGet("construction-feature-service")
	if err != nil {
		var eo services.IConstructionFeatureService
		return eo, err
	}
	o, ok := i.(services.IConstructionFeatureService)
	if !ok {
		return o, errors.New("could get 'construction-feature-service' because the object could not be cast to services.IConstructionFeatureService")
	}
	return o, nil
}

// GetConstructionFeatureService retrieves the "construction-feature-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-feature-service"
//	type: services.IConstructionFeatureService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionFeatureRepository) ["construction-feature-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetConstructionFeatureService() services.IConstructionFeatureService {
	o, err := c.SafeGetConstructionFeatureService()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetConstructionFeatureService retrieves the "construction-feature-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-feature-service"
//	type: services.IConstructionFeatureService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionFeatureRepository) ["construction-feature-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetConstructionFeatureService() (services.IConstructionFeatureService, error) {
	i, err := c.ctn.UnscopedSafeGet("construction-feature-service")
	if err != nil {
		var eo services.IConstructionFeatureService
		return eo, err
	}
	o, ok := i.(services.IConstructionFeatureService)
	if !ok {
		return o, errors.New("could get 'construction-feature-service' because the object could not be cast to services.IConstructionFeatureService")
	}
	return o, nil
}

// UnscopedGetConstructionFeatureService retrieves the "construction-feature-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-feature-service"
//	type: services.IConstructionFeatureService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionFeatureRepository) ["construction-feature-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetConstructionFeatureService() services.IConstructionFeatureService {
	o, err := c.UnscopedSafeGetConstructionFeatureService()
	if err != nil {
		panic(err)
	}
	return o
}

// ConstructionFeatureService retrieves the "construction-feature-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-feature-service"
//	type: services.IConstructionFeatureService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionFeatureRepository) ["construction-feature-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetConstructionFeatureService method.
// If the container can not be retrieved, it panics.
func ConstructionFeatureService(i interface{}) services.IConstructionFeatureService {
	return C(i).GetConstructionFeatureService()
}

// SafeGetConstructionHandler retrieves the "construction-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-handler"
//	type: handlers.ConstructionsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IConstructionManager) ["construction-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetConstructionHandler() (handlers.ConstructionsHandler, error) {
	i, err := c.ctn.SafeGet("construction-handler")
	if err != nil {
		var eo handlers.ConstructionsHandler
		return eo, err
	}
	o, ok := i.(handlers.ConstructionsHandler)
	if !ok {
		return o, errors.New("could get 'construction-handler' because the object could not be cast to handlers.ConstructionsHandler")
	}
	return o, nil
}

// GetConstructionHandler retrieves the "construction-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-handler"
//	type: handlers.ConstructionsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IConstructionManager) ["construction-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetConstructionHandler() handlers.ConstructionsHandler {
	o, err := c.SafeGetConstructionHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetConstructionHandler retrieves the "construction-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-handler"
//	type: handlers.ConstructionsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IConstructionManager) ["construction-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetConstructionHandler() (handlers.ConstructionsHandler, error) {
	i, err := c.ctn.UnscopedSafeGet("construction-handler")
	if err != nil {
		var eo handlers.ConstructionsHandler
		return eo, err
	}
	o, ok := i.(handlers.ConstructionsHandler)
	if !ok {
		return o, errors.New("could get 'construction-handler' because the object could not be cast to handlers.ConstructionsHandler")
	}
	return o, nil
}

// UnscopedGetConstructionHandler retrieves the "construction-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-handler"
//	type: handlers.ConstructionsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IConstructionManager) ["construction-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetConstructionHandler() handlers.ConstructionsHandler {
	o, err := c.UnscopedSafeGetConstructionHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// ConstructionHandler retrieves the "construction-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-handler"
//	type: handlers.ConstructionsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IConstructionManager) ["construction-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetConstructionHandler method.
// If the container can not be retrieved, it panics.
func ConstructionHandler(i interface{}) handlers.ConstructionsHandler {
	return C(i).GetConstructionHandler()
}

// SafeGetConstructionImageRepository retrieves the "construction-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-image-repository"
//	type: repositories.IConstructionImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetConstructionImageRepository() (repositories.IConstructionImageRepository, error) {
	i, err := c.ctn.SafeGet("construction-image-repository")
	if err != nil {
		var eo repositories.IConstructionImageRepository
		return eo, err
	}
	o, ok := i.(repositories.IConstructionImageRepository)
	if !ok {
		return o, errors.New("could get 'construction-image-repository' because the object could not be cast to repositories.IConstructionImageRepository")
	}
	return o, nil
}

// GetConstructionImageRepository retrieves the "construction-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-image-repository"
//	type: repositories.IConstructionImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetConstructionImageRepository() repositories.IConstructionImageRepository {
	o, err := c.SafeGetConstructionImageRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetConstructionImageRepository retrieves the "construction-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-image-repository"
//	type: repositories.IConstructionImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetConstructionImageRepository() (repositories.IConstructionImageRepository, error) {
	i, err := c.ctn.UnscopedSafeGet("construction-image-repository")
	if err != nil {
		var eo repositories.IConstructionImageRepository
		return eo, err
	}
	o, ok := i.(repositories.IConstructionImageRepository)
	if !ok {
		return o, errors.New("could get 'construction-image-repository' because the object could not be cast to repositories.IConstructionImageRepository")
	}
	return o, nil
}

// UnscopedGetConstructionImageRepository retrieves the "construction-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-image-repository"
//	type: repositories.IConstructionImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetConstructionImageRepository() repositories.IConstructionImageRepository {
	o, err := c.UnscopedSafeGetConstructionImageRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// ConstructionImageRepository retrieves the "construction-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-image-repository"
//	type: repositories.IConstructionImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetConstructionImageRepository method.
// If the container can not be retrieved, it panics.
func ConstructionImageRepository(i interface{}) repositories.IConstructionImageRepository {
	return C(i).GetConstructionImageRepository()
}

// SafeGetConstructionImageService retrieves the "construction-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-image-service"
//	type: services.IConstructionImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionImageRepository) ["construction-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetConstructionImageService() (services.IConstructionImageService, error) {
	i, err := c.ctn.SafeGet("construction-image-service")
	if err != nil {
		var eo services.IConstructionImageService
		return eo, err
	}
	o, ok := i.(services.IConstructionImageService)
	if !ok {
		return o, errors.New("could get 'construction-image-service' because the object could not be cast to services.IConstructionImageService")
	}
	return o, nil
}

// GetConstructionImageService retrieves the "construction-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-image-service"
//	type: services.IConstructionImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionImageRepository) ["construction-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetConstructionImageService() services.IConstructionImageService {
	o, err := c.SafeGetConstructionImageService()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetConstructionImageService retrieves the "construction-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-image-service"
//	type: services.IConstructionImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionImageRepository) ["construction-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetConstructionImageService() (services.IConstructionImageService, error) {
	i, err := c.ctn.UnscopedSafeGet("construction-image-service")
	if err != nil {
		var eo services.IConstructionImageService
		return eo, err
	}
	o, ok := i.(services.IConstructionImageService)
	if !ok {
		return o, errors.New("could get 'construction-image-service' because the object could not be cast to services.IConstructionImageService")
	}
	return o, nil
}

// UnscopedGetConstructionImageService retrieves the "construction-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-image-service"
//	type: services.IConstructionImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionImageRepository) ["construction-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetConstructionImageService() services.IConstructionImageService {
	o, err := c.UnscopedSafeGetConstructionImageService()
	if err != nil {
		panic(err)
	}
	return o
}

// ConstructionImageService retrieves the "construction-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-image-service"
//	type: services.IConstructionImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionImageRepository) ["construction-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetConstructionImageService method.
// If the container can not be retrieved, it panics.
func ConstructionImageService(i interface{}) services.IConstructionImageService {
	return C(i).GetConstructionImageService()
}

// SafeGetConstructionManager retrieves the "construction-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-manager"
//	type: managers.IConstructionManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IConstructionService) ["construction-service"]
//		- "1": Service(services.IConstructionFeatureService) ["construction-feature-service"]
//		- "2": Service(services.IConstructionImageService) ["construction-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetConstructionManager() (managers.IConstructionManager, error) {
	i, err := c.ctn.SafeGet("construction-manager")
	if err != nil {
		var eo managers.IConstructionManager
		return eo, err
	}
	o, ok := i.(managers.IConstructionManager)
	if !ok {
		return o, errors.New("could get 'construction-manager' because the object could not be cast to managers.IConstructionManager")
	}
	return o, nil
}

// GetConstructionManager retrieves the "construction-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-manager"
//	type: managers.IConstructionManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IConstructionService) ["construction-service"]
//		- "1": Service(services.IConstructionFeatureService) ["construction-feature-service"]
//		- "2": Service(services.IConstructionImageService) ["construction-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetConstructionManager() managers.IConstructionManager {
	o, err := c.SafeGetConstructionManager()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetConstructionManager retrieves the "construction-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-manager"
//	type: managers.IConstructionManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IConstructionService) ["construction-service"]
//		- "1": Service(services.IConstructionFeatureService) ["construction-feature-service"]
//		- "2": Service(services.IConstructionImageService) ["construction-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetConstructionManager() (managers.IConstructionManager, error) {
	i, err := c.ctn.UnscopedSafeGet("construction-manager")
	if err != nil {
		var eo managers.IConstructionManager
		return eo, err
	}
	o, ok := i.(managers.IConstructionManager)
	if !ok {
		return o, errors.New("could get 'construction-manager' because the object could not be cast to managers.IConstructionManager")
	}
	return o, nil
}

// UnscopedGetConstructionManager retrieves the "construction-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-manager"
//	type: managers.IConstructionManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IConstructionService) ["construction-service"]
//		- "1": Service(services.IConstructionFeatureService) ["construction-feature-service"]
//		- "2": Service(services.IConstructionImageService) ["construction-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetConstructionManager() managers.IConstructionManager {
	o, err := c.UnscopedSafeGetConstructionManager()
	if err != nil {
		panic(err)
	}
	return o
}

// ConstructionManager retrieves the "construction-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-manager"
//	type: managers.IConstructionManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IConstructionService) ["construction-service"]
//		- "1": Service(services.IConstructionFeatureService) ["construction-feature-service"]
//		- "2": Service(services.IConstructionImageService) ["construction-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetConstructionManager method.
// If the container can not be retrieved, it panics.
func ConstructionManager(i interface{}) managers.IConstructionManager {
	return C(i).GetConstructionManager()
}

// SafeGetConstructionRepository retrieves the "construction-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-repository"
//	type: repositories.IConstructionRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetConstructionRepository() (repositories.IConstructionRepository, error) {
	i, err := c.ctn.SafeGet("construction-repository")
	if err != nil {
		var eo repositories.IConstructionRepository
		return eo, err
	}
	o, ok := i.(repositories.IConstructionRepository)
	if !ok {
		return o, errors.New("could get 'construction-repository' because the object could not be cast to repositories.IConstructionRepository")
	}
	return o, nil
}

// GetConstructionRepository retrieves the "construction-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-repository"
//	type: repositories.IConstructionRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetConstructionRepository() repositories.IConstructionRepository {
	o, err := c.SafeGetConstructionRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetConstructionRepository retrieves the "construction-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-repository"
//	type: repositories.IConstructionRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetConstructionRepository() (repositories.IConstructionRepository, error) {
	i, err := c.ctn.UnscopedSafeGet("construction-repository")
	if err != nil {
		var eo repositories.IConstructionRepository
		return eo, err
	}
	o, ok := i.(repositories.IConstructionRepository)
	if !ok {
		return o, errors.New("could get 'construction-repository' because the object could not be cast to repositories.IConstructionRepository")
	}
	return o, nil
}

// UnscopedGetConstructionRepository retrieves the "construction-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-repository"
//	type: repositories.IConstructionRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetConstructionRepository() repositories.IConstructionRepository {
	o, err := c.UnscopedSafeGetConstructionRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// ConstructionRepository retrieves the "construction-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-repository"
//	type: repositories.IConstructionRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetConstructionRepository method.
// If the container can not be retrieved, it panics.
func ConstructionRepository(i interface{}) repositories.IConstructionRepository {
	return C(i).GetConstructionRepository()
}

// SafeGetConstructionService retrieves the "construction-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-service"
//	type: services.IConstructionService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionRepository) ["construction-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetConstructionService() (services.IConstructionService, error) {
	i, err := c.ctn.SafeGet("construction-service")
	if err != nil {
		var eo services.IConstructionService
		return eo, err
	}
	o, ok := i.(services.IConstructionService)
	if !ok {
		return o, errors.New("could get 'construction-service' because the object could not be cast to services.IConstructionService")
	}
	return o, nil
}

// GetConstructionService retrieves the "construction-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-service"
//	type: services.IConstructionService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionRepository) ["construction-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetConstructionService() services.IConstructionService {
	o, err := c.SafeGetConstructionService()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetConstructionService retrieves the "construction-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-service"
//	type: services.IConstructionService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionRepository) ["construction-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetConstructionService() (services.IConstructionService, error) {
	i, err := c.ctn.UnscopedSafeGet("construction-service")
	if err != nil {
		var eo services.IConstructionService
		return eo, err
	}
	o, ok := i.(services.IConstructionService)
	if !ok {
		return o, errors.New("could get 'construction-service' because the object could not be cast to services.IConstructionService")
	}
	return o, nil
}

// UnscopedGetConstructionService retrieves the "construction-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-service"
//	type: services.IConstructionService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionRepository) ["construction-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetConstructionService() services.IConstructionService {
	o, err := c.UnscopedSafeGetConstructionService()
	if err != nil {
		panic(err)
	}
	return o
}

// ConstructionService retrieves the "construction-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "construction-service"
//	type: services.IConstructionService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IConstructionRepository) ["construction-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetConstructionService method.
// If the container can not be retrieved, it panics.
func ConstructionService(i interface{}) services.IConstructionService {
	return C(i).GetConstructionService()
}

// SafeGetExternalConnections retrieves the "external-connections" object from the main scope.
//
// ---------------------------------------------
//
//	name: "external-connections"
//	type: connection.IConnection
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetExternalConnections() (connection.IConnection, error) {
	i, err := c.ctn.SafeGet("external-connections")
	if err != nil {
		var eo connection.IConnection
		return eo, err
	}
	o, ok := i.(connection.IConnection)
	if !ok {
		return o, errors.New("could get 'external-connections' because the object could not be cast to connection.IConnection")
	}
	return o, nil
}

// GetExternalConnections retrieves the "external-connections" object from the main scope.
//
// ---------------------------------------------
//
//	name: "external-connections"
//	type: connection.IConnection
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetExternalConnections() connection.IConnection {
	o, err := c.SafeGetExternalConnections()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetExternalConnections retrieves the "external-connections" object from the main scope.
//
// ---------------------------------------------
//
//	name: "external-connections"
//	type: connection.IConnection
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetExternalConnections() (connection.IConnection, error) {
	i, err := c.ctn.UnscopedSafeGet("external-connections")
	if err != nil {
		var eo connection.IConnection
		return eo, err
	}
	o, ok := i.(connection.IConnection)
	if !ok {
		return o, errors.New("could get 'external-connections' because the object could not be cast to connection.IConnection")
	}
	return o, nil
}

// UnscopedGetExternalConnections retrieves the "external-connections" object from the main scope.
//
// ---------------------------------------------
//
//	name: "external-connections"
//	type: connection.IConnection
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetExternalConnections() connection.IConnection {
	o, err := c.UnscopedSafeGetExternalConnections()
	if err != nil {
		panic(err)
	}
	return o
}

// ExternalConnections retrieves the "external-connections" object from the main scope.
//
// ---------------------------------------------
//
//	name: "external-connections"
//	type: connection.IConnection
//	scope: "main"
//	build: func
//	params: nil
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetExternalConnections method.
// If the container can not be retrieved, it panics.
func ExternalConnections(i interface{}) connection.IConnection {
	return C(i).GetExternalConnections()
}

// SafeGetFeatureHandler retrieves the "feature-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-handler"
//	type: handlers.FeaturesHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IFeatureManager) ["feature-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetFeatureHandler() (handlers.FeaturesHandler, error) {
	i, err := c.ctn.SafeGet("feature-handler")
	if err != nil {
		var eo handlers.FeaturesHandler
		return eo, err
	}
	o, ok := i.(handlers.FeaturesHandler)
	if !ok {
		return o, errors.New("could get 'feature-handler' because the object could not be cast to handlers.FeaturesHandler")
	}
	return o, nil
}

// GetFeatureHandler retrieves the "feature-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-handler"
//	type: handlers.FeaturesHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IFeatureManager) ["feature-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetFeatureHandler() handlers.FeaturesHandler {
	o, err := c.SafeGetFeatureHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetFeatureHandler retrieves the "feature-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-handler"
//	type: handlers.FeaturesHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IFeatureManager) ["feature-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetFeatureHandler() (handlers.FeaturesHandler, error) {
	i, err := c.ctn.UnscopedSafeGet("feature-handler")
	if err != nil {
		var eo handlers.FeaturesHandler
		return eo, err
	}
	o, ok := i.(handlers.FeaturesHandler)
	if !ok {
		return o, errors.New("could get 'feature-handler' because the object could not be cast to handlers.FeaturesHandler")
	}
	return o, nil
}

// UnscopedGetFeatureHandler retrieves the "feature-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-handler"
//	type: handlers.FeaturesHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IFeatureManager) ["feature-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetFeatureHandler() handlers.FeaturesHandler {
	o, err := c.UnscopedSafeGetFeatureHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// FeatureHandler retrieves the "feature-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-handler"
//	type: handlers.FeaturesHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IFeatureManager) ["feature-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetFeatureHandler method.
// If the container can not be retrieved, it panics.
func FeatureHandler(i interface{}) handlers.FeaturesHandler {
	return C(i).GetFeatureHandler()
}

// SafeGetFeatureManager retrieves the "feature-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-manager"
//	type: managers.IFeatureManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IFeatureService) ["feature-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetFeatureManager() (managers.IFeatureManager, error) {
	i, err := c.ctn.SafeGet("feature-manager")
	if err != nil {
		var eo managers.IFeatureManager
		return eo, err
	}
	o, ok := i.(managers.IFeatureManager)
	if !ok {
		return o, errors.New("could get 'feature-manager' because the object could not be cast to managers.IFeatureManager")
	}
	return o, nil
}

// GetFeatureManager retrieves the "feature-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-manager"
//	type: managers.IFeatureManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IFeatureService) ["feature-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetFeatureManager() managers.IFeatureManager {
	o, err := c.SafeGetFeatureManager()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetFeatureManager retrieves the "feature-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-manager"
//	type: managers.IFeatureManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IFeatureService) ["feature-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetFeatureManager() (managers.IFeatureManager, error) {
	i, err := c.ctn.UnscopedSafeGet("feature-manager")
	if err != nil {
		var eo managers.IFeatureManager
		return eo, err
	}
	o, ok := i.(managers.IFeatureManager)
	if !ok {
		return o, errors.New("could get 'feature-manager' because the object could not be cast to managers.IFeatureManager")
	}
	return o, nil
}

// UnscopedGetFeatureManager retrieves the "feature-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-manager"
//	type: managers.IFeatureManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IFeatureService) ["feature-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetFeatureManager() managers.IFeatureManager {
	o, err := c.UnscopedSafeGetFeatureManager()
	if err != nil {
		panic(err)
	}
	return o
}

// FeatureManager retrieves the "feature-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-manager"
//	type: managers.IFeatureManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IFeatureService) ["feature-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetFeatureManager method.
// If the container can not be retrieved, it panics.
func FeatureManager(i interface{}) managers.IFeatureManager {
	return C(i).GetFeatureManager()
}

// SafeGetFeatureRepository retrieves the "feature-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-repository"
//	type: repositories.IFeatureRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetFeatureRepository() (repositories.IFeatureRepository, error) {
	i, err := c.ctn.SafeGet("feature-repository")
	if err != nil {
		var eo repositories.IFeatureRepository
		return eo, err
	}
	o, ok := i.(repositories.IFeatureRepository)
	if !ok {
		return o, errors.New("could get 'feature-repository' because the object could not be cast to repositories.IFeatureRepository")
	}
	return o, nil
}

// GetFeatureRepository retrieves the "feature-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-repository"
//	type: repositories.IFeatureRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetFeatureRepository() repositories.IFeatureRepository {
	o, err := c.SafeGetFeatureRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetFeatureRepository retrieves the "feature-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-repository"
//	type: repositories.IFeatureRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetFeatureRepository() (repositories.IFeatureRepository, error) {
	i, err := c.ctn.UnscopedSafeGet("feature-repository")
	if err != nil {
		var eo repositories.IFeatureRepository
		return eo, err
	}
	o, ok := i.(repositories.IFeatureRepository)
	if !ok {
		return o, errors.New("could get 'feature-repository' because the object could not be cast to repositories.IFeatureRepository")
	}
	return o, nil
}

// UnscopedGetFeatureRepository retrieves the "feature-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-repository"
//	type: repositories.IFeatureRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetFeatureRepository() repositories.IFeatureRepository {
	o, err := c.UnscopedSafeGetFeatureRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// FeatureRepository retrieves the "feature-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-repository"
//	type: repositories.IFeatureRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetFeatureRepository method.
// If the container can not be retrieved, it panics.
func FeatureRepository(i interface{}) repositories.IFeatureRepository {
	return C(i).GetFeatureRepository()
}

// SafeGetFeatureService retrieves the "feature-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-service"
//	type: services.IFeatureService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IFeatureRepository) ["feature-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetFeatureService() (services.IFeatureService, error) {
	i, err := c.ctn.SafeGet("feature-service")
	if err != nil {
		var eo services.IFeatureService
		return eo, err
	}
	o, ok := i.(services.IFeatureService)
	if !ok {
		return o, errors.New("could get 'feature-service' because the object could not be cast to services.IFeatureService")
	}
	return o, nil
}

// GetFeatureService retrieves the "feature-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-service"
//	type: services.IFeatureService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IFeatureRepository) ["feature-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetFeatureService() services.IFeatureService {
	o, err := c.SafeGetFeatureService()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetFeatureService retrieves the "feature-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-service"
//	type: services.IFeatureService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IFeatureRepository) ["feature-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetFeatureService() (services.IFeatureService, error) {
	i, err := c.ctn.UnscopedSafeGet("feature-service")
	if err != nil {
		var eo services.IFeatureService
		return eo, err
	}
	o, ok := i.(services.IFeatureService)
	if !ok {
		return o, errors.New("could get 'feature-service' because the object could not be cast to services.IFeatureService")
	}
	return o, nil
}

// UnscopedGetFeatureService retrieves the "feature-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-service"
//	type: services.IFeatureService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IFeatureRepository) ["feature-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetFeatureService() services.IFeatureService {
	o, err := c.UnscopedSafeGetFeatureService()
	if err != nil {
		panic(err)
	}
	return o
}

// FeatureService retrieves the "feature-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "feature-service"
//	type: services.IFeatureService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IFeatureRepository) ["feature-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetFeatureService method.
// If the container can not be retrieved, it panics.
func FeatureService(i interface{}) services.IFeatureService {
	return C(i).GetFeatureService()
}

// SafeGetProductHandler retrieves the "product-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-handler"
//	type: handlers.ProductsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IProductManager) ["product-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetProductHandler() (handlers.ProductsHandler, error) {
	i, err := c.ctn.SafeGet("product-handler")
	if err != nil {
		var eo handlers.ProductsHandler
		return eo, err
	}
	o, ok := i.(handlers.ProductsHandler)
	if !ok {
		return o, errors.New("could get 'product-handler' because the object could not be cast to handlers.ProductsHandler")
	}
	return o, nil
}

// GetProductHandler retrieves the "product-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-handler"
//	type: handlers.ProductsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IProductManager) ["product-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetProductHandler() handlers.ProductsHandler {
	o, err := c.SafeGetProductHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetProductHandler retrieves the "product-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-handler"
//	type: handlers.ProductsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IProductManager) ["product-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetProductHandler() (handlers.ProductsHandler, error) {
	i, err := c.ctn.UnscopedSafeGet("product-handler")
	if err != nil {
		var eo handlers.ProductsHandler
		return eo, err
	}
	o, ok := i.(handlers.ProductsHandler)
	if !ok {
		return o, errors.New("could get 'product-handler' because the object could not be cast to handlers.ProductsHandler")
	}
	return o, nil
}

// UnscopedGetProductHandler retrieves the "product-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-handler"
//	type: handlers.ProductsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IProductManager) ["product-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetProductHandler() handlers.ProductsHandler {
	o, err := c.UnscopedSafeGetProductHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// ProductHandler retrieves the "product-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-handler"
//	type: handlers.ProductsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IProductManager) ["product-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetProductHandler method.
// If the container can not be retrieved, it panics.
func ProductHandler(i interface{}) handlers.ProductsHandler {
	return C(i).GetProductHandler()
}

// SafeGetProductImageRepository retrieves the "product-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-image-repository"
//	type: repositories.IProductImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetProductImageRepository() (repositories.IProductImageRepository, error) {
	i, err := c.ctn.SafeGet("product-image-repository")
	if err != nil {
		var eo repositories.IProductImageRepository
		return eo, err
	}
	o, ok := i.(repositories.IProductImageRepository)
	if !ok {
		return o, errors.New("could get 'product-image-repository' because the object could not be cast to repositories.IProductImageRepository")
	}
	return o, nil
}

// GetProductImageRepository retrieves the "product-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-image-repository"
//	type: repositories.IProductImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetProductImageRepository() repositories.IProductImageRepository {
	o, err := c.SafeGetProductImageRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetProductImageRepository retrieves the "product-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-image-repository"
//	type: repositories.IProductImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetProductImageRepository() (repositories.IProductImageRepository, error) {
	i, err := c.ctn.UnscopedSafeGet("product-image-repository")
	if err != nil {
		var eo repositories.IProductImageRepository
		return eo, err
	}
	o, ok := i.(repositories.IProductImageRepository)
	if !ok {
		return o, errors.New("could get 'product-image-repository' because the object could not be cast to repositories.IProductImageRepository")
	}
	return o, nil
}

// UnscopedGetProductImageRepository retrieves the "product-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-image-repository"
//	type: repositories.IProductImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetProductImageRepository() repositories.IProductImageRepository {
	o, err := c.UnscopedSafeGetProductImageRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// ProductImageRepository retrieves the "product-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-image-repository"
//	type: repositories.IProductImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetProductImageRepository method.
// If the container can not be retrieved, it panics.
func ProductImageRepository(i interface{}) repositories.IProductImageRepository {
	return C(i).GetProductImageRepository()
}

// SafeGetProductImageService retrieves the "product-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-image-service"
//	type: services.IProductImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IProductImageRepository) ["product-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetProductImageService() (services.IProductImageService, error) {
	i, err := c.ctn.SafeGet("product-image-service")
	if err != nil {
		var eo services.IProductImageService
		return eo, err
	}
	o, ok := i.(services.IProductImageService)
	if !ok {
		return o, errors.New("could get 'product-image-service' because the object could not be cast to services.IProductImageService")
	}
	return o, nil
}

// GetProductImageService retrieves the "product-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-image-service"
//	type: services.IProductImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IProductImageRepository) ["product-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetProductImageService() services.IProductImageService {
	o, err := c.SafeGetProductImageService()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetProductImageService retrieves the "product-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-image-service"
//	type: services.IProductImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IProductImageRepository) ["product-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetProductImageService() (services.IProductImageService, error) {
	i, err := c.ctn.UnscopedSafeGet("product-image-service")
	if err != nil {
		var eo services.IProductImageService
		return eo, err
	}
	o, ok := i.(services.IProductImageService)
	if !ok {
		return o, errors.New("could get 'product-image-service' because the object could not be cast to services.IProductImageService")
	}
	return o, nil
}

// UnscopedGetProductImageService retrieves the "product-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-image-service"
//	type: services.IProductImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IProductImageRepository) ["product-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetProductImageService() services.IProductImageService {
	o, err := c.UnscopedSafeGetProductImageService()
	if err != nil {
		panic(err)
	}
	return o
}

// ProductImageService retrieves the "product-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-image-service"
//	type: services.IProductImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IProductImageRepository) ["product-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetProductImageService method.
// If the container can not be retrieved, it panics.
func ProductImageService(i interface{}) services.IProductImageService {
	return C(i).GetProductImageService()
}

// SafeGetProductManager retrieves the "product-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-manager"
//	type: managers.IProductManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IProductService) ["product-service"]
//		- "1": Service(services.IProductImageService) ["product-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetProductManager() (managers.IProductManager, error) {
	i, err := c.ctn.SafeGet("product-manager")
	if err != nil {
		var eo managers.IProductManager
		return eo, err
	}
	o, ok := i.(managers.IProductManager)
	if !ok {
		return o, errors.New("could get 'product-manager' because the object could not be cast to managers.IProductManager")
	}
	return o, nil
}

// GetProductManager retrieves the "product-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-manager"
//	type: managers.IProductManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IProductService) ["product-service"]
//		- "1": Service(services.IProductImageService) ["product-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetProductManager() managers.IProductManager {
	o, err := c.SafeGetProductManager()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetProductManager retrieves the "product-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-manager"
//	type: managers.IProductManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IProductService) ["product-service"]
//		- "1": Service(services.IProductImageService) ["product-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetProductManager() (managers.IProductManager, error) {
	i, err := c.ctn.UnscopedSafeGet("product-manager")
	if err != nil {
		var eo managers.IProductManager
		return eo, err
	}
	o, ok := i.(managers.IProductManager)
	if !ok {
		return o, errors.New("could get 'product-manager' because the object could not be cast to managers.IProductManager")
	}
	return o, nil
}

// UnscopedGetProductManager retrieves the "product-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-manager"
//	type: managers.IProductManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IProductService) ["product-service"]
//		- "1": Service(services.IProductImageService) ["product-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetProductManager() managers.IProductManager {
	o, err := c.UnscopedSafeGetProductManager()
	if err != nil {
		panic(err)
	}
	return o
}

// ProductManager retrieves the "product-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-manager"
//	type: managers.IProductManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IProductService) ["product-service"]
//		- "1": Service(services.IProductImageService) ["product-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetProductManager method.
// If the container can not be retrieved, it panics.
func ProductManager(i interface{}) managers.IProductManager {
	return C(i).GetProductManager()
}

// SafeGetProductRepository retrieves the "product-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-repository"
//	type: repositories.IProductRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetProductRepository() (repositories.IProductRepository, error) {
	i, err := c.ctn.SafeGet("product-repository")
	if err != nil {
		var eo repositories.IProductRepository
		return eo, err
	}
	o, ok := i.(repositories.IProductRepository)
	if !ok {
		return o, errors.New("could get 'product-repository' because the object could not be cast to repositories.IProductRepository")
	}
	return o, nil
}

// GetProductRepository retrieves the "product-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-repository"
//	type: repositories.IProductRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetProductRepository() repositories.IProductRepository {
	o, err := c.SafeGetProductRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetProductRepository retrieves the "product-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-repository"
//	type: repositories.IProductRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetProductRepository() (repositories.IProductRepository, error) {
	i, err := c.ctn.UnscopedSafeGet("product-repository")
	if err != nil {
		var eo repositories.IProductRepository
		return eo, err
	}
	o, ok := i.(repositories.IProductRepository)
	if !ok {
		return o, errors.New("could get 'product-repository' because the object could not be cast to repositories.IProductRepository")
	}
	return o, nil
}

// UnscopedGetProductRepository retrieves the "product-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-repository"
//	type: repositories.IProductRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetProductRepository() repositories.IProductRepository {
	o, err := c.UnscopedSafeGetProductRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// ProductRepository retrieves the "product-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-repository"
//	type: repositories.IProductRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetProductRepository method.
// If the container can not be retrieved, it panics.
func ProductRepository(i interface{}) repositories.IProductRepository {
	return C(i).GetProductRepository()
}

// SafeGetProductService retrieves the "product-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-service"
//	type: services.IProductService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IProductRepository) ["product-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetProductService() (services.IProductService, error) {
	i, err := c.ctn.SafeGet("product-service")
	if err != nil {
		var eo services.IProductService
		return eo, err
	}
	o, ok := i.(services.IProductService)
	if !ok {
		return o, errors.New("could get 'product-service' because the object could not be cast to services.IProductService")
	}
	return o, nil
}

// GetProductService retrieves the "product-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-service"
//	type: services.IProductService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IProductRepository) ["product-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetProductService() services.IProductService {
	o, err := c.SafeGetProductService()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetProductService retrieves the "product-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-service"
//	type: services.IProductService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IProductRepository) ["product-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetProductService() (services.IProductService, error) {
	i, err := c.ctn.UnscopedSafeGet("product-service")
	if err != nil {
		var eo services.IProductService
		return eo, err
	}
	o, ok := i.(services.IProductService)
	if !ok {
		return o, errors.New("could get 'product-service' because the object could not be cast to services.IProductService")
	}
	return o, nil
}

// UnscopedGetProductService retrieves the "product-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-service"
//	type: services.IProductService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IProductRepository) ["product-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetProductService() services.IProductService {
	o, err := c.UnscopedSafeGetProductService()
	if err != nil {
		panic(err)
	}
	return o
}

// ProductService retrieves the "product-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "product-service"
//	type: services.IProductService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IProductRepository) ["product-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetProductService method.
// If the container can not be retrieved, it panics.
func ProductService(i interface{}) services.IProductService {
	return C(i).GetProductService()
}

// SafeGetVariantHandler retrieves the "variant-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-handler"
//	type: handlers.VariantsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IVariantManager) ["variant-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetVariantHandler() (handlers.VariantsHandler, error) {
	i, err := c.ctn.SafeGet("variant-handler")
	if err != nil {
		var eo handlers.VariantsHandler
		return eo, err
	}
	o, ok := i.(handlers.VariantsHandler)
	if !ok {
		return o, errors.New("could get 'variant-handler' because the object could not be cast to handlers.VariantsHandler")
	}
	return o, nil
}

// GetVariantHandler retrieves the "variant-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-handler"
//	type: handlers.VariantsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IVariantManager) ["variant-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetVariantHandler() handlers.VariantsHandler {
	o, err := c.SafeGetVariantHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetVariantHandler retrieves the "variant-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-handler"
//	type: handlers.VariantsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IVariantManager) ["variant-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetVariantHandler() (handlers.VariantsHandler, error) {
	i, err := c.ctn.UnscopedSafeGet("variant-handler")
	if err != nil {
		var eo handlers.VariantsHandler
		return eo, err
	}
	o, ok := i.(handlers.VariantsHandler)
	if !ok {
		return o, errors.New("could get 'variant-handler' because the object could not be cast to handlers.VariantsHandler")
	}
	return o, nil
}

// UnscopedGetVariantHandler retrieves the "variant-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-handler"
//	type: handlers.VariantsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IVariantManager) ["variant-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetVariantHandler() handlers.VariantsHandler {
	o, err := c.UnscopedSafeGetVariantHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// VariantHandler retrieves the "variant-handler" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-handler"
//	type: handlers.VariantsHandler
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(managers.IVariantManager) ["variant-manager"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetVariantHandler method.
// If the container can not be retrieved, it panics.
func VariantHandler(i interface{}) handlers.VariantsHandler {
	return C(i).GetVariantHandler()
}

// SafeGetVariantImageRepository retrieves the "variant-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-image-repository"
//	type: repositories.IVariantImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetVariantImageRepository() (repositories.IVariantImageRepository, error) {
	i, err := c.ctn.SafeGet("variant-image-repository")
	if err != nil {
		var eo repositories.IVariantImageRepository
		return eo, err
	}
	o, ok := i.(repositories.IVariantImageRepository)
	if !ok {
		return o, errors.New("could get 'variant-image-repository' because the object could not be cast to repositories.IVariantImageRepository")
	}
	return o, nil
}

// GetVariantImageRepository retrieves the "variant-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-image-repository"
//	type: repositories.IVariantImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetVariantImageRepository() repositories.IVariantImageRepository {
	o, err := c.SafeGetVariantImageRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetVariantImageRepository retrieves the "variant-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-image-repository"
//	type: repositories.IVariantImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetVariantImageRepository() (repositories.IVariantImageRepository, error) {
	i, err := c.ctn.UnscopedSafeGet("variant-image-repository")
	if err != nil {
		var eo repositories.IVariantImageRepository
		return eo, err
	}
	o, ok := i.(repositories.IVariantImageRepository)
	if !ok {
		return o, errors.New("could get 'variant-image-repository' because the object could not be cast to repositories.IVariantImageRepository")
	}
	return o, nil
}

// UnscopedGetVariantImageRepository retrieves the "variant-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-image-repository"
//	type: repositories.IVariantImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetVariantImageRepository() repositories.IVariantImageRepository {
	o, err := c.UnscopedSafeGetVariantImageRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// VariantImageRepository retrieves the "variant-image-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-image-repository"
//	type: repositories.IVariantImageRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetVariantImageRepository method.
// If the container can not be retrieved, it panics.
func VariantImageRepository(i interface{}) repositories.IVariantImageRepository {
	return C(i).GetVariantImageRepository()
}

// SafeGetVariantImageService retrieves the "variant-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-image-service"
//	type: services.IVariantImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IVariantImageRepository) ["variant-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetVariantImageService() (services.IVariantImageService, error) {
	i, err := c.ctn.SafeGet("variant-image-service")
	if err != nil {
		var eo services.IVariantImageService
		return eo, err
	}
	o, ok := i.(services.IVariantImageService)
	if !ok {
		return o, errors.New("could get 'variant-image-service' because the object could not be cast to services.IVariantImageService")
	}
	return o, nil
}

// GetVariantImageService retrieves the "variant-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-image-service"
//	type: services.IVariantImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IVariantImageRepository) ["variant-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetVariantImageService() services.IVariantImageService {
	o, err := c.SafeGetVariantImageService()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetVariantImageService retrieves the "variant-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-image-service"
//	type: services.IVariantImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IVariantImageRepository) ["variant-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetVariantImageService() (services.IVariantImageService, error) {
	i, err := c.ctn.UnscopedSafeGet("variant-image-service")
	if err != nil {
		var eo services.IVariantImageService
		return eo, err
	}
	o, ok := i.(services.IVariantImageService)
	if !ok {
		return o, errors.New("could get 'variant-image-service' because the object could not be cast to services.IVariantImageService")
	}
	return o, nil
}

// UnscopedGetVariantImageService retrieves the "variant-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-image-service"
//	type: services.IVariantImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IVariantImageRepository) ["variant-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetVariantImageService() services.IVariantImageService {
	o, err := c.UnscopedSafeGetVariantImageService()
	if err != nil {
		panic(err)
	}
	return o
}

// VariantImageService retrieves the "variant-image-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-image-service"
//	type: services.IVariantImageService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IVariantImageRepository) ["variant-image-repository"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetVariantImageService method.
// If the container can not be retrieved, it panics.
func VariantImageService(i interface{}) services.IVariantImageService {
	return C(i).GetVariantImageService()
}

// SafeGetVariantManager retrieves the "variant-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-manager"
//	type: managers.IVariantManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IVariantService) ["variant-service"]
//		- "1": Service(services.IFeatureService) ["feature-service"]
//		- "2": Service(services.IProductService) ["product-service"]
//		- "3": Service(services.IVariantImageService) ["variant-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetVariantManager() (managers.IVariantManager, error) {
	i, err := c.ctn.SafeGet("variant-manager")
	if err != nil {
		var eo managers.IVariantManager
		return eo, err
	}
	o, ok := i.(managers.IVariantManager)
	if !ok {
		return o, errors.New("could get 'variant-manager' because the object could not be cast to managers.IVariantManager")
	}
	return o, nil
}

// GetVariantManager retrieves the "variant-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-manager"
//	type: managers.IVariantManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IVariantService) ["variant-service"]
//		- "1": Service(services.IFeatureService) ["feature-service"]
//		- "2": Service(services.IProductService) ["product-service"]
//		- "3": Service(services.IVariantImageService) ["variant-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetVariantManager() managers.IVariantManager {
	o, err := c.SafeGetVariantManager()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetVariantManager retrieves the "variant-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-manager"
//	type: managers.IVariantManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IVariantService) ["variant-service"]
//		- "1": Service(services.IFeatureService) ["feature-service"]
//		- "2": Service(services.IProductService) ["product-service"]
//		- "3": Service(services.IVariantImageService) ["variant-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetVariantManager() (managers.IVariantManager, error) {
	i, err := c.ctn.UnscopedSafeGet("variant-manager")
	if err != nil {
		var eo managers.IVariantManager
		return eo, err
	}
	o, ok := i.(managers.IVariantManager)
	if !ok {
		return o, errors.New("could get 'variant-manager' because the object could not be cast to managers.IVariantManager")
	}
	return o, nil
}

// UnscopedGetVariantManager retrieves the "variant-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-manager"
//	type: managers.IVariantManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IVariantService) ["variant-service"]
//		- "1": Service(services.IFeatureService) ["feature-service"]
//		- "2": Service(services.IProductService) ["product-service"]
//		- "3": Service(services.IVariantImageService) ["variant-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetVariantManager() managers.IVariantManager {
	o, err := c.UnscopedSafeGetVariantManager()
	if err != nil {
		panic(err)
	}
	return o
}

// VariantManager retrieves the "variant-manager" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-manager"
//	type: managers.IVariantManager
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(services.IVariantService) ["variant-service"]
//		- "1": Service(services.IFeatureService) ["feature-service"]
//		- "2": Service(services.IProductService) ["product-service"]
//		- "3": Service(services.IVariantImageService) ["variant-image-service"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetVariantManager method.
// If the container can not be retrieved, it panics.
func VariantManager(i interface{}) managers.IVariantManager {
	return C(i).GetVariantManager()
}

// SafeGetVariantRepository retrieves the "variant-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-repository"
//	type: repositories.IVariantRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetVariantRepository() (repositories.IVariantRepository, error) {
	i, err := c.ctn.SafeGet("variant-repository")
	if err != nil {
		var eo repositories.IVariantRepository
		return eo, err
	}
	o, ok := i.(repositories.IVariantRepository)
	if !ok {
		return o, errors.New("could get 'variant-repository' because the object could not be cast to repositories.IVariantRepository")
	}
	return o, nil
}

// GetVariantRepository retrieves the "variant-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-repository"
//	type: repositories.IVariantRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetVariantRepository() repositories.IVariantRepository {
	o, err := c.SafeGetVariantRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetVariantRepository retrieves the "variant-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-repository"
//	type: repositories.IVariantRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetVariantRepository() (repositories.IVariantRepository, error) {
	i, err := c.ctn.UnscopedSafeGet("variant-repository")
	if err != nil {
		var eo repositories.IVariantRepository
		return eo, err
	}
	o, ok := i.(repositories.IVariantRepository)
	if !ok {
		return o, errors.New("could get 'variant-repository' because the object could not be cast to repositories.IVariantRepository")
	}
	return o, nil
}

// UnscopedGetVariantRepository retrieves the "variant-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-repository"
//	type: repositories.IVariantRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetVariantRepository() repositories.IVariantRepository {
	o, err := c.UnscopedSafeGetVariantRepository()
	if err != nil {
		panic(err)
	}
	return o
}

// VariantRepository retrieves the "variant-repository" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-repository"
//	type: repositories.IVariantRepository
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(connection.IConnection) ["external-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetVariantRepository method.
// If the container can not be retrieved, it panics.
func VariantRepository(i interface{}) repositories.IVariantRepository {
	return C(i).GetVariantRepository()
}

// SafeGetVariantService retrieves the "variant-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-service"
//	type: services.IVariantService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IVariantRepository) ["variant-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it returns an error.
func (c *Container) SafeGetVariantService() (services.IVariantService, error) {
	i, err := c.ctn.SafeGet("variant-service")
	if err != nil {
		var eo services.IVariantService
		return eo, err
	}
	o, ok := i.(services.IVariantService)
	if !ok {
		return o, errors.New("could get 'variant-service' because the object could not be cast to services.IVariantService")
	}
	return o, nil
}

// GetVariantService retrieves the "variant-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-service"
//	type: services.IVariantService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IVariantRepository) ["variant-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// If the object can not be retrieved, it panics.
func (c *Container) GetVariantService() services.IVariantService {
	o, err := c.SafeGetVariantService()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetVariantService retrieves the "variant-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-service"
//	type: services.IVariantService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IVariantRepository) ["variant-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it returns an error.
func (c *Container) UnscopedSafeGetVariantService() (services.IVariantService, error) {
	i, err := c.ctn.UnscopedSafeGet("variant-service")
	if err != nil {
		var eo services.IVariantService
		return eo, err
	}
	o, ok := i.(services.IVariantService)
	if !ok {
		return o, errors.New("could get 'variant-service' because the object could not be cast to services.IVariantService")
	}
	return o, nil
}

// UnscopedGetVariantService retrieves the "variant-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-service"
//	type: services.IVariantService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IVariantRepository) ["variant-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// This method can be called even if main is a sub-scope of the container.
// If the object can not be retrieved, it panics.
func (c *Container) UnscopedGetVariantService() services.IVariantService {
	o, err := c.UnscopedSafeGetVariantService()
	if err != nil {
		panic(err)
	}
	return o
}

// VariantService retrieves the "variant-service" object from the main scope.
//
// ---------------------------------------------
//
//	name: "variant-service"
//	type: services.IVariantService
//	scope: "main"
//	build: func
//	params:
//		- "0": Service(repositories.IVariantRepository) ["variant-repository"]
//		- "1": Service(connection.ICacheDB) ["cache-connections"]
//	unshared: false
//	close: false
//
// ---------------------------------------------
//
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it calls the GetVariantService method.
// If the container can not be retrieved, it panics.
func VariantService(i interface{}) services.IVariantService {
	return C(i).GetVariantService()
}
