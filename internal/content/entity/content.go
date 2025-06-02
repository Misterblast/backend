package entity

type Content struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	ImgURL  string `json:"img_url"`
	SiteURL string `json:"site_url"`
	Lang    string `json:"lang"`
}
