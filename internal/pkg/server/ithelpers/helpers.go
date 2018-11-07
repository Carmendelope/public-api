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
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"
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