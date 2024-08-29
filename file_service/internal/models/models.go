package models

type Image struct {
	ID         int     `json:"id"`
	TypeImage  string  `json:"type_image"`
	Size       float32 `json:"size"`
	URL        string  `json:"url"`
	UploadData string  `json:"upload_data"`
}
