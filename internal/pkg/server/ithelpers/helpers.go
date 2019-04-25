/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package ithelpers

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/authx/pkg/token"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"time"
)

const AuthHeader = "authorization"

func CreateOrganization(name string, orgClient grpc_organization_go.OrganizationsClient) *grpc_organization_go.Organization {
	toAdd := &grpc_organization_go.AddOrganizationRequest{
		Name: name,
	}
	added, err := orgClient.AddOrganization(context.Background(), toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	gomega.Expect(added).ToNot(gomega.BeNil())
	return added
}

func CreateCluster(organization *grpc_organization_go.Organization, clusterName string, clustClient grpc_infrastructure_go.ClustersClient) *grpc_infrastructure_go.Cluster {
	toAdd := &grpc_infrastructure_go.AddClusterRequest{
		RequestId:      organization.OrganizationId,
		OrganizationId: organization.OrganizationId,
		Name:           clusterName,
		Hostname:       clusterName,
		Labels:         nil,
	}

	added, err := clustClient.AddCluster(context.Background(), toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return added
}

func CreateNodes(cluster *grpc_infrastructure_go.Cluster, numNodes int, clustClient grpc_infrastructure_go.ClustersClient, nodeClient grpc_infrastructure_go.NodesClient) {
	for i := 0; i < numNodes; i++ {
		toAdd := &grpc_infrastructure_go.AddNodeRequest{
			RequestId:      cluster.ClusterId,
			OrganizationId: cluster.OrganizationId,
			Ip:             fmt.Sprintf("172.168.1.%d", i+1),
			Labels:         nil,
		}
		added, err := nodeClient.AddNode(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		attachRequest := &grpc_infrastructure_go.AttachNodeRequest{
			RequestId:      added.NodeId,
			OrganizationId: cluster.OrganizationId,
			ClusterId:      cluster.ClusterId,
			NodeId:         added.NodeId,
		}
		_, err = nodeClient.AttachNode(context.Background(), attachRequest)
		gomega.Expect(err).To(gomega.Succeed())
	}
}

func CreateRole(organizationID string, userManagerClient grpc_user_manager_go.UserManagerClient) *grpc_authx_go.Role {
	addRoleRequest := &grpc_user_manager_go.AddRoleRequest{
		OrganizationId: organizationID,
		Name:           "test role",
		Description:    "test role",
		Primitives:     []grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG},
	}
	added, err := userManagerClient.AddRole(context.Background(), addRoleRequest)
	gomega.Expect(err).To(gomega.Succeed())
	return added
}

func CreateUser(organizationID string, roleID string, userManagerClient grpc_user_manager_go.UserManagerClient) *grpc_user_manager_go.User {

	addUserRequest := &grpc_user_manager_go.AddUserRequest{
		OrganizationId:       organizationID,
		Email:                fmt.Sprintf("test-%d@mail.com", rand.Int()),
		Password:             "password",
		Name:                 "randomUser",
		RoleId:               roleID,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	added, err := userManagerClient.AddUser(context.Background(), addUserRequest)
	gomega.Expect(err).To(gomega.Succeed())
	return added
}

func GetAddDescriptorRequest(organizationID string) *grpc_application_go.AddAppDescriptorRequest {
	service := &grpc_application_go.Service{
		OrganizationId: organizationID,
		ServiceGroupId: "Sg_1",
		ServiceId:      "1",
		Name:           "simple-mysql-service",
		Type:           grpc_application_go.ServiceType_DOCKER,
		Image:          "mysql:5.6",
		Specs:          &grpc_application_go.DeploySpecs{Replicas: 1},
		Credentials:    &grpc_application_go.ImageCredentials{Username:"user_name", Password:"password", Email:"email@email.es"},
		Storage:        []*grpc_application_go.Storage{&grpc_application_go.Storage{MountPath: "/tmp"}},
		ExposedPorts: []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "mysqlport", InternalPort: 3306, ExposedPort: 3306,
		}},
		EnvironmentVariables: map[string]string{"MYSQL_ROOT_PASSWORD": "root"},
		Configs:              []*grpc_application_go.ConfigFile{&grpc_application_go.ConfigFile{MountPath: "/tmp"}},
		Labels:               map[string]string{"app": "simple-app", "component": "mysql"},
	}

	group1 := &grpc_application_go.ServiceGroup{
		Name: "g1",
		Services: []*grpc_application_go.Service{service},
		Specs: &grpc_application_go.ServiceGroupDeploymentSpecs{Replicas:1,MultiClusterReplica:false},
	}
	secRule := grpc_application_go.SecurityRule{
		Name:            "allow access to mysql",
		Access:          grpc_application_go.PortAccess_PUBLIC,
		RuleId:          "001",
		TargetPort:      3306,
		TargetServiceName: "1",
		TargetServiceGroupName: "g1",
	}

	return &grpc_application_go.AddAppDescriptorRequest{
		RequestId:      GenerateUUID(),
		OrganizationId: organizationID,
		Name:           "Sample application",
		Labels:         map[string]string{"app": "simple-app"},
		Rules:          []*grpc_application_go.SecurityRule{&secRule},
		Groups:      []*grpc_application_go.ServiceGroup{group1},
	}
}

func CreateAppDescriptor(organizationID string, appClient grpc_application_manager_go.ApplicationManagerClient) *grpc_application_go.AppDescriptor {
	toAdd := GetAddDescriptorRequest(organizationID)
	added, err := appClient.AddAppDescriptor(context.Background(), toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return added
}

func GenerateDeploy(organizationID string, appDescriptorID string) *grpc_application_manager_go.DeployRequest {
	return &grpc_application_manager_go.DeployRequest{
		OrganizationId:  organizationID,
		AppDescriptorId: appDescriptorID,
		Name:            GenerateUUID(),
	}
}

// GenerateUUID creates a new random UUID.
func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateToken(email string, organizationID string, roleName string, secret string, primitives []grpc_authx_go.AccessPrimitive) string {
	p := make([]string, 0)
	for _, prim := range primitives {
		p = append(p, prim.String())
	}

	pClaim := token.PersonalClaim{
		UserID:         email,
		Primitives:     p,
		RoleName:       roleName,
		OrganizationID: organizationID,
	}

	claim := token.NewClaim(pClaim, "it", time.Now(), time.Minute)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := t.SignedString([]byte(secret))
	gomega.Expect(err).To(gomega.Succeed())

	return tokenString
}

func GetContext(token string) (context.Context, context.CancelFunc) {
	md := metadata.New(map[string]string{AuthHeader: token})
	baseContext, cancel := context.WithTimeout(context.Background(), time.Minute)
	return metadata.NewOutgoingContext(baseContext, md), cancel
}

func GetAuthConfig(endpoints ...string) *interceptor.AuthorizationConfig {
	permissions := make(map[string]interceptor.Permission, 0)

	for _, e := range endpoints {
		permissions[e] = interceptor.Permission{
			Must: []string{grpc_authx_go.AccessPrimitive_ORG.String()},
		}
	}

	return &interceptor.AuthorizationConfig{
		AllowsAll:   false,
		Permissions: permissions,
	}
}

func GetAllAuthConfig() *interceptor.AuthorizationConfig {
	permissions := make(map[string]interceptor.Permission, 0)
	permissions["/public_api.Clusters/Install"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_RESOURCES.String()},
	}
	permissions["/public_api.Clusters/Info"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_RESOURCES.String()},
	}
	permissions["/public_api.Clusters/List"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_RESOURCES.String()},
	}
	permissions["/public_api.Clusters/Update"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_RESOURCES.String()},
	}
	permissions["/public_api.Nodes/List"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_RESOURCES.String()},
	}
	permissions["/public_api.Organizations/Info"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_RESOURCES.String(), grpc_authx_go.AccessPrimitive_APPS.String()},
	}
	permissions["/public_api.Resources/Summary"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_RESOURCES.String()},
	}
	permissions["/public_api.Roles/List"] = interceptor.Permission{
		Must: []string{grpc_authx_go.AccessPrimitive_ORG.String()},
	}
	permissions["/public_api.Roles/ListInternal"] = interceptor.Permission{
		Must: []string{grpc_authx_go.AccessPrimitive_ORG.String()},
	}
	permissions["/public_api.Roles/AssignRole"] = interceptor.Permission{
		Must: []string{grpc_authx_go.AccessPrimitive_ORG.String()},
	}
	permissions["/public_api.Users/Add"] = interceptor.Permission{
		Must: []string{grpc_authx_go.AccessPrimitive_ORG.String()},
	}
	permissions["/public_api.Users/Info"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_PROFILE.String()},
	}
	permissions["/public_api.Users/List"] = interceptor.Permission{
		Must: []string{grpc_authx_go.AccessPrimitive_ORG.String()},
	}
	permissions["/public_api.Users/Delete"] = interceptor.Permission{
		Must: []string{grpc_authx_go.AccessPrimitive_ORG.String()},
	}
	permissions["/public_api.Users/ChangePassword"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_PROFILE.String()},
	}
	permissions["/public_api.Users/Update"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_PROFILE.String()},
	}
	permissions["/public_api.Applications/AddAppDescriptor"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_APPS.String()},
	}
	permissions["/public_api.Applications/ListAppDescriptors"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_APPS.String()},
	}
	permissions["/public_api.Applications/GetAppDescriptor"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_APPS.String()},
	}
	permissions["/public_api.Applications/DeleteAppDescriptor"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_APPS.String()},
	}
	permissions["/public_api.Applications/Deploy"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_APPS.String()},
	}
	permissions["/public_api.Applications/Undeploy"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_APPS.String()},
	}
	permissions["/public_api.Applications/ListAppInstances"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_APPS.String()},
	}
	permissions["/public_api.Applications/GetAppInstance"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_APPS.String()},
	}
	permissions["/public_api.UnifiedLogging/Search"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_APPS.String()},
	}
	permissions["/public_api.Devices/AddDeviceGroup"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_DEVMNGR.String()},
	}
	permissions["/public_api.Devices/AddDeviceGroup"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_DEVMNGR.String()},
	}
	permissions["/public_api.Devices/UpdateDeviceGroup"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_DEVMNGR.String()},
	}
	permissions["/public_api.Devices/RemoveDeviceGroup"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_DEVMNGR.String()},
	}
	permissions["/public_api.Devices/ListDeviceGroups"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_DEVMNGR.String()},
	}
	permissions["/public_api.Devices/ListDevices"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_DEVMNGR.String()},
	}
	permissions["/public_api.Devices/AddLabelToDevice"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_DEVMNGR.String()},
	}
	permissions["/public_api.Devices/RemoveLabelFromDevice"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_DEVMNGR.String()},
	}
	permissions["/public_api.Devices/UpdateDevice"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_DEVMNGR.String()},
	}
	permissions["/public_api.Devices/RemoveDevice"] = interceptor.Permission{
		Should: []string{grpc_authx_go.AccessPrimitive_ORG.String(), grpc_authx_go.AccessPrimitive_DEVMNGR.String()},
	}
	return &interceptor.AuthorizationConfig{
		AllowsAll:   false,
		Permissions: permissions,
	}

}

