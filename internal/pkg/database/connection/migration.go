package database

import "github.com/yapi-teklif/internal/models"

func AutoMigrate(connection IConnection) error {
	err := connection.PsqlDB().AutoMigrate(
		&models.Company{},
		&models.Product{},
		&models.Variant{},
		&models.PFeature{},
		&models.ProductCategory{},
		&models.ProductFeature{},
		&models.ProductImage{},
		&models.Construction{},
		&models.CFeature{},
		&models.ConstructionCategory{},
		&models.ConstructionFeature{},
		&models.ConstructionImage{},
	)

	return err
}
