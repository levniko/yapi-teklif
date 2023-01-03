package forms

type LoginForm struct {
	Email    string `json:"email" validate:"required,email,max=75"`
	Password string `json:"password" validate:"required,max=100"`
}

type CompanyCreateForm struct {
	Name                     string  `json:"name" validate:"required,max=50"`
	CompanyType              string  `json:"company_type" validate:"required,oneof='Anonim Şirketi' 'Şahıs' 'Limited Şirketi' 'Kollektif Şirket' 'Adi Ortaklık' 'Adi Komandit Şirket' 'Sermayesi Paylara Bölünmüş Komandit Şirket' 'Diğer'"`
	WebSite                  *string `json:"web_site" validate:"omitempty"`
	Email                    string  `json:"email" validate:"required"`
	CompanyAuthorizedName    string  `json:"company_authorized_name" validate:"required,max=50"`
	CompanyAuthorizedSurname string  `json:"company_authorized_surname" validate:"required,max=50"`
	IsActive                 bool    `json:"is_active" validate:"omitempty"`
	IsSupplier               bool    `json:"is_supplier" validate:"required_without=IsConstructor,omitempty"`
	IsConstructor            bool    `json:"is_constructor" validate:"required_without=IsSupplier,omitempty"`
	Password                 string  `json:"password" validate:"required,max=100,eqfield=PasswordAgain"`
	PasswordAgain            string  `json:"password_again" validate:"required,max=100"`
}
