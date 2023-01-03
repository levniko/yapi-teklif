package models

import (
	"time"

	"gorm.io/gorm"
)

type ConstructionImage struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement;notNull"`
	ConstructionID uint           `json:"construction_id" gorm:"index;type:int"`
	RemoteLink     string         `json:"remote_link" gorm:"type:varchar(255)"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at"`
}

func (constructionImage *ConstructionImage) TableName() string {
	return "construction_images"
}

func (constructionImage *ConstructionImage) BeforeCreate(tx *gorm.DB) error {
	constructionImage.CreatedAt = time.Now()
	return nil
}

func (constructionImage *ConstructionImage) BeforeUpdate(tx *gorm.DB) error {
	constructionImage.UpdatedAt = time.Now()
	return nil
}
