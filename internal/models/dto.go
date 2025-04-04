package models

type Response[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type PaginationResponse[T any] struct {
	Page  int `json:"page"`
	Total int `json:"total"`
	Limit int `json:"limit"`
	Items []T `json:"items"`
}
