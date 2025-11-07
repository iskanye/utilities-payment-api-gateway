package auth

import (
	"context"
	"net"
	"strconv"

	"github.com/iskanye/utilities-payment-proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type clientApi struct {
	auth auth.AuthClient
}

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (token string, err error)
	Register(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
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

	return clientApi{auth.NewAuthClient(cc)}, nil
}

func (c *clientApi) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	resp, err := c.auth.Login(ctx, &auth.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	return resp.GetToken(), nil
}

func (c *clientApi) Register(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	resp, err := c.auth.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return 0, err
	}

	return resp.GetUserId(), nil
}
