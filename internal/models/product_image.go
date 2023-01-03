package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductImage struct {
	ID         uint           `json:"id" gorm:"primaryKey;autoIncrement;notNull"`
	ProductID  uint           `json:"product_id" gorm:"index;type:int"`
	RemoteLink string         `json:"remote_link" gorm:"type:varchar(255)"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

func (productImage *ProductImage) TableName() string {
	return "product_images"
}

func (productImage *ProductImage) BeforeCreate(tx *gorm.DB) error {
	productImage.CreatedAt = time.Now()
	return nil
}

func (productImage *ProductImage) BeforeUpdate(tx *gorm.DB) error {
	productImage.UpdatedAt = time.Now()
	return nil
}
