package cli

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Connection struct{
	Address string
	Port int
}

func NewConnection(address string, port int) * Connection{
	return &Connection{address, port}
}

func (c* Connection) GetConnection() (*grpc.ClientConn, derrors.Error){
	targetAddress := fmt.Sprintf("%s:%d", c.Address, c.Port)
	log.Debug().Str("address", targetAddress).Msg("creating connection")
	conn, err := grpc.Dial(targetAddress, grpc.WithInsecure())
	if err != nil{
		return nil, derrors.AsError(err, "cannot create connection with the public api")
	}
	return conn, nil
}

func (c * Connection) PrintResult(result interface{}) error {
	//Print descriptors
	res, err := json.MarshalIndent(result, "", "  ")
	if err == nil {
		fmt.Println(string(res))
	}
	return err
}