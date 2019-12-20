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

package cli

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-log-download-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"os"
	"time"
)

const OrderByField = "timestamp"
const FollowSleep = time.Second * 5

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
	msgFilter, from, to string, desc bool, redirectLog bool, follow bool, nFirst bool) {
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
		fromInt = fromTime.UnixNano()
	}
	if to != "" {
		toTime, err = dateparse.ParseLocal(to)
		if err != nil {
			log.Fatal().Err(err).Msg("invalid to time")
		}
		toInt = toTime.UnixNano()
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
		NFirst:                 nFirst,
	}

	toReturned := u.callSearch(searchRequest, redirectLog, client)

	if follow {
		ticker := time.NewTicker(FollowSleep)
		done := make(chan bool)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if toReturned != 0 {
					// next milisecond
					searchRequest.From = toReturned + 1000000
				}
				toReturned = u.callSearch(searchRequest, redirectLog, client)
			}
		}
	}

}

func (u *UnifiedLogging) Download(organizationId, descriptorId, instanceId, sgId, sgInstanceId, serviceId, serviceInstanceId,
	msgFilter, from, to string, desc bool, includeMetadata bool, outputPath string) {
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
		fromInt = fromTime.UnixNano()
	}
	if to != "" {
		toTime, err = dateparse.ParseLocal(to)
		if err != nil {
			log.Fatal().Err(err).Msg("invalid to time")
		}
		toInt = toTime.UnixNano()
	}

	u.load()

	client, conn := u.getClient()
	defer conn.Close()

	var order = grpc_common_go.OrderOptions{Order: grpc_common_go.Order_ASC, Field: OrderByField}
	if desc {
		order.Order = grpc_common_go.Order_DESC
	}

	request := &grpc_log_download_manager_go.DownloadLogRequest{
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
		IncludeMetadata:        includeMetadata,
	}
	ctx, cancel := u.GetContext()
	defer cancel()

	log.Info().Msg("download log entries...")
	response, err := client.DownloadLog(ctx, request)

	// check the operation status
	if err != nil {
		u.PrintResultOrError(response, err, "cannot download log entries")
	} else {
		// no error, check!
		rCtx, rCancel := context.WithTimeout(context.Background(), time.Minute*10)
		defer rCancel()
		ticker := time.NewTicker(FollowSleep)
		done := make(chan bool)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				checkResponse, err := u.callCheck(organizationId, response.RequestId)
				if err != nil {
					u.PrintResultOrError(checkResponse, err, "cannot check file status")
					return
				}
				log.Info().Str("state", checkResponse.StateName).Msg("check download status...")
				if checkResponse.State == grpc_log_download_manager_go.DownloadLogState_ERROR {
					u.PrintResultOrError(checkResponse, err, "")
					return
				}
				if checkResponse.State == grpc_log_download_manager_go.DownloadLogState_READY {
					u.callGet(checkResponse, outputPath)
					return
				}
				// check
			case <-rCtx.Done():
				log.Info().Msg("error")
				err = derrors.NewDeadlineExceededError("Download error. Try it manually")
				return
			}
		}
	}
}

func (u *UnifiedLogging) Check(organizationId, requestId string) {
	// Validate options
	if organizationId == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if requestId == "" {
		log.Fatal().Msg("requestId cannot be empty")
	}
	response, err := u.callCheck(organizationId, requestId)

	u.PrintResultOrError(response, err, "cannot check the status of the request")
}

func (u *UnifiedLogging) Get (organizationId, requestId, outputPath string ){
	// Validate options
	if organizationId == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if requestId == "" {
		log.Fatal().Msg("requestId cannot be empty")
	}

	// get the url
	response, err := u.callCheck(organizationId, requestId)
	if err != nil {
		u.PrintResultOrError(response, err, "cannot check the status of the request")
	}else{
		u.callGet(response, outputPath)
	}
	return
}

