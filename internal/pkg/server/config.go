package server

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// Port where the gRPC API service will listen requests.
	Port int
	// HTTPPort where the HTTP gRPC gateway will be listening.
	HTTPPort int
	// SystemModelAddress with the host:port to connect to System Model
	SystemModelAddress string
	// InfrastructureManagerAddress with the host:port to connect to the Infrastructure Manager.
	InfrastructureManagerAddress string
	// ApplicationsManagerAddress with the host:port to connect to the Applications manager.
	ApplicationsManagerAddress string
	// AccessManagerAddress with the host:port to connect to the Access manager.
	AccessManagerAddress string
	// LoggingAddress with the host:port to connect to the Logging component.
	LoggingAddress string
}

func (conf * Config) Validate() derrors.Error {
	return nil
}

func (conf *Config) Print() {
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Int("port", conf.HTTPPort).Msg("HTTP port")
	log.Info().Str("URL", conf.SystemModelAddress).Msg("System Model")
	log.Info().Str("URL", conf.InfrastructureManagerAddress).Msg("Infrastructure Manager")
	log.Info().Str("URL", conf.ApplicationsManagerAddress).Msg("Applications Manager")
	log.Info().Str("URL", conf.AccessManagerAddress).Msg("Access Manager")
	log.Info().Str("URL", conf.LoggingAddress).Msg("Logging Service")
}