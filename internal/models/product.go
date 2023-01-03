package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID                uint             `json:"id" gorm:"primaryKey;autoIncrement;notNull"`
	CompanyID         uint             `json:"company_id" gorm:"notNull"`
	SPU               string           `json:"spu" gorm:"varchar(11);notNull"`
	Name              string           `json:"name" gorm:"type:varchar(100)"`
	Description       string           `json:"description" gorm:"type:varchar(255)"`
	IsActive          bool             `json:"is_active" gorm:"type:bool"`
	HeroImage         string           `json:"hero_image" gorm:"type:varchar(255)"`
	ProductImages     []ProductImage   `json:"product_images" gorm:"foreingKey:ProductID"`
	ProductCategory   *ProductCategory `json:"product_category" gorm:"foreignkey:ProductCategoryID"`
	ProductCategoryID uint             `json:"category_id" gorm:"index;type:int"`
	Variants          []Variant        `json:"variants" gorm:"foreingKey:ProductID"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         gorm.DeletedAt   `json:"deleted_at"`
}

func (product *Product) TableName() string {
	return "products"
}

func (product *Product) BeforeCreate(tx *gorm.DB) error {
	product.CreatedAt = time.Now()
	return nil
}

func (product *Product) BeforeUpdate(tx *gorm.DB) error {
	product.UpdatedAt = time.Now()
	return nil
}
