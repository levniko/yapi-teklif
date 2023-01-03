package models

import (
	"time"

	"gorm.io/gorm"
)

type CFeature struct {
	ID                     uint                  `json:"id" gorm:"primaryKey;autoIncrement;notNull"`
	Name                   string                `json:"name" gorm:"type:varchar(100)"`
	Description            string                `json:"description" gorm:"type:varchar(255)"`
	IsRequired             bool                  `json:"is_required" gorm:"type:bool"`
	Type                   string                `json:"type" gorm:"type:varchar(255)"`
	ConstructionCategory   *ConstructionCategory `json:"construction_category" gorm:"foreignKey:ConstructionCategoryID"`
	ConstructionCategoryID uint                  `json:"construction_category_id" gorm:"index;type:int"`
	CreatedAt              time.Time             `json:"created_at"`
	UpdatedAt              time.Time             `json:"updated_at"`
	DeletedAt              gorm.DeletedAt        `json:"deleted_at"`
}

func (cFeature *CFeature) TableName() string {
	return "c_features"
}

func (cFeature *CFeature) BeforeCreate(tx *gorm.DB) error {
	cFeature.CreatedAt = time.Now()
	return nil
}

func (cFeature *CFeature) BeforeUpdate(tx *gorm.DB) error {
	cFeature.UpdatedAt = time.Now()
	return nil
}
