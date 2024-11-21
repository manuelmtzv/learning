package workers

import (
	"fmt"
	"order-processing/internal/orders"
)

type Manager interface {
	Manage(order <-chan *orders.OrderRequest)
}

type manager struct {
	pending    []*orders.OrderRequest
	processing []*orders.OrderRequest
}

func NewManager(orders int) Manager {
	return &manager{}
}

func (m *manager) Manage(order <-chan *orders.OrderRequest) {
	for {
		receivedOrder := <-order

		if m.orderExists(receivedOrder) {
			fmt.Println("Order already exists")
			continue
		}

		if receivedOrder.InProcess {
			m.processing = append(m.processing, receivedOrder)
		} else {
			m.pending = append(m.pending, receivedOrder)
		}
	}
}

func (m *manager) orderExists(order *orders.OrderRequest) bool {
	for _, o := range m.pending {
		if o.Order.ID == order.Order.ID {
			return true
		}
	}

	for _, o := range m.processing {
		if o.Order.ID == order.Order.ID {
			return true
		}
	}

	return false
}
