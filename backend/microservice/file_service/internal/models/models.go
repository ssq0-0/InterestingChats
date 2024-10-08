package models

// Image represents the structure of an image file.
type Image struct {
	ID         int     `json:"id"`
	TypeImage  string  `json:"type_image"`
	Size       float32 `json:"size"`
	URL        string  `json:"url"`
	UploadData string  `json:"upload_data"`
}

// Response represents the structure of an API response.
type Response struct {
	Errors []string `json:"Errors"`
	Data   any      `json:"Data"`
}

// UserFile represents the user's file with a corresponding URL.
type UserFile struct {
	UserID int          `json:"user_id"`
	URL    FileResponse `json:"file_url"`
}

// FileResponse represents the structure of a file URL response.
type FileResponse struct {
	Errors        []string `json:"Errors"`
	TemporaryLink string   `json:"temporary_url"`
	StaticLink    string   `json:"static_url"`
}
