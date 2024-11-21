package workers

import (
	"log"
	"order-processing/internal/orders"
)

type Auditor interface {
	Audit(order chan<- orders.OrderRequest)
}

type auditor struct {
}

func NewAudit(orders int) Auditor {
	return &auditor{}
}

func (a *auditor) Audit(order chan<- orders.OrderRequest) {
	for {
		order <- orders.OrderRequest{}
		log.Printf("Order logged\n")
	}
}
