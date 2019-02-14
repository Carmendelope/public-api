/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package cli

import (
	"github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/grpc-public-api-go"
        "github.com/golang/protobuf/ptypes"
        "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/araddon/dateparse"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type UnifiedLogging struct {
	Connection
	Credentials
}

func NewUnifiedLogging(address string, port int, insecure bool, caCertPath string) *UnifiedLogging {
	return &UnifiedLogging{
		Connection:  *NewConnection(address, port, insecure, caCertPath),
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
	t, err := dateparse.ParseAny(timeString)
	if err != nil {
		return nil, err
	}
	timeProto, err := ptypes.TimestampProto(t)
	if err != nil {
		return nil, err
	}

	return timeProto, nil
}

func (u *UnifiedLogging) Search(organizationId, instanceId, sgInstanceId, msgFilter, from, to string) {
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

	searchRequest := &grpc_unified_logging_go.SearchRequest{
		OrganizationId: organizationId,
		AppInstanceId: instanceId,
		ServiceGroupInstanceId: sgInstanceId,
		MsgQueryFilter: msgFilter,
		From: fromTime,
		To: toTime,
	}

	result, err := client.Search(ctx, searchRequest)
	u.PrintResultOrError(result, err, "cannot search logs")
}
