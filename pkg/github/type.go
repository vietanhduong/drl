package github

type Release struct {
	TagName string  `json:"tag_name"`
	Id      int     `json:"id"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	ContentType string `json:"content_type"`
}
