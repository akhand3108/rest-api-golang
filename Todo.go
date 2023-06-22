package main

type Todo struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	ID    int    `json:"id"`
}
