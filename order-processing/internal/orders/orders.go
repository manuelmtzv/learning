package orders

import "order-processing/internal/models"

type OrderRequest struct {
	Order     models.Order
	InProcess bool
	Success   chan bool
}
