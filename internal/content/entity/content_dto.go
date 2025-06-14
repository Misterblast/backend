package entity

type CreateAuthorRequest struct {
	Name        string  `json:"name" validate:"required"`
	ImgURL      string  `json:"img_url" validate:"required,url"`
	Description *string `json:"description"`
}

type UpdateAuthorRequest struct {
	Name        string  `json:"name" validate:"required"`
	ImgURL      string  `json:"img_url" validate:"required,url"`
	Description *string `json:"description"`
}
