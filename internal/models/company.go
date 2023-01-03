package models

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID                       uint           `json:"id" gorm:"primaryKey;autoIncrement;notNull"`
	Name                     string         `json:"name" gorm:"type:varchar(100)"`
	CompanyType              string         `json:"company_type" gorm:"type:varchar(255)"`
	WebSite                  string         `json:"web_site" gorm:"type:varchar(255)"`
	Email                    string         `json:"email" gorm:"type:varchar(255)"`
	CompanyAuthorizedName    string         `json:"company_authorized_name" gorm:"type:varchar(100)"`
	CompanyAuthorizedSurname string         `json:"company_authorized_surname" gorm:"type:varchar(100)"`
	IsActive                 bool           `json:"is_active" gorm:"type:bool"`
	IsSupplier               bool           `json:"is_supplier" gorm:"type:bool"`
	IsConstructor            bool           `json:"is_constructor" gorm:"type:bool"`
	Password                 string         `json:"password" gorm:"type:varchar(100)"`
	PasswordHash             string         `json:"password_hash" gorm:"type:varchar(100)"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `json:"deleted_at"`
}

type AccessDetails struct {
	AccessUuid string
	CompanyID  uint
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}
