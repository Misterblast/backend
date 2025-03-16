package entity

type SendOTP struct {
	Email string `json:"email" validate:"required,email"`
}

type SendDeeplink struct {
	Email string `json:"email" validate:"required,email"`
}

type CheckOTP struct {
	OTP string `json:"otp" validate:"required"`
	ID  int32  `json:"id" validate:"required"`
}

type DeeplinkResponse struct {
	UserID    int32  `json:"user_id"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}
