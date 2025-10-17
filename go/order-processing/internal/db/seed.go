package db

import (
	"context"
	"database/sql"
	"fmt"
	"order-processing/internal/models"
	"order-processing/internal/store"
	"sync"
)

func Seed(store *store.Storage, db *sql.DB) {
	ctx := context.Background()

	var wg sync.WaitGroup

	for i := 0; i < 100000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			order := &models.Order{
				Status: "created",
			}
			err := store.Orders.CreateOrder(ctx, order)
			if err != nil {
				fmt.Println(err)
			}
		}()

	}

	wg.Wait()

	fmt.Println("Seeded DB")
}
