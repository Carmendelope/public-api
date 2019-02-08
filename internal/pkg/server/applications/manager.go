package applications

import (
	"context"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-conductor-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/public-api/internal/pkg/entities"
)

type Manager struct {
	appClient grpc_application_manager_go.ApplicationManagerClient
}

func NewManager(appClient grpc_application_manager_go.ApplicationManagerClient) Manager {
	return Manager{appClient}
}

// AddAppDescriptor adds a new application descriptor to a given organization.
func (m *Manager) AddAppDescriptor(addRequest *grpc_application_go.AddAppDescriptorRequest) (*grpc_application_go.AppDescriptor, error) {
	return m.appClient.AddAppDescriptor(context.Background(), addRequest)
}

// ListAppDescriptors retrieves a list of application descriptors.
func (m *Manager) ListAppDescriptors(organizationID *grpc_organization_go.OrganizationId) (*grpc_application_go.AppDescriptorList, error) {
	return m.appClient.ListAppDescriptors(context.Background(), organizationID)
}

// GetAppDescriptor retrieves a given application descriptor.
func (m *Manager) GetAppDescriptor(appDescriptorID *grpc_application_go.AppDescriptorId) (*grpc_application_go.AppDescriptor, error) {
	return m.appClient.GetAppDescriptor(context.Background(), appDescriptorID)
}

// UpdateAppDescriptor allows the user to update the information of a registered descriptor.
func (m *Manager) UpdateAppDescriptor(request *grpc_application_go.UpdateAppDescriptorRequest) (*grpc_application_go.AppDescriptor, error) {
	return m.appClient.UpdateAppDescriptor(context.Background(), request)
}

// DeleteAppDescriptor deletes a given application descriptor.
func (m *Manager) DeleteAppDescriptor(appDescriptorID *grpc_application_go.AppDescriptorId) (*grpc_common_go.Success, error) {
	return m.appClient.RemoveAppDescriptor(context.Background(), appDescriptorID)
}

// Deploy an application descriptor.
func (m *Manager) Deploy(deployRequest *grpc_application_manager_go.DeployRequest) (*grpc_conductor_go.DeploymentResponse, error) {
	return m.appClient.Deploy(context.Background(), deployRequest)
}

// Undeploy a running application instance.
func (m *Manager) Undeploy(appInstanceID *grpc_application_go.AppInstanceId) (*grpc_common_go.Success, error) {
	return m.appClient.Undeploy(context.Background(), appInstanceID)
}

// ListAppInstances retrieves a list of application descriptors.
func (m *Manager) ListAppInstances(organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.AppInstanceList, error) {
	apps, err := m.appClient.ListAppInstances(context.Background(), organizationID)
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
	inst, err := m.appClient.GetAppInstance(context.Background(), appInstanceID)
	if err != nil {
		return nil, err
	}
	return entities.ToPublicAPIAppInstance(inst), nil
}
