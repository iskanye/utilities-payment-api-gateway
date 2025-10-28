package billing

import (
	"context"
	"net"
	"strconv"

	"github.com/iskanye/utilities-payment-proto/billing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type clientApi struct {
	billing billing.BillingClient
}

type Billing interface {
	AddBill(
		ctx context.Context,
		address string,
		amount int,
	) (billId int64, err error)
}

func New(
	host string,
	port int,
) (clientApi, error) {
	cc, err := grpc.NewClient(
		net.JoinHostPort(host, strconv.Itoa(port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return clientApi{}, status.Error(codes.Unavailable, err.Error())
	}

	return clientApi{billing.NewBillingClient(cc)}, nil
}

func (c *clientApi) AddBill(
	ctx context.Context,
	address string,
	amount int,
) (int64, error) {
	resp, err := c.billing.AddBill(ctx, &billing.Bill{
		Address: address,
		Amount:  int32(amount),
	})
	if err != nil {
		return 0, status.Error(codes.Internal, err.Error())
	}

	return resp.GetBillId(), nil
}
