package entity

type Content struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	ImgURL  string `json:"img_url"`
	SiteURL string `json:"site_url"`
	Lang    string `json:"lang"`
}

type Author struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	ImgURL      string  `json:"img_url"`
	Description *string `json:"description,omitempty"`
}
