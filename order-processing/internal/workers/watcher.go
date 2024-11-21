package workers

import "database/sql"

type OrderGetter interface {
}

type orderGetter struct {
	db *sql.DB
}
