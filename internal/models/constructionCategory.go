package models

import (
	"time"

	"gorm.io/gorm"
)

type ConstructionCategory struct {
	ID          uint                   `json:"id" gorm:"primaryKey;autoIncrement;notNull"`
	Name        string                 `json:"name" gorm:"type:varchar(100)"`
	Parent      *ConstructionCategory  `json:"parent" gorm:"foreignkey:ParentID"`
	ParentID    uint                   `json:"parent_id" gorm:"index;type:int"`
	Description string                 `json:"description" gorm:"type:varchar(255)"`
	Children    []ConstructionCategory `json:"children" gorm:"foreignKey:ID"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `json:"deleted_at"`
}

func (constructionCategory *ConstructionCategory) TableName() string {
	return "construction_categories"
}

func (constructionCategory *ConstructionCategory) BeforeCreate(tx *gorm.DB) error {
	constructionCategory.CreatedAt = time.Now()
	return nil
}

func (constructionCategory *ConstructionCategory) BeforeUpdate(tx *gorm.DB) error {
	constructionCategory.UpdatedAt = time.Now()
	return nil
}