func CreateDeviceGroup(organizationID string, name string, dmClient grpc_device_manager_go.DevicesClient) *grpc_device_manager_go.DeviceGroup{
	request := &grpc_device_manager_go.AddDeviceGroupRequest{
		OrganizationId:            organizationID,
		Name:                      name,
		Enabled:                   true,
		DefaultDeviceConnectivity: true,
	}
	added, err := dmClient.AddDeviceGroup(context.Background(), request)
	gomega.Expect(err).To(gomega.Succeed())
	return added
}

func CreateDevice (organizationID string, deviceGroupID string, groupApiKey string,
	devClient grpc_device_manager_go.DevicesClient) *grpc_device_manager_go.RegisterResponse {
	request := &grpc_device_manager_go.RegisterDeviceRequest{
		OrganizationId: organizationID,
		DeviceGroupId: deviceGroupID,
		DeviceGroupApiKey: groupApiKey,
		DeviceId: GenerateUUID(),
	}
	added, err := devClient.RegisterDevice(context.Background(), request)
	gomega.Expect(err).To(gomega.Succeed())
	return added
}

func GenerateLabels (tam int) map[string]string {
	labels := make (map[string]string, tam)
	for i:= 0; i< tam; i ++ {
		labels[fmt.Sprintf("label_%d", i)] = fmt.Sprintf("value_%d", i)
	}
	return labels
}

