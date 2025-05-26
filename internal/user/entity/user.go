package entity

import "mime/multipart"

type User struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsAdmin    bool   `json:"is_admin"`
	IsVerified bool   `json:"is_verified"`
}

type UserProfileImg struct {
	Img *multipart.FileHeader `json:"img" validate:"required"`
	Key string                `json:"key" validate:"required"`
}

type ImgResponse struct {
	Message string `json:"message"`
	Data    struct {
		URL string `json:"url"`
	} `json:"data"`
}
