package store

import (
	"context"
	"database/sql"
	"order-processing/internal/models"
)

type Storage struct {
	Orders interface {
		CreateOrder(context.Context, *models.Order) error
		GetCreatedOrders(context.Context) ([]*models.Order, error)
		GetOrder(context.Context, string) (*models.Order, error)
		ChangeOrderStatus(context.Context, int, string) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Orders: &OrderStorage{db: db},
	}
}
