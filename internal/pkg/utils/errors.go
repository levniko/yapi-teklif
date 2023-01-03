package utils

// Error messages
const (
	//Variant error messages
	VariantCanNotCreated      = "variant can not created"
	VariantFeatureIsWrongType = "variant feature is wrong type"
	VariantCanNotFound        = "variant can not found"
	VariantCanNotUpdated      = "variant can not updated"
	VariantCanNotDeleted      = "variant can not deleted"

	//Product error messages
	ProductCanNotCreated       = "product can not created"
	ProductCategoryCanNotFound = "product category can not found"
	ProductCanNotUpdated       = "product can not updated"
	ProductCanNotFound         = "product can not found"
	ProductCanNotDeleted       = "product can not deleted"
	ProductSpuMustBeUnique     = "product spu must be unique"

	//Product Images error messages
	ProductImageCanNotUpdated = "product image can not be updated"

	//Construction error messages
	ConstructionCanNotCreated      = "construction can not created"
	ConstructionFeatureIsWrongType = "construction feature is wrong type"
	ConstructionCanNotFound        = "construction can not found"
	ConstructionCanNotUpdated      = "construction can not updated"

	//Variant Image error messages
	VariantImageCanNotUpdated = "variant image can not updated"
	VariantImageCanNotDeleted = "variant image can not deleted"
	VariantImageCanNotCreated = "variant image can not created"
	VariantImageCanNotFound   = "variant image can not found"

	//Construction Image error messages
	ConstructionImageCanNotUpdated = "construction image can not updated"
	ConstructionImageCanNotDeleted = "construction image can not deleted"
	ConstructionImageCanNotFound   = "construction image can not found"
	ConstructionImageCanNotCreated = "construction image can not created"

	//Features error messages
	FeatureNotFound = "feature not found"

	//Company error messages
	CompanyRecordNotFound    = "company record not found"
	EmailOrPasswordIncorrect = "email or password is incorrect"
	EmailAlreadyExist        = "email already exists"
	CompanySaveError         = "company save error"

	//Form error
	FormValidationError = "form validation error"

	//Password error
	PasswordsAreNotSame = " passwords are not the same"

	required            = "This field is required. Please do not empty."
	requiredwithout     = "At least one of the fields is required. Please do not empty %v."
	number              = "%v is not a valid number. Please enter a valid number."
	numeric             = "%v is not a valid numeric number. Please enter a valid numeric number."
	alphanum            = "%v is not a valid alphanumeric character. Please enter a valid alphanumeric character."
	max                 = "The value has exceeded the maximum limit. Maximum Limit: %v"
	min                 = "The value has exceeded the minimum limit. Minumum Limit: %v"
	eqfield             = "This value has to be equal %v. Value: %v"
	email               = "This value is not a valid email. Value: %v"
	defaultErrorMessage = "This value is not valid. Value: %v"
)

// Errors Codes for error messages
const (
	//Variants
	VariantCanNotCreatedCode      = 400
	VariantFeatureIsWrongTypeCode = 401
	VariantCanNotFoundCode        = 402
	VariantCanNotUpdatedCode      = 403
	VariantCanNotDeletedCode      = 404

	//Products
	ProductCanNotCreatedCode       = 300
	ProductCategoryCanNotFoundCode = 301
	ProductCanNotUpdatedCode       = 302
	ProductCanNotFoundCode         = 303
	ProductCanNotDeletedCode       = 304
	ProductSpuMustBeUniqueCode     = 305

	//Product Images
	ProductImageCanNotUpdatedCode = 700

	//Constructions
	ConstructionCanNotCreatedCode      = 900
	ConstructionFeatureIsWrongTypeCode = 901
	ConstructionCanNotFoundCode        = 902
	ConstructionCanNotUpdatedCode      = 903
	ConstructionCanNotDeletedCode      = 904

	//Features
	FeatureNotFoundCode = 500

	//Variant Images
	VariantImageCanNotUpdatedCode = 1000
	VariantImageCanNotDeletedCode = 1001
	VariantImageCanNotCreatedCode = 1002
	VariantImageCanNotFoundCode   = 1003

	//Construction Images
	ConstructionImageCanNotUpdatedCode = 1100
	ConstructionImageCanNotDeletedCode = 1101
	ConstructionImageCanNotCreatedCode = 1102
	ConstructionImageCanNotFoundCode   = 1103
)
