package entity

type Attachment struct {
	ID   int32  `json:"-"`
	Type string `json:"type"`
	Url  string `json:"url"`
}
