package models

// Response defines the structure of the API response, containing errors and data.
type Response struct {
	Errors []string
	Data   any
}
