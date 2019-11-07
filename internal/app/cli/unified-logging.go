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

package cli

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type UnifiedLogging struct {
	Connection
	Credentials
}

func NewUnifiedLogging(address string, port int, insecure bool, useTLS bool, caCertPath string, output string, labelLength int) *UnifiedLogging {
	return &UnifiedLogging{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (u *UnifiedLogging) load() {
	err := u.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (u *UnifiedLogging) getClient() (grpc_public_api_go.UnifiedLoggingClient, *grpc.ClientConn) {
	conn, err := u.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	ulClient := grpc_public_api_go.NewUnifiedLoggingClient(conn)
	return ulClient, conn
}

func parseTime(timeString string) (*timestamp.Timestamp, error) {
	t, err := dateparse.ParseLocal(timeString)
	if err != nil {
		return nil, err
	}
	timeProto, err := ptypes.TimestampProto(t)
	if err != nil {
		return nil, err
	}

	return timeProto, nil
}

func (u *UnifiedLogging) Search(organizationId, instanceId, sgInstanceId, msgFilter, from, to string, desc bool, redirectLog bool) {
	// Validate options
	if organizationId == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	if instanceId == "" {
		log.Fatal().Msg("instanceID cannot be empty")
	}

	// Parse and validate timestamps
	var fromTime, toTime *timestamp.Timestamp
	var err error
	if from != "" {
		fromTime, err = parseTime(from)
		if err != nil {
			log.Fatal().Err(err).Msg("invalid from time")
		}
	}
	if to != "" {
		toTime, err = parseTime(to)
		if err != nil {
			log.Fatal().Err(err).Msg("invalid to time")
		}
	}

	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	var order = grpc_unified_logging_go.SortOrder_ASC
	if desc {
		order = grpc_unified_logging_go.SortOrder_DESC
	}

	searchRequest := &grpc_unified_logging_go.SearchRequest{
		OrganizationId:         organizationId,
		AppInstanceId:          instanceId,
		ServiceGroupInstanceId: sgInstanceId,
		MsgQueryFilter:         msgFilter,
		From:                   fromTime,
		To:                     toTime,
		Order:                  order,
	}

	result, err := client.Search(ctx, searchRequest)
	if redirectLog {
		if err != nil {
			log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("cannot search logs")
		} else {
			log.Info().Str("AppInstanceId", result.AppInstanceId).Str("from", result.From.String()).Str("to", result.To.String()).Msg("app log")
			for _, le := range result.Entries {
				log.Info().Msg(fmt.Sprintf("[%s] %s", le.Timestamp.String(), le.Msg))
			}
		}
	} else {
		u.PrintResultOrError(result, err, "cannot search logs")
	}

}
