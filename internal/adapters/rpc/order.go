package rpc

import (
	"context"
	"fmt"
	"service-hf-order-p5/internal/core/domain/entity/dto"
	"service-hf-order-p5/internal/core/domain/rpc"
	op "service-hf-order-p5/order_api_proto"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ rpc.OrderRPC = (*orderRPC)(nil)

type orderRPC struct {
	ctx  context.Context
	host string
	port string
}

func NewOrderRPC(ctx context.Context, host, port string) rpc.OrderRPC {
	return orderRPC{ctx: ctx, host: host, port: port}
}

func (p orderRPC) SaveOrder(order dto.RequestOrder) (*dto.OutputOrder, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", p.host, p.port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	input := op.CreateOrderRequest{
		ClientUuid:  order.ClientUUID,
		VoucherUuid: order.VoucherUUID,
	}

	cc := op.NewOrderClient(conn)

	resp, err := cc.CreateOrder(p.ctx, &input)

	if err != nil {
		return nil, err
	}

	var out = dto.OutputOrder{
		ClientUUID:       resp.ClientUuid,
		VoucherUUID:      resp.VoucherUuid,
		Status:           resp.Status,
		VerificationCode: resp.VerificationCode,
	}

	return &out, nil
}

func (p orderRPC) UpdateOrderByID(id string, order dto.RequestOrder) (*dto.OutputOrder, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", p.host, p.port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	input := op.UpdateOrderRequest{
		Id:               int64(order.ID),
		ClientUuid:       order.ClientUUID,
		VoucherUuid:      order.VoucherUUID,
		Status:           order.Status,
		VerificationCode: order.VerificationCode,
		CreatedAt:        order.CreatedAt,
	}

	cc := op.NewOrderClient(conn)

	resp, err := cc.UpdateOrder(p.ctx, &input)

	if err != nil {
		return nil, err
	}

	var out = dto.OutputOrder{
		ID:               int64(resp.Id),
		ClientUUID:       resp.ClientUuid,
		VoucherUUID:      resp.VoucherUuid,
		Status:           resp.Status,
		VerificationCode: resp.VerificationCode,
		CreatedAt:        resp.CreatedAt,
	}

	return &out, nil
}

func (p orderRPC) GetOrderByID(id string) ([]dto.OutputOrder, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", p.host, p.port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	input := op.GetOrderByIDRequest{
		Id: idInt,
	}

	cc := op.NewOrderClient(conn)

	resp, err := cc.GetOrderByID(p.ctx, &input)

	if err != nil {
		return nil, err
	}

	out := make([]dto.OutputOrder, 0)

	for i := range resp.Items {
		var outItem = dto.OutputOrder{
			ID:               int64(resp.Items[i].Id),
			ClientUUID:       resp.Items[i].ClientUUID,
			VoucherUUID:      resp.Items[i].VoucherUUID,
			Items:            resp.Items[i].Items,
			Status:           resp.Items[i].Status,
			VerificationCode: resp.Items[i].VerificationCode,
			CreatedAt:        resp.Items[i].CreatedAt,
		}
		out = append(out, outItem)
	}

	return out, nil
}

func (p orderRPC) GetOrders() (*dto.OutputOrder, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", p.host, p.port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	input := op.GetOrders{}

	cc := op.NewOrderClient(conn)

	resp, err := cc.GetOrders(p.ctx, &input)

	if err != nil {
		return nil, err
	}

	out := &dto.OutputOrder{
		ID:               int64(resp.Id),
		ClientUUID:       resp.ClientUUID,
		VoucherUUID:      resp.VoucherUUID,
		Items:            resp.Items,
		Status:           resp.Status,
		VerificationCode: resp.VerificationCode,
		CreatedAt:        resp.CreatedAt,
	}

	return out, nil
}
