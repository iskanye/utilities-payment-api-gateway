package payment

import (
	"context"
	"net"
	"strconv"

	"github.com/iskanye/utilities-payment-proto/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentStatus = payment.PaymentStatus

const (
	PAYMENT_PENDING PaymentStatus = iota
	PAYMENT_OK
	PAYMENT_FAILED
)

type clientApi struct {
	payment payment.PaymentClient
}

type Payment interface {
	ProcessPayment(
		ctx context.Context,
		amount int,
	) (PaymentStatus, error)
}

func New(
	host string,
	port int,
) (clientApi, error) {
	cc, err := grpc.NewClient(
		net.JoinHostPort(host, strconv.Itoa(port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return clientApi{}, err
	}

	return clientApi{payment.NewPaymentClient(cc)}, nil
}

func (c *clientApi) ProcessPayment(
	ctx context.Context,
	amount int,
) (PaymentStatus, error) {
	resp, err := c.payment.ProcessPayment(ctx, &payment.PaymentRequest{
		Amount: int32(amount),
	})
	if err != nil {
		return PAYMENT_FAILED, err
	}

	return resp.GetStatus(), nil
}
