package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductFeature struct {
	ID               uint           `json:"id" gorm:"primaryKey;autoIncrement;notNull"`
	Variant          *Variant       `json:"variant" gorm:"foreignkey:VariantID"`
	VariantID        uint           `json:"variant_id" gorm:"index;type:int"`
	ProductFeature   PFeature       `json:"product_feature" gorm:"foreignkey:ProductFeatureID"`
	ProductFeatureID uint           `json:"product_feature_id" gorm:"index;type:int"`
	Value            string         `json:"value" gorm:"type:varchar(128)"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at"`
}

func (productFeature *ProductFeature) TableName() string {
	return "product_features"
}

func (productFeature *ProductFeature) BeforeCreate(tx *gorm.DB) error {
	productFeature.CreatedAt = time.Now()
	return nil
}

func (productFeature *ProductFeature) BeforeUpdate(tx *gorm.DB) error {
	productFeature.UpdatedAt = time.Now()
	return nil
}
