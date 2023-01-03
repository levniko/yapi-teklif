package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductCategory struct {
	ID          uint              `json:"id" gorm:"primaryKey;autoIncrement;notNull"`
	Name        string            `json:"name" gorm:"type:varchar(100)"`
	Parent      *ProductCategory  `json:"parent" gorm:"foreignkey:ParentID"`
	ParentID    uint              `json:"parent_id" gorm:"index;type:int"`
	Description string            `json:"description" gorm:"type:varchar(255)"`
	Children    []ProductCategory `json:"children" gorm:"foreignKey:ID"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `json:"deleted_at"`
}

func (productCategory *ProductCategory) TableName() string {
	return "product_categories"
}

func (productCategory *ProductCategory) BeforeCreate(tx *gorm.DB) error {
	productCategory.CreatedAt = time.Now()
	return nil
}

func (productCategory *ProductCategory) BeforeUpdate(tx *gorm.DB) error {
	productCategory.UpdatedAt = time.Now()
	return nil
}
