package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Construction struct {
	ID                   uint                  `json:"id" gorm:"primaryKey;autoIncrement;notNull"`
	CompanyID            uint                  `json:"company_id" gorm:"notNull"`
	Name                 string                `json:"name" gorm:"type:varchar(255)"`
	CategoryID           uint                  `json:"category_id" gorm:"index;type:int"`
	GeographicRegion     string                `json:"geographic_region" gorm:"type:varchar(20)"`
	Province             string                `json:"province" gorm:"type:varchar(15)"`
	District             string                `json:"district" gorm:"type:varchar(20)"`
	Stage                string                `json:"stage" gorm:"type:varchar(20)"`
	Start                string                `json:"start" gorm:"type:varchar(25)"`
	End                  string                `json:"end" gorm:"type:varchar(25)"`
	CostOfProject        *decimal.Decimal      `json:"cost_of_project" gorm:"type:decimal(10,2);min=0"`
	LandArea             *decimal.Decimal      `json:"land_area" gorm:"type:decimal(10,2);min=0"`
	ConstructionZone     *decimal.Decimal      `json:"construction_zone" gorm:"type:decimal(10,2);min=0"`
	ConstructionImages   []ConstructionImage   `json:"construction_images" gorm:"foreingKey:ConstructionID"`
	ConstructionFeatures []ConstructionFeature `json:"construction_features"  gorm:"foreignKey:ConstructionID"`
	IsDeleted            bool                  `json:"is_deleted"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
	DeletedAt            gorm.DeletedAt        `json:"deleted_at"`
}

func (construction *Construction) TableName() string {
	return "constructions"
}

func (construction *Construction) BeforeCreate(tx *gorm.DB) error {
	construction.CreatedAt = time.Now()
	return nil
}

func (construction *Construction) BeforeUpdate(tx *gorm.DB) error {
	construction.UpdatedAt = time.Now()
	return nil
}
