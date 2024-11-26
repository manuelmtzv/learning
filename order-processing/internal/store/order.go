package store

import (
	"context"
	"database/sql"
	"order-processing/internal/models"
)

type OrderStorage struct {
	db *sql.DB
}

func (s *OrderStorage) CreateOrder(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (status)
		VALUES($1)
		RETURNING id, created_at
	`
	return s.db.QueryRowContext(ctx, query, order.Status).Scan(&order.ID, &order.CreatedAt)
}

func (s *OrderStorage) GetPendingOrders(ctx context.Context) ([]models.Order, error) {
	return nil, nil
}

func (s *OrderStorage) GetOrder(ctx context.Context, id string) (models.Order, error) {
	return models.Order{}, nil
}

func (s *OrderStorage) UpdateOrder(ctx context.Context, order models.Order) error {
	return nil
}
