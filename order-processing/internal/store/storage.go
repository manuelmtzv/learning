package store

import (
	"database/sql"
	"order-processing/internal/models"
)

type Storage struct {
	Orders interface {
		GetPendingOrders() ([]models.Order, error)
		GetOrder(id string) (models.Order, error)
		UpdateOrder(order models.Order) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Orders: &OrderStorage{db: db},
	}
}
