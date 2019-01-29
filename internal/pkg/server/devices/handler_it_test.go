/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package devices

import (
	"fmt"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
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
	"math/rand"
	"os"
)

/*
RUN_INTEGRATION_TEST=true
IT_SM_ADDRESS=localhost:8800
IT_DEVICE_MANAGER_ADDRESS=localhost:6010
*/

var _ = ginkgo.Describe("Devices", func() {

	const NumDevices= 10

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		systemModelAddress= os.Getenv("IT_SM_ADDRESS")
		deviceManagerAddress= os.Getenv("IT_DEVICE_MANAGER_ADDRESS")
	)

	if systemModelAddress == "" || deviceManagerAddress == "" {
		ginkgo.Fail("missing environment variables")
	}

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var orgClient grpc_organization_go.OrganizationsClient
	var dmClient grpc_device_manager_go.DevicesClient
	var deviceSMClient grpc_device_go.DevicesClient
	var smConn *grpc.ClientConn
	var dmConn *grpc.ClientConn
	var client grpc_public_api_go.DevicesClient

	// Target organization.
	var targetOrganization *grpc_organization_go.Organization
	var targetDeviceGroup *grpc_device_manager_go.DeviceGroup
	var ownerToken string
	var devManagerToken string
	var profileToken string

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()

		server = grpc.NewServer(interceptor.WithServerAuthxInterceptor(
			interceptor.NewConfig(ithelpers.GetAllAuthConfig(), "secret", ithelpers.AuthHeader)))

		smConn = utils.GetConnection(systemModelAddress)
		dmConn = utils.GetConnection(deviceManagerAddress)
		orgClient = grpc_organization_go.NewOrganizationsClient(smConn)
		dmClient = grpc_device_manager_go.NewDevicesClient(dmConn)
		deviceSMClient = grpc_device_go.NewDevicesClient(smConn)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		manager := NewManager(dmClient)
		handler := NewHandler(manager)
		grpc_public_api_go.RegisterDevicesServer(server, handler)
		test.LaunchServer(server, listener)

		client = grpc_public_api_go.NewDevicesClient(conn)
		rand.Seed(ginkgo.GinkgoRandomSeed())
		targetOrganization = ithelpers.CreateOrganization(fmt.Sprintf("testOrg-%d", rand.Int()), orgClient)
		targetDeviceGroup = ithelpers.CreateDeviceGroup(targetOrganization.OrganizationId, fmt.Sprintf("testDG-%d", rand.Int()), dmClient)
		ownerToken = ithelpers.GenerateToken("email@nalej.com",
			targetOrganization.OrganizationId, "Owner", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG})
		devManagerToken = ithelpers.GenerateToken("devmngr@nalej.com",
			targetOrganization.OrganizationId, "Device Manager", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_PROFILE, grpc_authx_go.AccessPrimitive_DEVMNGR})
		profileToken = ithelpers.GenerateToken("profile@nalej.com",
			targetOrganization.OrganizationId, "Profile", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_PROFILE})
	})

	ginkgo.AfterSuite(func() {

		ithelpers.NewTestCleaner(smConn).DeleteOrganizationClusters(targetOrganization.OrganizationId)

		server.Stop()
		listener.Close()
		smConn.Close()
	})

	ginkgo.It("should be able to add a device group", func(){
		tests := make([]utils.TestResult, 0)
		tests = append(tests, utils.TestResult{Token: ownerToken, Success: true, Msg: "Owner should be able to create a device group"})
		tests = append(tests, utils.TestResult{Token: devManagerToken, Success: true, Msg: "Device Manager should be able to create a device group"})
		tests = append(tests, utils.TestResult{Token: profileToken, Success: false, Msg: "Profile user should NOT be able to create a device group"})

		addRequest := &grpc_device_manager_go.AddDeviceGroupRequest{
			OrganizationId:            targetOrganization.OrganizationId,
			Name:                      fmt.Sprintf("dg-%d", rand.Int()),
			Enabled:                   false,
			DeviceDefaultConnectivity: false,
		}

		for _, test := range tests {
			ctx, cancel := ithelpers.GetContext(test.Token)
			defer cancel()
			added, err := client.AddDeviceGroup(ctx, addRequest)
			if test.Success {
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(added.DeviceGroupId).ShouldNot(gomega.BeEmpty())
				gomega.Expect(added.DeviceGroupApiKey).ShouldNot(gomega.BeEmpty())
			} else {
				gomega.Expect(err).NotTo(gomega.Succeed())
			}
		}
	})

})