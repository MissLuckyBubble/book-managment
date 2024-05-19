package models

type Book struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	ISBN   string `json:"isbn"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}
