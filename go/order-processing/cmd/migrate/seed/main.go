package main

import (
	"log"
	"order-processing/internal/db"
	"order-processing/internal/env"
	"order-processing/internal/store"
	"time"
)

func main() {
	err := env.Load()
	if err != nil {
		log.Panic(err)
	}

	add := env.GetString("DB_ADDR", "")
	conn, err := db.New(add, 25, 25, 15*time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store, conn)
}
