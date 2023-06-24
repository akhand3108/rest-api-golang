package models

type Todo struct {
	Title  string `json:"title"`
	ID     int    `json:"id"`
	Done   bool   `json:"done"`
	UserID int    `json:"userId"`
}
