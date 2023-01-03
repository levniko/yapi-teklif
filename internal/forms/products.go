package forms

type ProductCreateForm struct {
	Name          string             `json:"name" validate:"required"`
	SPU           string             `json:"spu" validate:"required"`
	IsActive      bool               `json:"is_active" validate:"required"`
	Description   string             `json:"description" validate:"omitempty"`
	HeroImage     string             `json:"hero_image" validate:"omitempty"`
	ProductImages []ProductImageForm `json:"product_images" validate:"omitempty"`
	CategoryID    uint               `json:"category_id" validate:"required"`
}

type ProductImageForm struct {
	RemoteLink string `json:"remote_link" validate:"required"`
}

type ProductUpdateForm struct {
	Name          string             `json:"name" validate:"required"`
	SPU           string             `json:"spu" validate:"required"`
	IsActive      *bool              `json:"is_active" validate:"required"`
	Description   *string            `json:"description" validate:"omitempty"`
	HeroImage     *string            `json:"hero_image" validate:"omitempty"`
	ProductImages []ProductImageForm `json:"product_images" validate:"omitempty"`
}
