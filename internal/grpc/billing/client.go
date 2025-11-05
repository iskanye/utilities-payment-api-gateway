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
	) (int64, error)
	GetBills(
		ctx context.Context,
		userID int64,
	) ([]models.Bill, error)
	GetBill(
		ctx context.Context,
		billID int64,
	) (models.Bill, error)
	PayBill(
		ctx context.Context,
		billID int64,
	) error
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
			DueDate: bill.GetDueDate(),
		})
	}

	return bills, nil
}

func (c *clientApi) GetBill(
	ctx context.Context,
	billID int64,
) (models.Bill, error) {
	resp, err := c.billing.GetBill(ctx, &billing.BillRequest{
		BillId: billID,
	})
	if err != nil {
		return models.Bill{}, status.Error(codes.Internal, err.Error())
	}

	return models.Bill{
		ID:      resp.GetBillId(),
		Address: resp.GetAddress(),
		Amount:  int(resp.GetAmount()),
		UserID:  resp.GetUserId(),
		DueDate: resp.GetDueDate(),
	}, nil
}

func (c *clientApi) PayBill(
	ctx context.Context,
	billID int64,
) error {
	_, err := c.billing.PayBill(ctx, &billing.PayRequest{
		BillId: billID,
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
