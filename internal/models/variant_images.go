package models

import (
	"time"

	"gorm.io/gorm"
)

type VariantImage struct {
	ID         uint           `json:"id" gorm:"primaryKey;autoIncrement;notNull"`
	VariantID  uint           `json:"variant_id" gorm:"index;type:int"`
	RemoteLink string         `json:"remote_link" gorm:"type:varchar(255)"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

func (variantImage *VariantImage) TableName() string {
	return "variant_images"
}

func (variantImage *VariantImage) BeforeCreate(tx *gorm.DB) error {
	variantImage.CreatedAt = time.Now()
	return nil
}

func (variantImage *VariantImage) BeforeUpdate(tx *gorm.DB) error {
	variantImage.UpdatedAt = time.Now()
	return nil
}
