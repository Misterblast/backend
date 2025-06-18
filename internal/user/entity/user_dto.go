package entity

import "mime/multipart"

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type Register struct {
	Name     string  `json:"name" validate:"required,min=2,max=20"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6,max=20"`
	ImgUrl   *string `json:"img_url,omitempty"`
}

type EditUser struct {
	Name   string  `json:"name" validate:"required,min=2,max=20"`
	Email  string  `json:"email" validate:"required,email"`
	ImgUrl *string `json:"img_url,omitempty"`
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

type ChangePassword struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type RegisterDTO struct {
	Name     string                `form:"name" validate:"required,min=2,max=20"`
	Email    string                `form:"email" validate:"required,email"`
	Password string                `form:"password" validate:"required,min=6,max=20"`
	Img      *multipart.FileHeader `form:"img" validate:"omitempty"`
}

type EditDTO struct {
	Name  string                `form:"name" validate:"required,min=2,max=20"`
	Email string                `form:"email" validate:"required,email"`
	Img   *multipart.FileHeader `form:"img,omitempty"`
}

type UserSummary struct {
	TotalQuizAttempts int     `json:"total_quiz_attempts"`
	TotalTaskAttempts int     `json:"total_task_attempts"`
	AvgQuizScore      float64 `json:"avg_quiz_score"`
	AvgTaskScore      float64 `json:"avg_task_score"`
}

// "https://stg.file.go-assessment.link/file?key=prod/user/profile-img/71/misterblast-1.png"
