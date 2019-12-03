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

package applications

import (
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/public-api/internal/pkg/entities"
	"github.com/nalej/public-api/internal/pkg/server/common"
)

type Manager struct {
	appClient grpc_application_manager_go.ApplicationManagerClient
}

func NewManager(appClient grpc_application_manager_go.ApplicationManagerClient) Manager {
	return Manager{appClient}
}

// AddAppDescriptor adds a new application descriptor to a given organization.
func (m *Manager) AddAppDescriptor(addRequest *grpc_application_go.AddAppDescriptorRequest) (*grpc_application_go.AppDescriptor, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.appClient.AddAppDescriptor(ctx, addRequest)
}

// ListAppDescriptors retrieves a list of application descriptors.
func (m *Manager) ListAppDescriptors(organizationID *grpc_organization_go.OrganizationId) (*grpc_application_go.AppDescriptorList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.appClient.ListAppDescriptors(ctx, organizationID)
}

// GetAppDescriptor retrieves a given application descriptor.
func (m *Manager) GetAppDescriptor(appDescriptorID *grpc_application_go.AppDescriptorId) (*grpc_application_go.AppDescriptor, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.appClient.GetAppDescriptor(ctx, appDescriptorID)
}

// UpdateAppDescriptor allows the user to update the information of a registered descriptor.
func (m *Manager) UpdateAppDescriptor(request *grpc_application_go.UpdateAppDescriptorRequest) (*grpc_application_go.AppDescriptor, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.appClient.UpdateAppDescriptor(ctx, request)
}

// DeleteAppDescriptor deletes a given application descriptor.
func (m *Manager) DeleteAppDescriptor(appDescriptorID *grpc_application_go.AppDescriptorId) (*grpc_common_go.Success, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.appClient.RemoveAppDescriptor(ctx, appDescriptorID)
}

// Deploy an application descriptor.
func (m *Manager) Deploy(deployRequest *grpc_application_manager_go.DeployRequest) (*grpc_application_manager_go.DeploymentResponse, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.appClient.Deploy(ctx, deployRequest)
}

// Undeploy a running application instance.
func (m *Manager) Undeploy(undeployRequest *grpc_application_manager_go.UndeployRequest) (*grpc_common_go.Success, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.appClient.Undeploy(ctx, undeployRequest)
}

// ListAppInstances retrieves a list of application descriptors.
func (m *Manager) ListAppInstances(organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.AppInstanceList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	apps, err := m.appClient.ListAppInstances(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	result := make([]*grpc_public_api_go.AppInstance, 0)
	for _, app := range apps.Instances {
		result = append(result, entities.ToPublicAPIAppInstance(app))
	}
	return &grpc_public_api_go.AppInstanceList{
		Instances: result,
	}, nil
}

// GetAppDescriptor retrieves a given application descriptor.
func (m *Manager) GetAppInstance(appInstanceID *grpc_application_go.AppInstanceId) (*grpc_public_api_go.AppInstance, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	inst, err := m.appClient.GetAppInstance(ctx, appInstanceID)
	if err != nil {
		return nil, err
	}
	return entities.ToPublicAPIAppInstance(inst), nil
}

// ListInstanceParameters retrieves a list of instance parameters
func (m *Manager) ListInstanceParameters(appInstanceID *grpc_application_go.AppInstanceId) (*grpc_application_go.InstanceParameterList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.appClient.ListInstanceParameters(ctx, appInstanceID)
}

// ListDescriptorAppParameters retrieves a list of parameters of an application
func (m *Manager) ListDescriptorAppParameters(appDescriptorID *grpc_application_go.AppDescriptorId) (*grpc_public_api_go.AppParameterList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	params, err := m.appClient.ListDescriptorAppParameters(ctx, appDescriptorID)
	if err != nil {
		return nil, err
	}

	result := make([]*grpc_public_api_go.AppParameter, 0)
	for _, p := range params.Parameters {
		result = append(result, entities.ToPublicAPIAppParam(p))
	}
	return &grpc_public_api_go.AppParameterList{
		Parameters: result,
	}, nil
}
