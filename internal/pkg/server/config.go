package server

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// Address where the API service will listen requests.
	Port int
}

func (conf * Config) Validate() derrors.Error {
	return nil
}

func (conf *Config) Print() {
	log.Info().Int("port", conf.Port).Msg("gRPC port")
}