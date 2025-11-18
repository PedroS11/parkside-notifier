package interfaces

type Flyer struct {
	Url          string    `json:"url"`
	Name         string    `json:"name"`
	Date         string    `json:"date"`
	PreviewImage string    `json:"preview_image"`
	Images       []string  `json:"images"`
	Products     []Product `json:"products"`
}
