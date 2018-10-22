/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */



package users

import (
	"context"
	"fmt"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-go"
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

var _ = ginkgo.Describe("Users", func() {

	const NumUsers = 5

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
	var client grpc_public_api_go.UsersClient

	// Target organization.
	var targetOrganization * grpc_organization_go.Organization
	var targetUser * grpc_user_go.User

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		smConn = utils.GetConnection(systemModelAddress)
		orgClient = grpc_organization_go.NewOrganizationsClient(smConn)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		manager := NewManager()
		handler := NewHandler(manager)
		grpc_public_api_go.RegisterUsersServer(server, handler)
		test.LaunchServer(server, listener)

		client = grpc_public_api_go.NewUsersClient(conn)
		targetOrganization = ithelpers.CreateOrganization(fmt.Sprintf("testOrg-%d", ginkgo.GinkgoRandomSeed()), orgClient)

	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
		smConn.Close()
	})

	ginkgo.PIt("should be able to retrieve the user information", func(){
		userID := &grpc_user_go.UserId{
			OrganizationId:       targetUser.OrganizationId,
			Email:                targetUser.Email,
		}
		info, err := client.Info(context.Background(), userID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(info.OrganizationId).Should(gomega.Equal(targetUser.OrganizationId))
		gomega.Expect(info.Email).Should(gomega.Equal(targetUser.Email))
		gomega.Expect(info.Name).Should(gomega.Equal(targetUser.Name))
		gomega.Expect(info.RoleName).Should(gomega.Equal("RoleName"))
	})

	ginkgo.PIt("should be able list users in an organization", func(){
	    organizationID := &grpc_organization_go.OrganizationId{
	    	OrganizationId: targetOrganization.OrganizationId,
		}
	    list, err := client.List(context.Background(), organizationID)
	    gomega.Expect(err).To(gomega.Succeed())
	    gomega.Expect(len(list.Users)).Should(gomega.Equal(NumUsers))
	})

	ginkgo.PIt("should be able to delete a user", func(){
		userID := &grpc_user_go.UserId{
			OrganizationId:       targetUser.OrganizationId,
			Email:                targetUser.Email,
		}
		success, err := client.Delete(context.Background(), userID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(success).ToNot(gomega.BeNil())
	})

	ginkgo.PIt("should be able to reset the password of a user", func(){
		userID := &grpc_user_go.UserId{
			OrganizationId:       targetUser.OrganizationId,
			Email:                targetUser.Email,
		}
		reset, err := client.ResetPassword(context.Background(), userID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(reset.NewPassword).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.PIt("should be able to update an existing user", func(){
	    updateUserRequest := &grpc_user_go.UpdateUserRequest{
			OrganizationId:       targetUser.OrganizationId,
			Email:                targetUser.Email,
			Name:                 "newName",
			Role:                 "newRole",
		}
		success, err := client.Update(context.Background(), updateUserRequest)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(success).ToNot(gomega.BeNil())
	})
})