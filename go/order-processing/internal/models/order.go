package models

type Order struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}
