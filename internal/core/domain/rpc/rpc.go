package rpc

import "service-hf-order-p5/internal/core/domain/entity/dto"

type OrderRPC interface {
	GetOrderByID(uuid string) (*dto.OutputOrder, error)
	SaveOrder(order dto.RequestOrder) (*dto.OutputOrder, error)
	UpdateOrderByID(id string, order dto.RequestOrder) (*dto.OutputOrder, error)
	GetOrders(id string) ([]dto.OutputOrder, error)
}

type OrderWorkerRPC interface {
	GetOrderByID(uuid string) (*dto.OutputOrder, error)
	SaveOrder(order dto.RequestOrder) (*dto.OutputOrder, error)
	UpdateOrderByID(id string, order dto.RequestOrder) (*dto.OutputOrder, error)
	GetOrders(id string) ([]dto.OutputOrder, error)
}
