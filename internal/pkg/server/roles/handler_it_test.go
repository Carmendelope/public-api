/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package roles

import (
	"context"
	"fmt"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/public-api/internal/pkg/server/ithelpers"
	"github.com/nalej/public-api/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"os"
)

var _ = ginkgo.Describe("Roles", func() {

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		systemModelAddress = os.Getenv("IT_SM_ADDRESS")
	)

	if systemModelAddress == "" {
		ginkgo.Fail("missing environment variables")
	}

	// gRPC server
	var server * grpc.Server
	// grpc test listener
	var listener * bufconn.Listener
	// client
	var orgClient grpc_organization_go.OrganizationsClient
	var smConn * grpc.ClientConn
	var client grpc_public_api_go.RolesClient

	// Target organization.
	var targetOrganization * grpc_organization_go.Organization

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		smConn = utils.GetConnection(systemModelAddress)
		orgClient = grpc_organization_go.NewOrganizationsClient(smConn)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		manager := NewManager()
		handler := NewHandler(manager)
		grpc_public_api_go.RegisterRolesServer(server, handler)
		test.LaunchServer(server, listener)

		client = grpc_public_api_go.NewRolesClient(conn)
		targetOrganization = ithelpers.CreateOrganization(fmt.Sprintf("testOrg-%d", ginkgo.GinkgoRandomSeed()), orgClient)

	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
		smConn.Close()
	})

	ginkgo.PIt("should be able to list the roles in the system", func(){
		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId:       targetOrganization.OrganizationId,
		}
	    roleList, err := client.List(context.Background(), organizationID)
	    gomega.Expect(err).To(gomega.Succeed())
	    gomega.Expect(len(roleList.Roles)).Should(gomega.Equal(1))
	})

})