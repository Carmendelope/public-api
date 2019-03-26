package cli

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"
)

// Connection structure for the public API
type Connection struct {
	// Address to connect to.
	Address    string
	// Port where the public API is listening
	Port       int
	// Insecure accepts any CA.
	Insecure   bool
	// UseTLS specifies whether the target address uses TLS connections
	UseTLS bool
	// CACertPath contains the path of the CA that will be used for verifications.
	CACertPath string
	// Output specifies how the output will be shown.
	output string
}

// NewConnection creates a new connection object that will establish the communication with the public API.
func NewConnection(address string, port int, insecure bool, useTLS bool, caCertPath string, output string) *Connection {
	return &Connection{address, port, insecure, useTLS,caCertPath, output}
}

// GetPath resolves a given path by adding support for relative paths.
func GetPath(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, _ := user.Current()
		return strings.Replace(path, "~", usr.HomeDir, 1)
	}
	if strings.HasPrefix(path, "../"){
		abs, _ := filepath.Abs("../")
		return strings.Replace(path, "..", abs, 1)
	}
	if strings.HasPrefix(path, "."){
		abs, _ := filepath.Abs("./")
		return strings.Replace(path, ".", abs, 1)
	}
	return path
}

// GetSecureConnection returns a secure connection.
func (c *Connection) GetSecureConnection() (*grpc.ClientConn, derrors.Error) {

	var creds credentials.TransportCredentials

	if c.Insecure {
		cfg := &tls.Config{ServerName: "", InsecureSkipVerify: true}
		creds = credentials.NewTLS(cfg)
		log.Warn().Msg("CA validation will be skipped")
	}else{
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
	if c.UseTLS{
		if c.Insecure || c.CACertPath != "" {
			return c.GetSecureConnection()
		}else{
			return nil, derrors.NewInvalidArgumentError("expecting CA certificate path or insecure connection")
		}
	}
	return c.GetNoTLSConnection()
}

func (c *Connection) PrintResultOrError(result interface{}, err error, errMsg string) {
	if err != nil {
		converted := conversions.ToDerror(err)
		if zerolog.GlobalLevel() == zerolog.DebugLevel{
			log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg(errMsg)
		}else{
			log.Fatal().Str("err", converted.Error()).Msg(errMsg)
		}
	} else {
		if c.asText(){
			c.PrintResultAsTable(result)
		}else{
			c.PrintResult(result)
		}
	}
}

func (c *Connection) ExitOnError(err error, errMsg string) {
	if err != nil {
		converted := conversions.ToDerror(err)
		if zerolog.GlobalLevel() == zerolog.DebugLevel{
			log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg(errMsg)
		}else{
			log.Fatal().Str("err", converted.Error()).Msg(errMsg)
		}
	}
}

// TODO Refactor a move print methods to other entity.
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

func (c * Connection) asText() bool {
	return strings.ToLower(c.output) == "text"
}

func (c * Connection) PrintResultAsTable(result interface{}) {
	table := AsTable(result)
	table.Print()
}