// DeleteAllInstances from system model.
func DeleteAllInstances(organizationID string, smAppClient grpc_application_go.ApplicationsClient) {
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}
	instances, err := smAppClient.ListAppInstances(context.Background(), orgID)
	gomega.Expect(err).To(gomega.Succeed())
	for _, inst := range instances.Instances {

		instance := &grpc_application_go.AppInstanceId{
			OrganizationId: organizationID,
			AppInstanceId:  inst.AppInstanceId,
		}
		// TODO: ask Dani if I have to ask for the instance another time
		inst2, err2 := smAppClient.GetAppInstance(context.Background(), instance)
		gomega.Expect(err2).To(gomega.Succeed())

		if inst2.Status == grpc_application_go.ApplicationStatus_QUEUED {
			log.Debug().Str("app_instance", inst.AppInstanceId).Msg("QUEUED, Waiting 3s to DEPLOY")
			time.Sleep(time.Duration(3) * time.Second)
		}

		toRemove := &grpc_application_go.AppInstanceId{
			OrganizationId: inst.OrganizationId,
			AppInstanceId:  inst.AppInstanceId,
		}
		_, err := smAppClient.RemoveAppInstance(context.Background(), toRemove)
		gomega.Expect(err).To(gomega.Succeed())
	}
}

// DeleteAllUsers from system model
func DeleteAllUsers(organizationID string, smUserClient grpc_user_manager_go.UserManagerClient) {
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}
	users, err := smUserClient.ListUsers(context.Background(), orgID)
	gomega.Expect(err).To(gomega.Succeed())

	for _, user := range users.Users {
		toRemove := &grpc_user_go.UserId{
			OrganizationId: user.OrganizationId,
			Email:          user.Email,
		}
		_, err := smUserClient.RemoveUser(context.Background(), toRemove)
		gomega.Expect(err).To(gomega.Succeed())
	}
}

