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

func (s *OrderStorage) GetCreatedOrders(ctx context.Context) ([]*models.Order, error) {
	query := `
		SELECT id, status, created_at FROM orders
		WHERE status = 'created'
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	orders := []*models.Order{}

	for rows.Next() {
		order := &models.Order{}

		err := rows.Scan(&order.ID, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (s *OrderStorage) ChangeOrderStatus(ctx context.Context, id int, newStatus string) error {
	query := `
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`

	args := []interface{}{newStatus, id}
	_, err := s.db.ExecContext(ctx, query, args...)

	return err
}

func (s *OrderStorage) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	return nil, nil
}
