package billing

import (
	"context"
	"errors"
	"io"
	"net"
	"strconv"

	"github.com/iskanye/utilities-payment-proto/billing"
	"github.com/iskanye/utilities-payment/pkg/models"
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
		userID int64,
	) (billId int64, err error)
	GetBills(
		ctx context.Context,
		userID int64,
	) (bills []models.Bill, err error)
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
	userID int64,
) (int64, error) {
	resp, err := c.billing.AddBill(ctx, &billing.Bill{
		Address: address,
		Amount:  int32(amount),
		UserId:  userID,
	})
	if err != nil {
		return 0, status.Error(codes.Internal, err.Error())
	}

	return resp.GetBillId(), nil
}

func (c *clientApi) GetBills(
	ctx context.Context,
	userID int64,
) ([]models.Bill, error) {
	resp, err := c.billing.GetBills(ctx, &billing.BillsRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	bills := make([]models.Bill, 0)

	for {
		bill, err := resp.Recv()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, status.Error(codes.Internal, err.Error())
			}
			break
		}

		bills = append(bills, models.Bill{
			ID:      bill.GetBillId(),
			Address: bill.GetAddress(),
			Amount:  int(bill.GetAmount()),
			UserID:  userID,
			DueDate: *bill.DueDate,
		})
	}

	return bills, nil
}