func (u *UnifiedLogging) List (organizationId string, watch bool) {
	// Validate options
	if organizationId == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	u.load()

	client, conn := u.getClient()
	defer conn.Close()
	ctx, cancel := u.GetContext()
	defer cancel()

	request := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationId,
	}
	response, err := client.List(ctx, request)

	toCompare := make (map[string]*grpc_public_api_go.DownloadLogResponse, 0)
	for _, resp := range response.Responses {
		toCompare[resp.RequestId] = resp
	}

	u.PrintResultOrError(response, err, "cannot list the status of the requests")

	if watch {
		ticker := time.NewTicker(FollowSleep)
		for {
			select {
			case <-ticker.C:
				toShow := make([]*grpc_public_api_go.DownloadLogResponse, 0)

				watchCtx, watchCancel := u.GetContext()
				operations, err := client.List(watchCtx, request)
				if err != nil {
					u.PrintResultOrError(operations, err, "cannot list the status of the requests")
				}else {

					for _, retrieved := range operations.Responses {
						found, exists := toCompare[retrieved.RequestId]
						if !exists {
							toShow = append(toShow, retrieved)
						} else if found.State != retrieved.State || found.Info != retrieved.Info {
							toShow = append(toShow, retrieved)
						}
						toCompare[retrieved.RequestId] = retrieved
					}
				}

				if len(toShow) > 0 {
					operations.Responses = toShow
					fmt.Println("")
					u.PrintResultOrError(operations, err, "cannot list the status of the requests")
				}

				watchCancel()
			}
		}
	}
}

func (u *UnifiedLogging) callGet(checkResponse *grpc_public_api_go.DownloadLogResponse, outputPath string) {

	tlsConfigInsecure := &tls.Config{InsecureSkipVerify: true}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    tlsConfigInsecure,
	}
	client := &http.Client{
		Transport: tr,
	}
	req, err := http.NewRequest("GET", checkResponse.Url, nil)
	req.Header.Add("authorization", u.Token)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		u.PrintResultOrError(req, err, "error getting the file")
	}

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		u.PrintResultOrError(req, derrors.NewInternalError(resp.Status), buf.String())
		return
	}
	// Create the file
	outputFilePath := fmt.Sprintf("%s%s.zip",outputPath, checkResponse.RequestId)
	out, err := os.Create(outputFilePath)
	if err != nil {
		u.PrintResultOrError(checkResponse, err, "error creating the file")
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		u.PrintResultOrError(checkResponse, err, "error copying the file")
		return
	}
	fmt.Printf("\nLog Entries file: %s\n", outputFilePath)
}

func (u *UnifiedLogging) callCheck(organizationId, requestId string) (*grpc_public_api_go.DownloadLogResponse, error) {
	u.load()

	client, conn := u.getClient()
	defer conn.Close()
	ctx, cancel := u.GetContext()
	defer cancel()

	return client.Check(ctx, &grpc_log_download_manager_go.DownloadRequestId{
		OrganizationId: organizationId,
		RequestId:      requestId,
	})

}


// returns searchTo field (we need this value to update the next search)
func (u *UnifiedLogging) callSearch(searchRequest *grpc_public_api_go.SearchRequest, redirectLog bool, client grpc_public_api_go.UnifiedLoggingClient) int64 {

	followCtx, followCancel := u.GetContext()
	defer followCancel()
	result, err := client.Search(followCtx, searchRequest)
	if redirectLog {
		if err != nil {
			log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("cannot search logs")
		} else {
			log.Info().Str("OrganizationId", result.OrganizationId).Str("from", string(result.From)).
				Str("to", string(result.To)).Msg("app log")
			for _, le := range result.Entries {
				log.Info().Msg(fmt.Sprintf("[%s] %s", string(le.Timestamp), le.Msg))
			}
		}
	} else {
		u.PrintResultOrError(result, err, "cannot search logs")
	}
	return result.To
}
