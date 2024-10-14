package account

import (
	"context"

	accountv1 "github.com/gatepoint/gatepoint/api/account/v1"
	v1 "github.com/gatepoint/gatepoint/api/gatepoint/v1"
	"github.com/gatepoint/gatepoint/pkg/utils"
	"github.com/gatepoint/gatepoint/pkg/utils/password"
)

var _ v1.AccountServiceServer = new(Server)

type Server struct {
	v1.UnimplementedAccountServiceServer
	accountManager *AccountManager
}

func NewServer() *Server {
	return &Server{
		accountManager: NewAccountManager(),
	}
}

func (s Server) UpdatePassword(ctx context.Context, request *accountv1.UpdatePasswordRequest) (*accountv1.UpdatePasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) CreateAccount(ctx context.Context, request *accountv1.CreateAccountRequest) (*accountv1.CreateAccountResponse, error) {
	hashedPassword, err := password.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}
	err = s.accountManager.CreateAccount(&Account{
		Username:      request.Username,
		PasswordHash:  hashedPassword,
		PasswordMtime: utils.NowUTC(),
	})
	if err != nil {
		return nil, err
	}
	return &accountv1.CreateAccountResponse{}, nil
}
