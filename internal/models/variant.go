package models

import (
	"time"

	"gorm.io/gorm"
)

type Variant struct {
	ID            uint             `json:"id" gorm:"primaryKey;autoIncrement;notNull"`
	Product       *Product         `json:"product" gorm:"foreignKey:ProductID"`
	ProductID     uint             `json:"product_id" gorm:"index;type:int"`
	SKU           string           `json:"sku" gorm:"varchar(11);notNull"`
	VariantImages []VariantImage   `json:"variant_images" gorm:"foreingKey:VariantID"`
	Features      []ProductFeature `json:"features" gorm:"foreignKey:VariantID"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	DeletedAt     gorm.DeletedAt   `json:"deleted_at"`
}

func (variant *Variant) TableName() string {
	return "variants"
}

func (variant *Variant) BeforeCreate(tx *gorm.DB) error {
	variant.CreatedAt = time.Now()
	return nil
}

func (variant *Variant) BeforeUpdate(tx *gorm.DB) error {
	variant.UpdatedAt = time.Now()
	return nil
}
