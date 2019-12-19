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

package unified_logging

import (
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-log-download-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/entities"
	"github.com/nalej/public-api/internal/pkg/server/common"
	"github.com/nalej/public-api/internal/pkg/server/decorators"
	"github.com/rs/zerolog/log"
)

type Manager struct {
	unifiedLoggingClient grpc_application_manager_go.UnifiedLoggingClient
	logDownloadClient    grpc_log_download_manager_go.LogDownloadManagerClient
}

func NewManager(unifiedLoggingClient grpc_application_manager_go.UnifiedLoggingClient,
	logDownloadClient grpc_log_download_manager_go.LogDownloadManagerClient) Manager {
	return Manager{unifiedLoggingClient, logDownloadClient}
}

func (m *Manager) Search(request *grpc_public_api_go.SearchRequest) (*grpc_application_manager_go.LogResponse, error) {
	log.Debug().Interface("request", request).Msg("Search request")
	ctx, cancel := common.GetContext()
	defer cancel()
	convertedLog, err := m.unifiedLoggingClient.Search(ctx, entities.NewSearchRequest(request))

	if err != nil {
		return nil, err
	}

	// if sorting requested -> apply the decorator
	if request.Order != nil {
		sortOptions := decorators.OrderOptions{Field: request.Order.Field, Asc: request.Order.Order == grpc_public_api_go.Order_ASC}
		sortingResponse := decorators.ApplyDecorator(convertedLog.Entries, decorators.NewOrderDecorator(sortOptions))
		if sortingResponse.Error != nil {
			return nil, conversions.ToGRPCError(sortingResponse.Error)
		} else {
			convertedLog.Entries = sortingResponse.LogResponseList
		}
	}
	return convertedLog, nil
}

// Check checks the state of the download operation
func (m *Manager) Check(requestId *grpc_log_download_manager_go.DownloadRequestId) (*grpc_public_api_go.DownloadLogResponse, error) {
	log.Debug().Interface("request", requestId).Msg("Check request")
	ctx, cancel := common.GetContext()
	defer cancel()
	response, err := m.logDownloadClient.Check(ctx, requestId)
	if err != nil {
		return nil, err
	}
	return entities.ToPublicAPIDownloadLogReponse(response), nil

}

// DownloadLog ask for log entries and store them into a zip file
func (m *Manager) DownloadLog(request *grpc_log_download_manager_go.DownloadLogRequest) (*grpc_public_api_go.DownloadLogResponse, error) {
	log.Debug().Interface("request", request).Msg("DownloadLog request")
	ctx, cancel := common.GetContext()
	defer cancel()
	response, err := m.logDownloadClient.DownloadLog(ctx, request)
	if err != nil {
		return nil, err
	}
	return entities.ToPublicAPIDownloadLogReponse(response), nil
}

func (m *Manager) List(request *grpc_organization_go.OrganizationId) (*grpc_public_api_go.DownloadLogResponseList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	responses, err := m.logDownloadClient.List(ctx, request)
	if err != nil {
		return nil, err
	}

	responseList := make([]*grpc_public_api_go.DownloadLogResponse, len(responses.Responses))
	for i, resp := range responses.Responses {
		responseList[i] = entities.ToPublicAPIDownloadLogReponse(resp)
	}
	return &grpc_public_api_go.DownloadLogResponseList{
		Responses: responseList,
	}, nil
}
