package models

type Post struct {
	ID string `json:"id"`
	UserID string `json:"user_id"`
	Text string `json:"text"`
}