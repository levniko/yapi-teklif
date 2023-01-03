package models

import (
	"time"

	"gorm.io/gorm"
)

type PFeature struct {
	ID                uint             `json:"id" gorm:"primaryKey;autoIncrement;notNull"`
	Name              string           `json:"name" gorm:"type:varchar(100)"`
	Description       string           `json:"description" gorm:"type:varchar(255)"`
	IsRequired        bool             `json:"required" gorm:"type:bool"`
	Type              string           `json:"type" gorm:"type:varchar(255)"`
	ProductCategory   *ProductCategory `json:"product_category" gorm:"foreignKey:ProductCategoryID"`
	ProductCategoryID uint             `json:"product_category_id" gorm:"index;type:int"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         gorm.DeletedAt   `json:"deleted_at"`
}

func (pFeature *PFeature) TableName() string {
	return "p_features"
}

func (pFeature *PFeature) BeforeCreate(tx *gorm.DB) error {
	pFeature.CreatedAt = time.Now()
	return nil
}

func (pFeature *PFeature) BeforeUpdate(tx *gorm.DB) error {
	pFeature.UpdatedAt = time.Now()
	return nil
}
