package useCase

import "service-hf-order-p5/internal/core/domain/entity/dto"

type OrderUseCase interface {
	SaveOrder(reqOrder dto.RequestOrder) error
	GetOrderByID(uuid int64) error
	UpdateOrderByID(uuid int64, order dto.RequestOrder) error
	GetOrders() error
}
