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

package devices

import (
	"context"
	"fmt"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-device-controller-go"
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

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		systemModelAddress   = os.Getenv("IT_SM_ADDRESS")
		deviceManagerAddress = os.Getenv("IT_DEVICE_MANAGER_ADDRESS")
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
	var latClient grpc_device_manager_go.LatencyClient
	var smConn *grpc.ClientConn
	var dmConn *grpc.ClientConn
	var latConn *grpc.ClientConn
	var client grpc_public_api_go.DevicesClient

	// Target organization.
	var targetOrganization *grpc_organization_go.Organization
	//var targetDeviceGroup *grpc_device_manager_go.DeviceGroup
	var ownerToken string
	var devManagerToken string
	var profileToken string

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()

		server = grpc.NewServer(interceptor.WithServerAuthxInterceptor(
			interceptor.NewConfig(ithelpers.GetAllAuthConfig(), "secret", ithelpers.AuthHeader)))

		smConn = utils.GetConnection(systemModelAddress)
		dmConn = utils.GetConnection(deviceManagerAddress)
		latConn = utils.GetConnection(deviceManagerAddress)
		orgClient = grpc_organization_go.NewOrganizationsClient(smConn)
		dmClient = grpc_device_manager_go.NewDevicesClient(dmConn)
		latClient = grpc_device_manager_go.NewLatencyClient(latConn)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		manager := NewManager(dmClient)
		handler := NewHandler(manager)

		grpc_public_api_go.RegisterDevicesServer(server, handler)
		test.LaunchServer(server, listener)

		client = grpc_public_api_go.NewDevicesClient(conn)
		rand.Seed(ginkgo.GinkgoRandomSeed())
		targetOrganization = ithelpers.CreateOrganization(fmt.Sprintf("testOrg-%d", rand.Int()), orgClient)
		//targetDeviceGroup = ithelpers.CreateDeviceGroup(targetOrganization.OrganizationId, fmt.Sprintf("testDG-%d", rand.Int()), dmClient)
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
	ginkgo.Context("Device Groups", func() {

		ginkgo.It("should be able to add a device group", func() {
			tests := make([]utils.TestResult, 0)
			tests = append(tests, utils.TestResult{Token: ownerToken, Success: true, Msg: "Owner should be able to create a device group"})
			tests = append(tests, utils.TestResult{Token: devManagerToken, Success: true, Msg: "Device Manager should be able to create a device group"})
			tests = append(tests, utils.TestResult{Token: profileToken, Success: false, Msg: "Profile user should NOT be able to create a device group"})

			addRequest := &grpc_device_manager_go.AddDeviceGroupRequest{
				OrganizationId:            targetOrganization.OrganizationId,
				Name:                      fmt.Sprintf("dg-%d", rand.Int()),
				Enabled:                   false,
				DefaultDeviceConnectivity: false,
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

		ginkgo.It("should be able to remove a device group", func() {
			tests := make([]utils.TestResult, 0)
			tests = append(tests, utils.TestResult{Token: ownerToken, Success: true, Msg: "Owner should be able to remove a device group"})
			tests = append(tests, utils.TestResult{Token: devManagerToken, Success: true, Msg: "Device Manager should be able to remove a device group"})
			tests = append(tests, utils.TestResult{Token: profileToken, Success: false, Msg: "Profile user should NOT be able to remove a device group"})

			addRequest := &grpc_device_manager_go.AddDeviceGroupRequest{
				OrganizationId:            targetOrganization.OrganizationId,
				Name:                      fmt.Sprintf("dg-%d", rand.Int()),
				Enabled:                   false,
				DefaultDeviceConnectivity: false,
			}

			for _, test := range tests {

				// 1) Add a device group (owner)
				ctxOwner, cancel := ithelpers.GetContext(ownerToken)
				defer cancel()
				added, err := client.AddDeviceGroup(ctxOwner, addRequest)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(added.DeviceGroupId).ShouldNot(gomega.BeEmpty())
				gomega.Expect(added.DeviceGroupApiKey).ShouldNot(gomega.BeEmpty())

				// 2) remove
				ctx, cancel := ithelpers.GetContext(test.Token)
				defer cancel()

				removeGroup := &grpc_device_go.DeviceGroupId{
					OrganizationId: added.OrganizationId,
					DeviceGroupId:  added.DeviceGroupId,
				}
				success, err := client.RemoveDeviceGroup(ctx, removeGroup)

				if test.Success {
					gomega.Expect(err).To(gomega.Succeed())
					gomega.Expect(success).ShouldNot(gomega.BeNil())
				} else {
					gomega.Expect(err).NotTo(gomega.Succeed())
					gomega.Expect(success).Should(gomega.BeNil())

				}
			}
		})

		ginkgo.It("should be able to update a device group", func() {
			tests := make([]utils.TestResult, 0)
			tests = append(tests, utils.TestResult{Token: ownerToken, Success: true, Msg: "Owner should be able to update a device group"})
			tests = append(tests, utils.TestResult{Token: devManagerToken, Success: true, Msg: "Device Manager should be able to update a device group"})
			tests = append(tests, utils.TestResult{Token: profileToken, Success: false, Msg: "Profile user should NOT be able to update a device group"})

			addRequest := &grpc_device_manager_go.AddDeviceGroupRequest{
				OrganizationId:            targetOrganization.OrganizationId,
				Name:                      fmt.Sprintf("dg-%d", rand.Int()),
				Enabled:                   false,
				DefaultDeviceConnectivity: false,
			}

			for _, test := range tests {

				// 1) Add a device group (owner)
				ctxOwner, cancel := ithelpers.GetContext(ownerToken)
				defer cancel()
				added, err := client.AddDeviceGroup(ctxOwner, addRequest)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(added.DeviceGroupId).ShouldNot(gomega.BeEmpty())
				gomega.Expect(added.DeviceGroupApiKey).ShouldNot(gomega.BeEmpty())

				// 2) update
				ctx, cancel := ithelpers.GetContext(test.Token)
				defer cancel()

				updateGroup := &grpc_device_manager_go.UpdateDeviceGroupRequest{
					OrganizationId: added.OrganizationId,
					DeviceGroupId:  added.DeviceGroupId,
					UpdateEnabled:  true,
					Enabled:        true,
				}
				updated, err := client.UpdateDeviceGroup(ctx, updateGroup)

				if test.Success {
					gomega.Expect(err).To(gomega.Succeed())
					gomega.Expect(updated).ShouldNot(gomega.BeNil())
					gomega.Expect(updated.Enabled).Should(gomega.Equal(updateGroup.Enabled))
					gomega.Expect(updated.DefaultDeviceConnectivity).Should(gomega.Equal(added.DefaultDeviceConnectivity))
				} else {
					gomega.Expect(err).NotTo(gomega.Succeed())

				}
			}
		})

		ginkgo.It("should be able to list a device groups on an organization", func() {
			tests := make([]utils.TestResult, 0)
			tests = append(tests, utils.TestResult{Token: ownerToken, Success: true, Msg: "Owner should be able to list a device groups on an organization"})
			tests = append(tests, utils.TestResult{Token: devManagerToken, Success: true, Msg: "Device Manager should be able to list a device groups on an organization"})
			tests = append(tests, utils.TestResult{Token: profileToken, Success: false, Msg: "Profile user should NOT be able to list a device groups on an organization"})

			addRequest := &grpc_device_manager_go.AddDeviceGroupRequest{
				OrganizationId:            targetOrganization.OrganizationId,
				Name:                      fmt.Sprintf("dg-%d", rand.Int()),
				Enabled:                   false,
				DefaultDeviceConnectivity: false,
			}

			for _, test := range tests {

				// 1) Add a device group (owner)
				ctxOwner, cancel := ithelpers.GetContext(ownerToken)
				defer cancel()
				added, err := client.AddDeviceGroup(ctxOwner, addRequest)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(added.DeviceGroupId).ShouldNot(gomega.BeEmpty())
				gomega.Expect(added.DeviceGroupApiKey).ShouldNot(gomega.BeEmpty())

				// 2) list
				ctx, cancel := ithelpers.GetContext(test.Token)
				defer cancel()

				organizationID := &grpc_organization_go.OrganizationId{
					OrganizationId: added.OrganizationId,
				}
				list, err := client.ListDeviceGroups(ctx, organizationID)

				if test.Success {
					gomega.Expect(err).To(gomega.Succeed())
					gomega.Expect(list.Groups).ShouldNot(gomega.BeEmpty())
				} else {
					gomega.Expect(err).NotTo(gomega.Succeed())

				}
			}
		})

	})
	ginkgo.Context("Devices", func() {

		var targetDeviceGroup *grpc_device_manager_go.DeviceGroup
		var targetDevice *grpc_device_manager_go.RegisterResponse // (DeviceId, DeviceApiKey)

		ginkgo.BeforeEach(func() {
			// create a device group
			targetDeviceGroup = ithelpers.CreateDeviceGroup(targetOrganization.OrganizationId, fmt.Sprintf("testDG-%d", rand.Int()), dmClient)
			// create a device in the group above
			targetDevice = ithelpers.CreateDevice(targetOrganization.OrganizationId,
				targetDeviceGroup.DeviceGroupId,
				targetDeviceGroup.DeviceGroupApiKey,
				dmClient)
		})

		ginkgo.It("should be able to list devices on a group", func() {
			tests := make([]utils.TestResult, 0)
			tests = append(tests, utils.TestResult{Token: ownerToken, Success: true, Msg: "Owner should be able to list devices on a group"})
			tests = append(tests, utils.TestResult{Token: devManagerToken, Success: true, Msg: "Device Manager should be able to list devices on a group"})
			tests = append(tests, utils.TestResult{Token: profileToken, Success: false, Msg: "Profile user should NOT be able to list devices on a group"})

			request := &grpc_device_go.DeviceGroupId{
				OrganizationId: targetOrganization.OrganizationId,
				DeviceGroupId:  targetDeviceGroup.DeviceGroupId,
			}

			for _, test := range tests {
				ctx, cancel := ithelpers.GetContext(test.Token)
				defer cancel()
				list, err := client.ListDevices(ctx, request)
				if test.Success {
					gomega.Expect(err).To(gomega.Succeed())
					gomega.Expect(list.Devices).ShouldNot(gomega.BeEmpty())
				} else {
					gomega.Expect(err).NotTo(gomega.Succeed())
				}
			}
		})

		ginkgo.PIt("should be able to add labels in a device (pending until device-manager implements this)", func() {
			tests := make([]utils.TestResult, 0)
			tests = append(tests, utils.TestResult{Token: ownerToken, Success: true, Msg: "Owner should be able to add labels in a device"})
			tests = append(tests, utils.TestResult{Token: devManagerToken, Success: true, Msg: "Device Manager should be able to add labels in a device"})
			tests = append(tests, utils.TestResult{Token: profileToken, Success: false, Msg: "Profile user should NOT be able to add labels in a device"})

			for _, test := range tests {

				tam := rand.Intn(5) + 1

				request := &grpc_device_manager_go.DeviceLabelRequest{
					OrganizationId: targetOrganization.OrganizationId,
					DeviceGroupId:  targetDeviceGroup.DeviceGroupId,
					DeviceId:       targetDevice.DeviceId,
					Labels:         ithelpers.GenerateLabels(tam),
				}
				ctx, cancel := ithelpers.GetContext(test.Token)
				defer cancel()
				success, err := client.AddLabelToDevice(ctx, request)
				if test.Success {
					gomega.Expect(err).To(gomega.Succeed())
					gomega.Expect(success).NotTo(gomega.BeNil())
				} else {
					gomega.Expect(err).NotTo(gomega.Succeed())
					gomega.Expect(success).To(gomega.BeNil())
				}
			}
		})
		ginkgo.PIt("should be able to remove labels in a device (pending until device-manager implements this)", func() {
		})

		ginkgo.It("should be able to update a device", func() {
			tests := make([]utils.TestResult, 0)
			tests = append(tests, utils.TestResult{Token: ownerToken, Success: true, Msg: "Owner should be able to add labels in a device"})
			tests = append(tests, utils.TestResult{Token: devManagerToken, Success: true, Msg: "Device Manager should be able to add labels in a device"})
			tests = append(tests, utils.TestResult{Token: profileToken, Success: false, Msg: "Profile user should NOT be able to add labels in a device"})

			enabled := !targetDeviceGroup.DefaultDeviceConnectivity

			for _, test := range tests {
				request := &grpc_device_manager_go.UpdateDeviceRequest{
					OrganizationId: targetOrganization.OrganizationId,
					DeviceGroupId:  targetDeviceGroup.DeviceGroupId,
					DeviceId:       targetDevice.DeviceId,
					Enabled:        enabled,
				}

				ctx, cancel := ithelpers.GetContext(test.Token)
				defer cancel()
				device, err := client.UpdateDevice(ctx, request)
				if test.Success {
					gomega.Expect(err).To(gomega.Succeed())
					gomega.Expect(device).NotTo(gomega.BeNil())
					if enabled {
						gomega.Expect(device.Enabled).To(gomega.BeTrue())
					} else {
						gomega.Expect(device.Enabled).NotTo(gomega.BeTrue())
					}
					// change the value
					enabled = !enabled
				} else {
					gomega.Expect(err).NotTo(gomega.Succeed())
				}

			}
		})

		ginkgo.It("should be able to get the device status", func() {

			ping := &grpc_device_controller_go.RegisterLatencyRequest{
				OrganizationId: targetOrganization.OrganizationId,
				DeviceGroupId:  targetDeviceGroup.DeviceGroupId,
				DeviceId:       targetDevice.DeviceId,
				Latency:        20,
			}

			success, err := latClient.RegisterLatency(context.Background(), ping)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(success).NotTo(gomega.BeNil())

			request := &grpc_device_go.DeviceGroupId{
				OrganizationId: targetOrganization.OrganizationId,
				DeviceGroupId:  targetDeviceGroup.DeviceGroupId,
			}

			ctx, cancel := ithelpers.GetContext(ownerToken)
			defer cancel()
			list, err := client.ListDevices(ctx, request)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list.Devices).ShouldNot(gomega.BeEmpty())
			gomega.Expect(list.Devices[0].DeviceStatusName).Should(gomega.Equal(grpc_device_manager_go.DeviceStatus_ONLINE.String()))

		})
	})

})
