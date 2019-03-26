/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cli

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/nalej/authx/pkg/token"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-login-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
)

type Login struct {
	Connection
}

// NewLogin creates a new Login structure.
func NewLogin(address string, port int, insecure bool, useTLS bool, caCertPath string, output string) *Login {
	return &Login{
		*NewConnection(address, port, insecure, useTLS, caCertPath, output),
	}
}

// Login into the platform using email and password.
func (l *Login) Login(email string, password string) (*Credentials, derrors.Error) {
	c, err := l.GetConnection()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	loginClient := grpc_login_api_go.NewLoginClient(c)
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	loginRequest := &grpc_authx_go.LoginWithBasicCredentialsRequest{
		Username: email,
		Password: password,
	}
	response, lErr := loginClient.LoginWithBasicCredentials(ctx, loginRequest)
	if lErr != nil {
		return nil, conversions.ToDerror(lErr)
	}
	log.Debug().Str("token", response.Token).Msg("Login success")
	credentials := NewCredentials(DefaultPath, response.Token, response.RefreshToken)
	sErr := credentials.Store()
	if sErr != nil {
		return nil, sErr
	}
	return credentials, nil
}

func (l *Login) GetPersonalClaims(credentials *Credentials) (*token.Claim, derrors.Error) {
	parser := jwt.Parser{
		SkipClaimsValidation: true,
	}
	tk, _, err := parser.ParseUnverified(credentials.Token, &token.Claim{})
	if err != nil {
		return nil, derrors.AsError(err, "cannot parse token")
	}
	return tk.Claims.(*token.Claim), nil
}
