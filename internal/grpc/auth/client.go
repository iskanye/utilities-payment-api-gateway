package auth

import (
	"context"
	"net"
	"strconv"

	"github.com/iskanye/utilities-payment-proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type clientApi struct {
	auth auth.AuthClient
}

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (token string, userId int64, err error)
	Register(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	Validate(
		ctx context.Context,
		token string,
	) (isValid bool, err error)
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

	return clientApi{auth.NewAuthClient(cc)}, nil
}

func (c *clientApi) Login(
	ctx context.Context,
	email string,
	password string,
) (string, int64, error) {
	resp, err := c.auth.Login(ctx, &auth.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", 0, status.Error(codes.Internal, err.Error())
	}

	return resp.GetToken(), resp.GetUserId(), err
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
		return 0, status.Error(codes.Internal, err.Error())
	}

	return resp.GetUserId(), err
}

func (c *clientApi) Validate(
	ctx context.Context,
	token string,
) (bool, error) {
	resp, err := c.auth.Validate(ctx, &auth.ValidateRequest{
		Token: token,
	})
	if err != nil {
		return false, status.Error(codes.Internal, err.Error())
	}

	return resp.GetIsValid(), err
}
