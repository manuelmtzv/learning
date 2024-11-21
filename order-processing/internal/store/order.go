package store

import (
	"database/sql"
	"order-processing/internal/models"
)

type OrderStorage struct {
	db *sql.DB
}

func (s *OrderStorage) GetPendingOrders() ([]models.Order, error) {
	return nil, nil
}

func (s *OrderStorage) GetOrder(id string) (models.Order, error) {
	return models.Order{}, nil
}

func (s *OrderStorage) UpdateOrder(order models.Order) error {
	return nil
}
