package models

import (
	"time"

	"gorm.io/gorm"
)

type ConstructionFeature struct {
	ID                    uint           `json:"id" gorm:"primaryKey;autoIncrement;notNull"`
	Construction          *Construction  `json:"construction" gorm:"foreignkey:ConstructionID"`
	ConstructionID        uint           `json:"construction_id" gorm:"index;type:int"`
	ConstructionFeature   CFeature       `json:"construction_feature" gorm:"foreignkey:ConstructionFeatureID"`
	ConstructionFeatureID uint           `json:"construction_feature_id" gorm:"index;type:int"`
	Value                 string         `json:"value" gorm:"type:varchar(128)"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at"`
}

func (constructionFeature *ConstructionFeature) TableName() string {
	return "construction_features"
}

func (constructionFeature *ConstructionFeature) BeforeCreate(tx *gorm.DB) error {
	constructionFeature.CreatedAt = time.Now()
	return nil
}

func (constructionFeature *ConstructionFeature) BeforeUpdate(tx *gorm.DB) error {
	constructionFeature.UpdatedAt = time.Now()
	return nil
}
