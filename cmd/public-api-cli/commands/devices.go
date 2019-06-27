/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package commands

import (
	"fmt"
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var deviceGroupCmd = &cobra.Command{
	Use:     "devicegroups",
	Aliases: []string{"devicegroup", "dg"},
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
	updateDeviceGroupCmd.Flags().BoolVar(&enabled, "enable", false, "Whether the group is enabled")
	updateDeviceGroupCmd.Flags().BoolVar(&disabled, "disable", false, "Whether the group is disabled")
	updateDeviceGroupCmd.Flags().BoolVar(&enabledDefaultConnectivity, "enableDefaultConnectivity", false, "Default connectivity for devices joining the device group (enabled)")
	updateDeviceGroupCmd.Flags().BoolVar(&disabledDefaultConnectivity, "disableDefaultConnectivity", false, "Default connectivity for devices joining the device group (disabled)")
	deviceGroupCmd.AddCommand(updateDeviceGroupCmd)

	removeDeviceGroupCmd.Flags().StringVar(&deviceGroupID, "deviceGroupId", "", "Device group identifier")
	deviceGroupCmd.AddCommand(removeDeviceGroupCmd)

	deviceGroupCmd.AddCommand(listDeviceGroupsCmd)

	// Devices
	rootCmd.AddCommand(devicesCmd)
	devicesCmd.PersistentFlags().StringVar(&deviceGroupID, "deviceGroupId", "", "Device group identifier")

	devicesCmd.AddCommand(listDevicesCmd)

	devicesCmd.AddCommand(deviceInfoCmd)

	devicesCmd.AddCommand(deviceLabelsCmd)
	deviceLabelsCmd.PersistentFlags().StringVar(&deviceID, "deviceId", "", "Device identifier")
	deviceLabelsCmd.PersistentFlags().StringVar(&rawLabels, "labels", "", "Labels separated by ; as in key1:value;key2:value")

	deviceLabelsCmd.AddCommand(addLabelToDeviceCmd)
	deviceLabelsCmd.AddCommand(removeLabelFromDeviceCmd)

	updateDeviceCmd.Flags().StringVar(&deviceID, "deviceId", "", "Device identifier")
	updateDeviceCmd.Flags().BoolVar(&enabled, "enabled", false, "Whether the device is enabled")
	updateDeviceCmd.Flags().BoolVar(&disabled, "disabled", false, "Whether the device is disabled")
	devicesCmd.AddCommand(updateDeviceCmd)

	devicesCmd.AddCommand(removeDeviceCmd)
	removeDeviceCmd.PersistentFlags().StringVar(&deviceID, "deviceId", "", "Device identifier")

}

var addDeviceGroupCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new device group",
	Long:  `Add a new device group to an organization`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"name"}, args, []string{name})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			n.AddDeviceGroup(options.Resolve("organizationID", organizationID),
				targetValues[0], enabled, disabled, enabledDefaultConnectivity, disabledDefaultConnectivity)
		}
	},
}

var updateDeviceGroupCmd = &cobra.Command{
	Use:   "update [deviceGroupID]",
	Short: "Update a device group",
	Long:  `Update the options of a device group`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"deviceGroupID"}, args, []string{deviceGroupID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			n.UpdateDeviceGroup(options.Resolve("organizationID", organizationID),
				targetValues[0], enabled, disabled, enabledDefaultConnectivity, disabledDefaultConnectivity)
		}
	},
}

var removeDeviceGroupCmd = &cobra.Command{
	Use:   "delete [deviceGroupID]",
	Aliases: []string{"remove", "del"},
	Short: "Remove a device group",
	Long:  `Remove a device group`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"deviceGroupID"}, args, []string{deviceGroupID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			n.RemoveDeviceGroup(options.Resolve("organizationID", organizationID),
				targetValues[0])
		}
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
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		n.ListDeviceGroups(options.Resolve("organizationID", organizationID))
	},
}

var listDevicesCmd = &cobra.Command{
	Use:   "list",
	Short: "List the devices in a device group",
	Long:  `List the devices in a device group`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"deviceGroupID"}, args, []string{deviceGroupID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			n.ListDevices(options.Resolve("organizationID", organizationID),
				targetValues[0])
		}

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
	Use:   "add [deviceGroupID] [deviceID] [labels]",
	Short: "Add a set of labels to a device",
	Long:  `Add a set of labels to a device`,
	Args: cobra.MaximumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"deviceGroupID", "deviceID", "labels"}, args, []string{deviceGroupID, deviceID, rawLabels})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			n.AddLabelToDevice(options.Resolve("organizationID", organizationID),
				targetValues[0], targetValues[1], targetValues[2])
		}
	},
}

var removeLabelFromDeviceCmd = &cobra.Command{
	Use:   "delete [deviceGroupID] [deviceID] [labels]",
	Aliases: []string{"remove", "del"},
	Short: "Remove a set of labels from a device",
	Long:  `Remove a set of labels from a device`,
	Args: cobra.MaximumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"deviceGroupID", "deviceID", "labels"}, args, []string{deviceGroupID, deviceID, rawLabels})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			n.RemoveLabelFromDevice(options.Resolve("organizationID", organizationID),
				targetValues[0], targetValues[1], targetValues[2])
		}
	},
}

var updateDeviceCmd = &cobra.Command{
	Use:   "update [deviceGroupID] [deviceID]",
	Short: "Update the information of a device",
	Long:  `Update the information of a device`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"deviceGroupID", "deviceID"}, args, []string{deviceGroupID, deviceID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			n.UpdateDevice(options.Resolve("organizationID", organizationID),
				targetValues[0], targetValues[1], enabled, disabled)
		}
	},
}

var removeDeviceCmd = &cobra.Command{
	Use:   "delete [deviceGroupID] [deviceID]",
	Aliases: []string{"remove", "del"},
	Short: "Remove a device",
	Long:  `Remove a device`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"deviceGroupID", "deviceID"}, args, []string{deviceGroupID, deviceID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			n.RemoveDevice(options.Resolve("organizationID", organizationID),
				targetValues[0], targetValues[1])
		}
	},
}

var deviceInfoCmd = &cobra.Command{
	Use: "info [deviceGroupID] [deviceID]",
	Short: "Show device info",
	Long:  `Show device info`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"deviceGroupID", "deviceID"}, args, []string{deviceGroupID, deviceID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			n.GetDeviceInfo(options.Resolve("organizationID", organizationID), targetValues[0], targetValues[1])
		}
	},
}