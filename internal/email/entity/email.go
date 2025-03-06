package entity

type SendOTP struct {
	Email string `json:"email" validate:"required,email"`
}

type CheckOTP struct {
	OTP string `json:"otp" validate:"required"`
	ID  int32  `json:"id" validate:"required"`
}
