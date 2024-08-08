package useCase

import (
	"errors"
	"service-hf-order-p5/internal/core/domain/entity/dto"
	"service-hf-order-p5/internal/core/domain/useCase"
)

var _ useCase.OrderUseCase = (*orderUseCase)(nil)

type orderUseCase struct {
}

func NewOrderUseCase() orderUseCase {
	return orderUseCase{}
}

func (p orderUseCase) SaveOrder(reqOrder dto.RequestOrder) error {
	order := reqOrder.Order()

	if err := order.Status.Validate(); err != nil {
		return err
	}

	if len(order.VerificationCode.Value) == 0 {
		order.VerificationCode.Generate()
	}

	if err := order.VerificationCode.Validate(); err != nil {
		order.VerificationCode.Generate()
	}

	reqOrder.VerificationCode = order.VerificationCode.Value

	return nil

}

func (o orderUseCase) UpdateOrderByID(id int64, reqOrder dto.RequestOrder) error {
	if id < 1 {
		return errors.New("the id is not valid for consult")
	}

	order := reqOrder.Order()

	if err := order.Status.Validate(); err != nil {
		return err
	}

	if err := order.VerificationCode.Validate(); err != nil {
		order.VerificationCode.Generate()
	}

	reqOrder.VerificationCode = order.VerificationCode.Value

	return nil
}

func (o orderUseCase) GetOrders() error {
	return nil
}

func (o orderUseCase) GetOrderByID(id int64) error {
	if id < 1 {
		return errors.New("the id is not valid for consult")
	}

	return nil
}
