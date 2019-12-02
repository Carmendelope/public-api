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
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"time"
)

const OrderByField = "timestamp"
const FollowSleep = time.Second * 3


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

//l.Search(cliOptions.Resolve("organizationID", organizationID), descriptorID, instanceID, sgID, sgInstanceID, serviceID, serviceInstanceID, message, from, to, desc, redirectLog)
func (u *UnifiedLogging) Search(organizationId, descriptorId, instanceId, sgId, sgInstanceId, serviceId, serviceInstanceId,
	msgFilter, from, to string, desc bool, redirectLog bool, follow bool) {
	// Validate options
	if organizationId == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	// Parse and validate timestamps
	var fromTime, toTime time.Time
	var fromInt, toInt int64
	fromInt = 0
	toInt = 0
	var err error
	if from != "" {
		fromTime, err = dateparse.ParseLocal(from)
		if err != nil {
			log.Fatal().Err(err).Msg("invalid from time")
		}
		fromInt = fromTime.Unix()
	}
	if to != "" {
		toTime, err = dateparse.ParseLocal(to)
		if err != nil {
			log.Fatal().Err(err).Msg("invalid to time")
		}
		toInt = toTime.Unix()
	}

	if follow && (toInt != 0 || fromInt != 0) {
		log.Fatal().Msg("time range can not be informed with follow option")
	}

	u.load()

	client, conn := u.getClient()
	defer conn.Close()

	var order = grpc_public_api_go.OrderOptions{Order: grpc_public_api_go.Order_ASC, Field: OrderByField}
	if desc {
		order.Order = grpc_public_api_go.Order_DESC
	}

	searchRequest := &grpc_public_api_go.SearchRequest{
		OrganizationId:         organizationId,
		AppDescriptorId:        descriptorId,
		AppInstanceId:          instanceId,
		ServiceGroupId:         sgId,
		ServiceGroupInstanceId: sgInstanceId,
		ServiceId:              serviceId,
		ServiceInstanceId:      serviceInstanceId,
		MsgQueryFilter:         msgFilter,
		From:                   fromInt,
		To:                     toInt, 
		Order:                  &order,
	}

	toReturned := u.callSearch(searchRequest, redirectLog, client)

	ticker := time.NewTicker(FollowSleep)
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if toReturned != 0 {
				searchRequest.From = toReturned + time.Unix(1, 0, ).Unix()
			}
			toReturned = u.callSearch(searchRequest, redirectLog, client)
		}
	}

}

// returns searchTo field (we need this value to update the next search)
func (u *UnifiedLogging) callSearch(searchRequest *grpc_public_api_go.SearchRequest, redirectLog bool, client grpc_public_api_go.UnifiedLoggingClient) int64{
	followCtx, followCancel := u.GetContext()
	defer followCancel()
	result, err := client.Search(followCtx, searchRequest)
	if redirectLog {
		if err != nil {
			log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("cannot search logs")
		} else {
			log.Info().Str("OrganizationId", result.OrganizationId).Str("from", time.Unix(result.From, 0).String()).
				Str("to", time.Unix(result.To, 0).String()).Msg("app log")
			for _, le := range result.Entries {
				log.Info().Msg(fmt.Sprintf("[%s] %s", time.Unix(le.Timestamp, 0).String(), le.Msg))
			}
		}
	} else {
		u.PrintResultOrError(result, err, "cannot search logs")
	}
	return result.To
}