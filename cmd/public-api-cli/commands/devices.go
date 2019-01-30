/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var deviceGroupCmd = &cobra.Command{
	Use:     "devicegroup",
	Aliases: []string{"dg"},
	Short:   "Manage device groups",
	Long:    `Manage device groups`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var devicesCmd = &cobra.Command{
	Use:     "devices",
	Aliases: []string{"device", "dev"},
	Short:   "Manage devices",
	Long:    `Manage devices`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	// Device groups
	rootCmd.AddCommand(deviceGroupCmd)
	addDeviceGroupCmd.Flags().StringVar(&name, "name", "", "Device group name")
	addDeviceGroupCmd.Flags().BoolVar(&enabled, "enabled", false, "Whether the group is enabled")
	addDeviceGroupCmd.Flags().BoolVar(&disabled, "disabled", false, "Whether the group is disabled")
	addDeviceGroupCmd.Flags().BoolVar(&enabledDefaultConnectivity, "enabledDefaultConnectivity", false, "Default connectivity for devices joining the device group (enabled)")
	addDeviceGroupCmd.Flags().BoolVar(&disabledDefaultConnectivity, "disabledDefaultConnectivity", false, "Default connectivity for devices joining the device group (disabled)")
	deviceGroupCmd.AddCommand(addDeviceGroupCmd)

	updateDeviceGroupCmd.Flags().StringVar(&deviceGroupID, "deviceGroupId", "", "Device group identifier")
	updateDeviceGroupCmd.Flags().BoolVar(&enabled, "enabled", false, "Whether the group is enabled")
	updateDeviceGroupCmd.Flags().BoolVar(&disabled, "disabled", false, "Whether the group is disabled")
	updateDeviceGroupCmd.Flags().BoolVar(&enabledDefaultConnectivity, "enabledDefaultConnectivity", false, "Default connectivity for devices joining the device group (enabled)")
	updateDeviceGroupCmd.Flags().BoolVar(&disabledDefaultConnectivity, "disabledDefaultConnectivity", false, "Default connectivity for devices joining the device group (disabled)")
	deviceGroupCmd.AddCommand(updateDeviceGroupCmd)

	removeDeviceGroupCmd.Flags().StringVar(&deviceGroupID, "deviceGroupId", "", "Device group identifier")
	deviceGroupCmd.AddCommand(removeDeviceGroupCmd)

	deviceGroupCmd.AddCommand(listDeviceGroupsCmd)

	// Devices
	rootCmd.AddCommand(devicesCmd)
	devicesCmd.PersistentFlags().StringVar(&deviceGroupID, "deviceGroupId", "", "Device group identifier")

	devicesCmd.AddCommand(listDevicesCmd)

	devicesCmd.AddCommand(deviceLabelsCmd)
	deviceLabelsCmd.PersistentFlags().StringVar(&deviceID, "deviceId", "", "Device identifier")
	deviceLabelsCmd.PersistentFlags().StringVar(&rawLabels, "labels", "", "Labels separated by ; as in key1:value;key2:value")

	deviceLabelsCmd.AddCommand(addLabelToDeviceCmd)
	deviceLabelsCmd.AddCommand(removeLabelFromDeviceCmd)

	updateDeviceCmd.Flags().StringVar(&deviceID, "deviceId", "", "Device identifier")
	updateDeviceCmd.Flags().BoolVar(&enabled, "enabled", false, "Whether the device is enabled")
	updateDeviceCmd.Flags().BoolVar(&disabled, "disabled", false, "Whether the device is disabled")
	devicesCmd.AddCommand(updateDeviceCmd)

}

var addDeviceGroupCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new device group",
	Long:  `Add a new device group to an organization`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		n.AddDeviceGroup(options.Resolve("organizationID", organizationID),
			name, enabled, disabled, enabledDefaultConnectivity, disabledDefaultConnectivity)
	},
}

var updateDeviceGroupCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a device group",
	Long:  `Update the options of a device group`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		n.UpdateDeviceGroup(options.Resolve("organizationID", organizationID),
			deviceGroupID, enabled, disabled, enabledDefaultConnectivity, disabledDefaultConnectivity)
	},
}

var removeDeviceGroupCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a device group",
	Long:  `Remove a device group`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		n.RemoveDeviceGroup(options.Resolve("organizationID", organizationID),
			deviceGroupID)
	},
}

var listDeviceGroupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List the device groups in an organization",
	Long:  `List the device groups in an organization`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		n.ListDeviceGroups(options.Resolve("organizationID", organizationID))
	},
}

var listDevicesCmd = &cobra.Command{
	Use:   "list",
	Short: "List the devices in a device group",
	Long:  `List the devices in a device group`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		n.ListDevices(options.Resolve("organizationID", organizationID),
			deviceGroupID)
	},
}

var deviceLabelsCmd = &cobra.Command{
	Use:   "label",
	Aliases: []string{"labels", "l"},
	Short: "Manage device labels",
	Long:  `Manage device labels`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var addLabelToDeviceCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a set of labels to a device",
	Long:  `Add a set of labels to a device`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		n.AddLabelToDevice(options.Resolve("organizationID", organizationID),
			deviceGroupID, deviceID, rawLabels)
	},
}

var removeLabelFromDeviceCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a set of labels to a device",
	Long:  `Remove a set of labels to a device`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		n.RemoveLabelFromDevice(options.Resolve("organizationID", organizationID),
			deviceGroupID, deviceID, rawLabels)
	},
}

var updateDeviceCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the information of a device",
	Long:  `Update the information of a device`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		n.UpdateDevice(options.Resolve("organizationID", organizationID),
			deviceGroupID, deviceID, enabled, disabled)
	},
}