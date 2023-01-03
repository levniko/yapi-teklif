package repositories

type IProductFeatureRepository interface {
}

type ProductFeatureRepository struct {
}

func GetProductFeatureRepository() *ProductFeatureRepository {
	return &ProductFeatureRepository{}
}
