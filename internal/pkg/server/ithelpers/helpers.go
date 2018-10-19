/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package ithelpers

import (
	"context"
	"fmt"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/onsi/gomega"
)

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