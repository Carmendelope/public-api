/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package cli2

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/nalej/authx/pkg/token"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-login-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/app/options"
	"github.com/nalej/public-api/internal/app/output"
	"github.com/rs/zerolog/log"
)

type Login struct {
	*Connection
	*output.Output
}

// NewLogin creates a new Login structure.
func NewLogin(connection *Connection, output *output.Output) *Login {
	return &Login{connection, output}
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
	credentials := NewCredentials(options.DefaultPath, response.Token, response.RefreshToken)
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
