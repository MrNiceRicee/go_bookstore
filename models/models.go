package models

type Book struct {
	Id        int    `json:"_id,omitempty"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Genre     string `json:"genre"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Response struct {
	Data interface{} `json:"data,omitempty"`
}
