package application

import (
	"context"
	"errors"
	l "service-hf-order-p5/external/logger"
	ps "service-hf-order-p5/external/strings"
	"service-hf-order-p5/internal/core/domain/entity/dto"
	"service-hf-order-p5/internal/core/domain/rpc"
)

type Application interface {
	GetOrderByID(msgID string, uuid string) (*dto.OutputOrder, error)
	SaveOrder(msgID string, order dto.RequestOrder) (*dto.OutputOrder, error)
	UpdateOrderByID(msgID string, id string, order dto.RequestOrder) (*dto.OutputOrder, error)
	GetOrders(msgID string, order string) ([]dto.OutputOrder, error)
}

type application struct {
	ctx            context.Context
	orderRPC       rpc.OrderRPC
	orderWorkerRPC rpc.OrderWorkerRPC
}

func (app *application) setMessageIDCtx(msgID string) {

}

func NewApplication(ctx context.Context, orderRPC rpc.OrderRPC, orderWorkerRPC rpc.OrderWorkerRPC) Application {
	return application{
		ctx:            ctx,
		orderRPC:       orderRPC,
		orderWorkerRPC: orderWorkerRPC,
	}
}

func (app application) GetOrderByID(msgID string, uuid string) (*dto.OutputOrder, error) {
	app.setMessageIDCtx(msgID)

	l.Infof(msgID, "GetOrderByIDApp: ", " | ", uuid)

	go app.orderRPC.GetOrderByID(uuid) // pub

	orderRpc, err := app.orderWorkerRPC.GetOrderByID(uuid)

	if err != nil {
		l.Errorf(msgID, "GetOrderByIDApp error: ", " | ", err)
		return nil, err
	}

	if orderRpc == nil {
		l.Infof(msgID, "GetOrderByIDApp output: ", " | ", nil)
		return nil, nil
	}

	order := &dto.OutputOrder{
		ID:             int64(orderRpc.Id),
		ClientUUID:       orderRpc.ClientUUID,
		VoucherUUID:      orderRpc.VoucherUUID,
		Items:            orderRpc.Items,
		Status:           orderRpc.Status,
		VerificationCode: orderRpc.VerificationCode,
		CreatedAt:        orderRpc.CreatedAt,
	}

	l.Infof(msgID, "GetOrderByIDApp output: ", " | ", order)
	return order, nil
}

func (app application) SaveOrder(msgID string, order dto.RequestOrder) (*dto.OutputOrder, error) {
	app.setMessageIDCtx(msgID)

	l.Infof(msgID, "SaveOrderApp: ", " | ", ps.MarshalString(order))

	go app.orderRPC.SaveOrder(order) // pub

	pRepo, err := app.orderWorkerRPC.SaveOrder(order)

	if err != nil {
		l.Errorf(msgID, "SaveOrderApp error: ", " | ", err)
		return nil, err
	}

	if pRepo == nil {
		l.Infof(msgID, "SaveOrderApp output: ", " | ", nil)
		return nil, errors.New("is not possible to save order because it's null")
	}

	out := &dto.OutputOrder{
		ID:             int64(pRepo.Id),
		ClientUUID:       pRepo.ClientUUID,
		VoucherUUID:      pRepo.VoucherUUID,
		Items:            pRepo.Items,
		Status:           pRepo.Status,
		VerificationCode: pRepo.VerificationCode,
		CreatedAt:        pRepo.CreatedAt,
	}

	l.Infof(msgID, "SaveOrderApp output: ", " | ", ps.MarshalString(out))
	return out, nil
}

func (app application) UpdateOrderByID(msgID string, id string, order dto.RequestOrder) (*dto.OutputOrder, error) {
	app.setMessageIDCtx(msgID)

	l.Infof(msgID, "UpdateOrderByIDApp: ", " | ", id, " | ", ps.MarshalString(order))

	go app.orderRPC.UpdateOrderByID(id, order)

	p, err := app.orderWorkerRPC.UpdateOrderByID(id, order)

	if err != nil {
		l.Errorf(msgID, "UpdateOrderByIDApp error: ", " | ", err)
		return nil, err
	}

	if p == nil {
		orderNullErr := errors.New("is not possible to save order because it's null")
		l.Errorf(msgID, "UpdateOrderByIDApp output: ", " | ", orderNullErr)
		return nil, orderNullErr
	}

	out := &dto.OutputOrder{
		ID:             int64(p.Id),
		ClientUUID:       p.ClientUUID,
		VoucherUUID:      p.VoucherUUID,
		Items:            p.Items,
		Status:           p.Status,
		VerificationCode: p.VerificationCode,
		CreatedAt:        p.CreatedAt,
	}

	l.Infof(msgID, "UpdateOrderByIDApp output: ", " | ", ps.MarshalString(out))
	return out, nil
}

func (app application) GetOrders(msgID string, uuid string) (*dto.OutputOrder, error) {
	app.setMessageIDCtx(msgID)

	l.Infof(msgID, "GetOrdersApp: ", " | ", uuid)

	go app.orderRPC.GetOrders(uuid) // pub

	orderRpc, err := app.orderWorkerRPC.GetOrders(uuid)

	if err != nil {
		l.Errorf(msgID, "GetOrdersApp error: ", " | ", err)
		return nil, err
	}

	if orderRpc == nil {
		l.Infof(msgID, "GetOrdersApp output: ", " | ", nil)
		return nil, nil
	}

	order := &dto.OutputOrder{
		ID:             int64(orderRpc.Id),
		ClientUUID:       orderRpc.ClientUUID,
		VoucherUUID:      orderRpc.VoucherUUID,
		Items:            orderRpc.Items,
		Status:           orderRpc.Status,
		VerificationCode: orderRpc.VerificationCode,
		CreatedAt:        orderRpc.CreatedAt,
	}

	l.Infof(msgID, "GetOrdersApp output: ", " | ", order)
	return order, nil
}
