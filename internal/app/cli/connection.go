package cli

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"
)

type Connection struct {
	Address    string
	Port       int
	Insecure   bool
	CACertPath string
}

func NewConnection(address string, port int, insecure bool, caCertPath string) *Connection {
	return &Connection{address, port, insecure, caCertPath}
}

func GetPath(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, _ := user.Current()
		return strings.Replace(path, "~", usr.HomeDir, 1)
	}
	if strings.HasPrefix(path, "."){
		abs, _ := filepath.Abs("./")
		return strings.Replace(path, ".", abs, 1)
	}
	return path
}

func (c *Connection) GetSecureConnection() (*grpc.ClientConn, derrors.Error) {
	rootCAs := x509.NewCertPool()
	caPath := GetPath(c.CACertPath)
	log.Debug().Str("caCertPath", caPath).Msg("loading CA cert")
	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, derrors.NewInternalError("Error loading CA certificate")
	}
	added := rootCAs.AppendCertsFromPEM(caCert)
	if !added {
		return nil, derrors.NewInternalError("cannot add CA certificate to the pool")
	}

	targetAddress := fmt.Sprintf("%s:%d", c.Address, c.Port)
	log.Debug().Str("address", targetAddress).Msg("creating connection")

	creds := credentials.NewClientTLSFromCert(rootCAs, "")
	log.Debug().Interface("creds", creds.Info()).Msg("Secure credentials")
	sConn, dErr := grpc.Dial(targetAddress, grpc.WithTransportCredentials(creds))
	if dErr != nil {
		return nil, derrors.AsError(dErr, "cannot create connection with the signup service")
	}
	return sConn, nil
}

func (c *Connection) GetInsecureConnection() (*grpc.ClientConn, derrors.Error) {
	log.Warn().Msg("Using insecure connection, TLS options will be ignored")
	targetAddress := fmt.Sprintf("%s:%d", c.Address, c.Port)
	log.Debug().Str("address", targetAddress).Msg("creating connection")
	conn, err := grpc.Dial(targetAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the public api")
	}
	return conn, nil
}

func (c *Connection) GetConnection() (*grpc.ClientConn, derrors.Error) {
	if c.Insecure {
		return c.GetInsecureConnection()
	} else if c.CACertPath != "" {
		return c.GetSecureConnection()
	}
	return nil, derrors.NewInvalidArgumentError("type of connection must be set, either insecure or a CA cert must be provided")
}

func (c *Connection) PrintResultOrError(result interface{}, err error, errMsg string) {
	if err != nil {
		log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg(errMsg)
	} else {
		c.PrintResult(result)
	}
}

func (c *Connection) PrintSuccessOrError(err error, errMsg string, successMsg string) {
	if err != nil {
		log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg(errMsg)
	} else {
		fmt.Println(fmt.Sprintf("{\"msg\":\"%s\"}", successMsg))
	}
}

func (c *Connection) PrintResult(result interface{}) error {
	//Print descriptors
	res, err := json.MarshalIndent(result, "", "  ")
	if err == nil {
		fmt.Println(string(res))
	}
	return err
}