type TestCleaner struct {
	client *grpc.ClientConn
}

func NewTestCleaner(smConn *grpc.ClientConn) *TestCleaner {
	return &TestCleaner{smConn}
}

// DeleteAllInstances from system model
func (tu *TestCleaner) DeleteAllInstances(organizationID string) {
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}

	smAppClient := grpc_application_go.NewApplicationsClient(tu.client)

	instances, err := smAppClient.ListAppInstances(context.Background(), orgID)
	gomega.Expect(err).To(gomega.Succeed())

	log.Debug().Msg("Waiting 2s to conductor queue")
	time.Sleep(time.Duration(2) * time.Second)

	for _, inst := range instances.Instances {

		toRemove := &grpc_application_go.AppInstanceId{
			OrganizationId: inst.OrganizationId,
			AppInstanceId:  inst.AppInstanceId,
		}
		_, err := smAppClient.RemoveAppInstance(context.Background(), toRemove)
		gomega.Expect(err).To(gomega.Succeed())
	}
}

// DeleteAppDescriptors from system model
func (tu *TestCleaner) DeleteAppDescriptors(organizationID string) {

	smAppClient := grpc_application_go.NewApplicationsClient(tu.client)

	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}
	apps, err := smAppClient.ListAppDescriptors(context.Background(), orgID)
	gomega.Expect(err).To(gomega.Succeed())

	if apps.Descriptors != nil {
		for _, app := range apps.Descriptors {

			toRemove := &grpc_application_go.AppDescriptorId{
				OrganizationId:  app.OrganizationId,
				AppDescriptorId: app.AppDescriptorId,
			}

			_, err := smAppClient.RemoveAppDescriptor(context.Background(), toRemove)
			gomega.Expect(err).To(gomega.Succeed())

		}
	}

}

// DeleteAllUsers from system model
func (tu *TestCleaner) DeleteAllUsers(organizationID string) {

	client := grpc_user_manager_go.NewUserManagerClient(tu.client)

	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}
	users, err := client.ListUsers(context.Background(), orgID)
	gomega.Expect(err).To(gomega.Succeed())

	for _, user := range users.Users {
		toRemove := &grpc_user_go.UserId{
			OrganizationId: user.OrganizationId,
			Email:          user.Email,
		}
		_, err := client.RemoveUser(context.Background(), toRemove)
		gomega.Expect(err).To(gomega.Succeed())
	}
}

// Delete the nodes of a cluster from system model
func (tu *TestCleaner) DeleteClusterNodes(organizationID string, clusterID string) {

	client := grpc_infrastructure_go.NewNodesClient(tu.client)

	clusterId := &grpc_infrastructure_go.ClusterId{
		ClusterId:      clusterID,
		OrganizationId: organizationID,
	}
	nodes, err := client.ListNodes(context.Background(), clusterId)
	gomega.Expect(err).To(gomega.Succeed())

	if nodes != nil {
		nodeList := make([]string, 0)
		for _, node := range nodes.Nodes {
			nodeList = append(nodeList, node.NodeId)
		}

		toRemove := &grpc_infrastructure_go.RemoveNodesRequest{
			RequestId:      GenerateUUID(),
			OrganizationId: organizationID,
			Nodes:          nodeList,
		}
		_, err = client.RemoveNodes(context.Background(), toRemove)
		gomega.Expect(err).To(gomega.Succeed())
	}

}

// delete a cluster from system model
func (tu *TestCleaner) DeleteOrganizationClusters(organizationID string) {

	client := grpc_infrastructure_go.NewClustersClient(tu.client)

	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}
	clusters, err := client.ListClusters(context.Background(), orgID)
	gomega.Expect(err).To(gomega.Succeed())

	for _, cluster := range clusters.Clusters {

		tu.DeleteClusterNodes(organizationID, cluster.ClusterId)

		toRemove := &grpc_infrastructure_go.RemoveClusterRequest{
			OrganizationId: cluster.OrganizationId,
			ClusterId:      cluster.ClusterId,
		}

		_, err := client.RemoveCluster(context.Background(), toRemove)
		gomega.Expect(err).To(gomega.Succeed())
	}

}

func (tu *TestCleaner) DeleteOrganizationDescriptors(organizationID string) {
	// delete Instances
	tu.DeleteAllInstances(organizationID)
	// delete appDescriptors
	tu.DeleteAppDescriptors(organizationID)

}

// delete all of an organization
func (tu *TestCleaner) DeleteOrganization(organizationID string) {

	//
	tu.DeleteOrganizationDescriptors(organizationID)
	tu.DeleteOrganizationClusters(organizationID)
	tu.DeleteAllUsers(organizationID)
}
