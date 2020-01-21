/*
 * Copyright 2020 Nalej
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
 */

package cli2

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/public-api/internal/app/options"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

// Connection structure for the public API
type Connection struct {
	// Address to connect to.
	Address string
	// Port where the public API is listening
	Port int
	// Insecure accepts any CA.
	Insecure bool
	// UseTLS specifies whether the target address uses TLS connections
	UseTLS bool
	// CACertPath contains the path of the CA that will be used for verifications.
	CACertPath string
}

// NewConnection creates a new connection object that will establish the communication with the public API.
func NewConnection(address string, port int, insecure bool, useTLS bool, caCertPath string) *Connection {
	return &Connection{address, port, insecure, useTLS, caCertPath}
}

// GetSecureConnection returns a secure connection.
func (c *Connection) GetSecureConnection() (*grpc.ClientConn, derrors.Error) {

	var creds credentials.TransportCredentials

	if c.Insecure {
		cfg := &tls.Config{ServerName: "", InsecureSkipVerify: true}
		creds = credentials.NewTLS(cfg)
		log.Warn().Msg("CA validation will be skipped")
	} else {
		rootCAs := x509.NewCertPool()
		caPath := options.GetPath(c.CACertPath)
		log.Debug().Str("caCertPath", caPath).Msg("loading CA cert")
		caCert, err := ioutil.ReadFile(caPath)
		if err != nil {
			return nil, derrors.NewInternalError("Error loading CA certificate")
		}
		added := rootCAs.AppendCertsFromPEM(caCert)
		if !added {
			return nil, derrors.NewInternalError("cannot add CA certificate to the pool")
		}

		creds = credentials.NewClientTLSFromCert(rootCAs, "")
		log.Debug().Interface("creds", creds.Info()).Msg("Secure credentials")
	}

	targetAddress := fmt.Sprintf("%s:%d", c.Address, c.Port)
	log.Debug().Str("address", targetAddress).Msg("creating connection")

	sConn, dErr := grpc.Dial(targetAddress, grpc.WithTransportCredentials(creds))
	if dErr != nil {
		return nil, derrors.AsError(dErr, "cannot create connection with the signup service")
	}
	return sConn, nil
}

// GetNoTLSConnection creates a connection to a non TLS based endpoint.
func (c *Connection) GetNoTLSConnection() (*grpc.ClientConn, derrors.Error) {
	log.Warn().Msg("Using insecure connection to a non TLS endpoint")
	targetAddress := fmt.Sprintf("%s:%d", c.Address, c.Port)
	log.Debug().Str("address", targetAddress).Msg("creating connection")
	conn, err := grpc.Dial(targetAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the public api")
	}
	return conn, nil
}

// GetConnection creates the appropriate connection type based on the established preferences.
func (c *Connection) GetConnection() (*grpc.ClientConn, derrors.Error) {
	if c.UseTLS {
		if c.Insecure || c.CACertPath != "" {
			return c.GetSecureConnection()
		} else {
			return nil, derrors.NewInvalidArgumentError("expecting CA certificate path or insecure connection")
		}
	}
	return c.GetNoTLSConnection()
}
