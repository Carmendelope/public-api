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
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"time"
)

const AuthHeader = "authorization"

func CreateOrganization(name string, orgClient grpc_organization_go.OrganizationsClient) * grpc_organization_go.Organization {
	toAdd := &grpc_organization_go.AddOrganizationRequest{
		Name:                 name,
	}
	added, err := orgClient.AddOrganization(context.Background(), toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	gomega.Expect(added).ToNot(gomega.BeNil())
	return added
}

func CreateCluster(organization * grpc_organization_go.Organization, clusterName string, clustClient grpc_infrastructure_go.ClustersClient) * grpc_infrastructure_go.Cluster {
	toAdd := &grpc_infrastructure_go.AddClusterRequest{
		RequestId:            organization.OrganizationId,
		OrganizationId:       organization.OrganizationId,
		Name:                 clusterName,
		Description:          clusterName,
		Hostname:             clusterName,
		Labels:               nil,
	}

	added, err := clustClient.AddCluster(context.Background(), toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return added
}

func CreateNodes(cluster * grpc_infrastructure_go.Cluster, numNodes int, clustClient grpc_infrastructure_go.ClustersClient, nodeClient grpc_infrastructure_go.NodesClient) {
	for i := 0; i < numNodes; i++ {
		toAdd := &grpc_infrastructure_go.AddNodeRequest{
			RequestId:            cluster.ClusterId,
			OrganizationId:       cluster.OrganizationId,
			Ip:                   fmt.Sprintf("172.168.1.%d", i+1),
			Labels:               nil,
		}
		added, err := nodeClient.AddNode(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		attachRequest := &grpc_infrastructure_go.AttachNodeRequest{
			RequestId:            added.NodeId,
			OrganizationId:       cluster.OrganizationId,
			ClusterId:            cluster.ClusterId,
			NodeId:               added.NodeId,
		}
		_, err = nodeClient.AttachNode(context.Background(), attachRequest)
		gomega.Expect(err).To(gomega.Succeed())
	}
}

func CreateRole(organizationID string, userManagerClient grpc_user_manager_go.UserManagerClient) *grpc_authx_go.Role {
	addRoleRequest := &grpc_user_manager_go.AddRoleRequest{
		OrganizationId:       organizationID,
		Name:                 "test role",
		Description:          "test role",
		Primitives:           []grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG},
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
		ServiceId:            "1",
		Name:                 "Simple MySQL service",
		Description:          "A MySQL instance",
		Type:                 grpc_application_go.ServiceType_DOCKER,
		Image:                "mysql:5.6",
		Specs:                &grpc_application_go.DeploySpecs{Replicas: 1},
		Storage:              []*grpc_application_go.Storage{&grpc_application_go.Storage{MountPath: "/tmp",}},
		ExposedPorts:         []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "mysqlport", InternalPort: 3306, ExposedPort: 3306,
		}},
		EnvironmentVariables: map[string]string{"MYSQL_ROOT_PASSWORD":"root"},
		Configs:              []*grpc_application_go.ConfigFile{&grpc_application_go.ConfigFile{MountPath:"/tmp"}},
		Labels:                map[string]string { "app":"simple-app", "component":"mysql"},
	}

	secRule := grpc_application_go.SecurityRule{
		OrganizationId: organizationID,
		Name:"all open",
		Access: grpc_application_go.PortAccess_PUBLIC,
		RuleId: "001",
		SourcePort: 3306,
		SourceServiceId: "1",
	}

	return &grpc_application_go.AddAppDescriptorRequest{
		RequestId: GenerateUUID(),
		OrganizationId: organizationID,
		Name:                 "Sample application",
		Description:          "This is a basic descriptor of an application",
		Labels:               map[string]string{"app":"simple-app"},
		Rules:                []*grpc_application_go.SecurityRule{&secRule},
		Services:             []*grpc_application_go.Service{service},
	}
}

func CreateAppDescriptor(organizationID string, appClient grpc_application_manager_go.ApplicationManagerClient) * grpc_application_go.AppDescriptor {
	toAdd := GetAddDescriptorRequest(organizationID)
	added, err := appClient.AddAppDescriptor(context.Background(), toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return added
}

// GenerateUUID creates a new random UUID.
func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateToken(email string, organizationID string, roleName string, secret string, primitives []grpc_authx_go.AccessPrimitive) string {
	p := make([]string, 0)
	for _, prim := range primitives{
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

func GetAuthConfig(endpoints ... string) *interceptor.AuthorizationConfig {
	permissions := make(map[string]interceptor.Permission, 0)
	for _, e := range endpoints{
		permissions[e] = interceptor.Permission{
			Must:    []string{grpc_authx_go.AccessPrimitive_ORG.String()},
		}
	}
	return &interceptor.AuthorizationConfig{
		AllowsAll:   false,
		Permissions: permissions,
	}
}