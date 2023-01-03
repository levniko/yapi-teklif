package forms

type ConstructionCreateForm struct {
	Name                   string                    `json:"name" validate:"required,min=0"`
	ConstructionCategoryId uint                      `json:"construction_category_id" validate:"required,number"`
	GeographicRegion       string                    `json:"geographic_region" validate:"required,oneof='Marmara' 'Ege' 'İç Anadolu' 'Akdeniz' 'Karadeniz' 'Doğu Anadolu' 'Güneydoğu Anadolu'"`
	Province               string                    `json:"province" validate:"required"`
	District               string                    `json:"district" validate:"required"`
	Stage                  string                    `json:"stage" validate:"required,oneof='Proje' 'Temel' 'Kaba' 'İnce' 'Tamamlandı' 'Beklemede' 'Devam Ediyor' 'Planlanan'"`
	Start                  string                    `json:"start" validate:"required"`
	End                    string                    `json:"end" validate:"required"`
	WebSite                string                    `json:"web_site" validate:"omitempty,url"`
	CostOfProject          float64                   `json:"cost_of_project" validate:"required,number"`
	LandArea               float64                   `json:"land_area" validate:"required,number"`
	ConstructionZone       float64                   `json:"Construction_zone" validate:"required,number"`
	ConstructionImages     []ConstructionImageForm   `json:"construction_images" validate:"omitempty,dive"`
	ConstructionFeatures   []ConstructionFeatureForm `json:"construction_features" validate:"omitempty,dive"`
}

type ConstructionImageForm struct {
	RemoteLink string `json:"remote_link" validate:"required"`
}

type ConstructionFeatureForm struct {
	FeatureID uint   `json:"feature_id" validate:"required,number"`
	Value     string `json:"value" validate:"required"`
}

type ConstructionUpdateForm struct {
	Name                 string                    `json:"name" validate:"required,min=0"`
	GeographicRegion     *string                   `json:"geographic_region" validate:"omitempty,oneof='Marmara' 'Ege' 'İç Anadolu' 'Akdeniz' 'Karadeniz' 'Doğu Anadolu' 'Güneydoğu Anadolu'"`
	Province             *string                   `json:"province" validate:"omitempty"`
	District             *string                   `json:"district" validate:"omitempty"`
	Stage                *string                   `json:"stage" validate:"omitempty,oneof='Proje' 'Temel' 'Kaba' 'İnce' 'Tamamlandı' 'Beklemede' 'Devam Ediyor' 'Planlanan'"`
	Start                *string                   `json:"start" validate:"omitempty,startEndValidator"`
	End                  *string                   `json:"end" validate:"omitempty,startEndValidator"`
	WebSite              *string                   `json:"web_site" validate:"omitempty,url"`
	CostOfProject        *float64                  `json:"cost_of_project" validate:"omitempty,number"`
	LandArea             *float64                  `json:"land_area" validate:"omitempty,number"`
	ConstructionZone     *float64                  `json:"Construction_zone" validate:"omitempty,number"`
	ConstructionImages   []ConstructionImageForm   `json:"construction_images" validate:"omitempty,dive"`
	ConstructionFeatures []ConstructionFeatureForm `json:"construction_features" validate:"omitempty,dive"`
}
