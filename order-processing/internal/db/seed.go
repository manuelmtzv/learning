package db

import (
	"context"
	"database/sql"
	"fmt"
	"order-processing/internal/models"
	"order-processing/internal/store"
)

func Seed(store *store.Storage, db *sql.DB) {
	ctx := context.Background()

	for i := 0; i < 1000; i++ {
		order := &models.Order{
			Status: "created",
		}
		err := store.Orders.CreateOrder(ctx, order)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Seeded DB")
}
