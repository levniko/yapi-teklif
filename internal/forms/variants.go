package forms

type VariantCreateForm struct {
	Name          string             `json:"name" validate:"required"`
	SKU           string             `json:"sku" validate:"required"`
	ProductId     int                `json:"product_id" validate:"required,number"`
	Description   string             `json:"description" validate:"required"`
	IsActive      bool               `json:"is_active" validate:"omitempty"`
	VariantImages []VariantImageForm `json:"variant_images" validate:"omitempty,dive"`
	Features      []FeatureForm      `json:"features" validate:"omitempty,dive"`
}

type VariantImageForm struct {
	RemoteLink string `json:"remote_link" validate:"required"`
}

type FeatureForm struct {
	FeatureID uint   `json:"feature_id" validate:"required,number"`
	Value     string `json:"value" validate:"required"`
}

type VariantUpdateForm struct {
	Name          string             `json:"name" validate:"required"`
	SKU           string             `json:"sku" validate:"required"`
	Description   *string            `json:"description" validate:"omitempty"`
	IsActive      *bool              `json:"is_active" validate:"omitempty"`
	VariantImages []VariantImageForm `json:"variant_images" validate:"omitempty,dive"`
	Features      []FeatureForm      `json:"features" validate:"omitempty,dive"`
}
