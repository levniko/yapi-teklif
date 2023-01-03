package forms

type ProductImageCreateForm struct {
	ProductID  uint   `json:"product_id" validate:"required"`
	RemoteLink string `json:"remote_link" validate:"required"`
}
