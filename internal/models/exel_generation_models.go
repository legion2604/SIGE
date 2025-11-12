package models

type JSONToExcelRequest struct {
	ID        []int    `json:"ID"`
	Path      []string `json:"path"`
	CreatedAt []string `json:"created_at"`
}
