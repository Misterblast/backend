package entity

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type Register struct {
	Name     string `json:"name" validate:"required,min=2,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type EditUser struct {
	Name     string  `json:"name" validate:"required,min=2,max=20"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password,omitempty" validate:"omitempty,min=6,max=20"`
	ImgUrl   *string `json:"img_url,omitempty"`
}

type ListUser struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	ImgUrl string `json:"img_url"`
}

type DetailUser struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	ImgUrl string `json:"img_url"`
}

type UserJWT struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsAdmin    bool   `json:"is_admin"`
	IsVerified bool   `json:"is_verified"`
}

type LoginResponse struct {
	ID         int32  `json:"id"`
	Email      string `json:"email"`
	IsAdmin    bool   `json:"is_admin"`
	IsVerified bool   `json:"is_verified"`
}

type UserAuth struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	ImgUrl     string `json:"img_url"`
	IsAdmin    bool   `json:"is_admin"`
	IsVerified bool   `json:"is_verified"`
}

// ADMIN

type AdminActivation struct {
	ID    int32  `json:"id"`
	Email string `json:"email" validate:"required,email"`
	OTP   int32  `json:"otp"`
}

type RegisterAdmin struct {
	Name  string `json:"name" validate:"required,min=2,max=20"`
	Email string `json:"email" validate:"required,email"`
}
