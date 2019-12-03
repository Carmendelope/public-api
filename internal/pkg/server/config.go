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
 */

package server

import (
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/derrors"
	"github.com/nalej/public-api/version"
	"github.com/rs/zerolog/log"
	"strings"
)

type Config struct {
	// Debug level is active.
	Debug bool
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
	// UserManagerAddress with the host:port to connect to the Access manager.
	UserManagerAddress string
	// DeviceManagerAddress with the host:port to connect to the Device Manager component.
	DeviceManagerAddress string
	// UnifiedLoggingAddress with the host:port to connect to the Unified Logging Coordinator component.
	UnifiedLoggingAddress string
	// MonitoringManagerAddress with the host:port to connect to the Monitoring Manager component.
	MonitoringManagerAddress string
	// InventoryManagerAddress with the host:port to connect to the Inventory Manager component.
	InventoryManagerAddress string
	// ProvisionerManagerAddress with the host:port to connect to the Provisioner Manager component.
	ProvisionerManagerAddress string
	// AuthSecret contains the shared authx secret.
	AuthSecret string
	// AuthHeader contains the name of the target header.
	AuthHeader string
	// AuthConfigPath contains the path of the file with the authentication configuration.
	AuthConfigPath string
}

func (conf *Config) Validate() derrors.Error {

	if conf.Port <= 0 || conf.HTTPPort <= 0 {
		return derrors.NewInvalidArgumentError("ports must be valid")
	}

	if conf.SystemModelAddress == "" {
		return derrors.NewInvalidArgumentError("systemModelAddress must be set")
	}

	if conf.InfrastructureManagerAddress == "" {
		return derrors.NewInvalidArgumentError("infrastructureManagerAddress must be set")
	}

	if conf.ApplicationsManagerAddress == "" {
		return derrors.NewInvalidArgumentError("applicationsManagerAddress must be set")
	}

	if conf.UserManagerAddress == "" {
		return derrors.NewInvalidArgumentError("userManagerAddress must be set")
	}

	if conf.DeviceManagerAddress == "" {
		return derrors.NewInvalidArgumentError("deviceManagerAddress must be set")
	}

	if conf.UnifiedLoggingAddress == "" {
		return derrors.NewInvalidArgumentError("unifiedLoggingAddress must be set")
	}

	if conf.MonitoringManagerAddress == "" {
		return derrors.NewInvalidArgumentError("monitoringManagerAddress must be set")
	}

	if conf.InventoryManagerAddress == "" {
		return derrors.NewInvalidArgumentError("inventoryManagerAddress must be set")
	}

	if conf.ProvisionerManagerAddress == "" {
		return derrors.NewInvalidArgumentError("provisionerManagerAddress must be set")
	}

	if conf.AuthHeader == "" || conf.AuthSecret == "" {
		return derrors.NewInvalidArgumentError("Authorization header and secret must be set")
	}

	if conf.AuthConfigPath == "" {
		return derrors.NewInvalidArgumentError("authConfigPath must be set")
	}

	return nil
}

// LoadAuthConfig loads the security configuration.
func (conf *Config) LoadAuthConfig() (*interceptor.AuthorizationConfig, derrors.Error) {
	return interceptor.LoadAuthorizationConfig(conf.AuthConfigPath)
}

func (conf *Config) Print() {
	log.Info().Str("app", version.AppVersion).Str("commit", version.Commit).Msg("Version")
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Int("port", conf.HTTPPort).Msg("HTTP port")
	log.Info().Str("URL", conf.SystemModelAddress).Msg("System Model")
	log.Info().Str("URL", conf.InfrastructureManagerAddress).Msg("Infrastructure Manager")
	log.Info().Str("URL", conf.ApplicationsManagerAddress).Msg("Applications Manager")
	log.Info().Str("URL", conf.UserManagerAddress).Msg("User Manager")
	log.Info().Str("URL", conf.UnifiedLoggingAddress).Msg("Unified Logging Coordinator Service")
	log.Info().Str("URL", conf.MonitoringManagerAddress).Msg("Monitoring Manager Service")
	log.Info().Str("URL", conf.DeviceManagerAddress).Msg("Device Manager Service")
	log.Info().Str("URL", conf.InventoryManagerAddress).Msg("Inventory Manager Service")
	log.Info().Str("URL", conf.ProvisionerManagerAddress).Msg("Provisioner Manager service")
	log.Info().Str("header", conf.AuthHeader).Str("secret", strings.Repeat("*", len(conf.AuthSecret))).Msg("Authorization")
	log.Info().Str("path", conf.AuthConfigPath).Msg("Permissions file")

}
