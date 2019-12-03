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

package commands

import (
	"fmt"
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var deviceGroupCmd = &cobra.Command{
	Use:     "devicegroup",
	Aliases: []string{"devicegroups", "dg"},
	Short:   "Manage device groups",
	Long:    `Manage device groups`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var devicesCmd = &cobra.Command{
	Use:     "device",
	Aliases: []string{"devices", "dev"},
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

	updateDeviceGroupCmd.Flags().BoolVar(&enabled, "enable", false, "Whether the group is enabled")
	updateDeviceGroupCmd.Flags().BoolVar(&disabled, "disable", false, "Whether the group is disabled")
	updateDeviceGroupCmd.Flags().BoolVar(&enabledDefaultConnectivity, "enableDefaultConnectivity", false, "Default connectivity for devices joining the device group (enabled)")
	updateDeviceGroupCmd.Flags().BoolVar(&disabledDefaultConnectivity, "disableDefaultConnectivity", false, "Default connectivity for devices joining the device group (disabled)")
	deviceGroupCmd.AddCommand(updateDeviceGroupCmd)

	removeDeviceGroupCmd.Flags().StringVar(&deviceGroupID, "deviceGroupID", "", "Device group identifier")
	deviceGroupCmd.AddCommand(removeDeviceGroupCmd)

	deviceGroupCmd.AddCommand(listDeviceGroupsCmd)

	// Devices
	rootCmd.AddCommand(devicesCmd)
	devicesCmd.PersistentFlags().StringVar(&deviceGroupID, "deviceGroupID", "", "Device group identifier")

	devicesCmd.AddCommand(listDevicesCmd)

	devicesCmd.AddCommand(deviceInfoCmd)

	devicesCmd.AddCommand(deviceLabelsCmd)
	deviceLabelsCmd.PersistentFlags().StringVar(&deviceID, "deviceID", "", "Device identifier")
	deviceLabelsCmd.PersistentFlags().StringVar(&rawLabels, "labels", "", "Labels separated by ; as in key1:value;key2:value")

	deviceLabelsCmd.AddCommand(addLabelToDeviceCmd)
	deviceLabelsCmd.AddCommand(removeLabelFromDeviceCmd)

	updateDeviceCmd.Flags().StringVar(&deviceID, "deviceID", "", "Device identifier")
	updateDeviceCmd.Flags().BoolVar(&enabled, "enabled", false, "Whether the device is enabled")
	updateDeviceCmd.Flags().BoolVar(&disabled, "disabled", false, "Whether the device is disabled")
	devicesCmd.AddCommand(updateDeviceCmd)

	devicesCmd.AddCommand(removeDeviceCmd)
	removeDeviceCmd.PersistentFlags().StringVar(&deviceID, "deviceID", "", "Device identifier")

}

var addDeviceGroupCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new device group",
	Long:  `Add a new device group to an organization`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"name"}, args, []string{name})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			n.AddDeviceGroup(cliOptions.Resolve("organizationID", organizationID),
				targetValues[0], enabled, disabled, enabledDefaultConnectivity, disabledDefaultConnectivity)
		}
	},
}

var updateDeviceGroupCmd = &cobra.Command{
	Use:   "update [deviceGroupID]",
	Short: "Update a device group",
	Long:  `Update the options of a device group`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"deviceGroupID"}, args, []string{deviceGroupID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			n.UpdateDeviceGroup(cliOptions.Resolve("organizationID", organizationID),
				targetValues[0], enabled, disabled, enabledDefaultConnectivity, disabledDefaultConnectivity)
		}
	},
}

var removeDeviceGroupCmd = &cobra.Command{
	Use:     "delete [deviceGroupID]",
	Aliases: []string{"remove", "del", "rm"},
	Short:   "Remove a device group",
	Long:    `Remove a device group`,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"deviceGroupID"}, args, []string{deviceGroupID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			n.RemoveDeviceGroup(cliOptions.Resolve("organizationID", organizationID),
				targetValues[0])
		}
	},
}

var listDeviceGroupsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the device groups in an organization",
	Long:    `List the device groups in an organization`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		n.ListDeviceGroups(cliOptions.Resolve("organizationID", organizationID))
	},
}

var listDevicesCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the devices in a device group",
	Long:    `List the devices in a device group`,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"deviceGroupID"}, args, []string{deviceGroupID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			n.ListDevices(cliOptions.Resolve("organizationID", organizationID),
				targetValues[0])
		}

	},
}

var deviceLabelsCmd = &cobra.Command{
	Use:     "label",
	Aliases: []string{"labels", "l"},
	Short:   "Manage device labels",
	Long:    `Manage device labels`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var addLabelToDeviceCmd = &cobra.Command{
	Use:   "add [deviceGroupID] [deviceID] [labels]",
	Short: "Add a set of labels to a device",
	Long:  `Add a set of labels to a device`,
	Args:  cobra.MaximumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"deviceGroupID", "deviceID", "labels"}, args, []string{deviceGroupID, deviceID, rawLabels})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			n.AddLabelToDevice(cliOptions.Resolve("organizationID", organizationID),
				targetValues[0], targetValues[1], targetValues[2])
		}
	},
}

var removeLabelFromDeviceCmd = &cobra.Command{
	Use:     "delete [deviceGroupID] [deviceID] [labels]",
	Aliases: []string{"remove", "del", "rm"},
	Short:   "Remove a set of labels from a device",
	Long:    `Remove a set of labels from a device`,
	Args:    cobra.MaximumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"deviceGroupID", "deviceID", "labels"}, args, []string{deviceGroupID, deviceID, rawLabels})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			n.RemoveLabelFromDevice(cliOptions.Resolve("organizationID", organizationID),
				targetValues[0], targetValues[1], targetValues[2])
		}
	},
}

var updateDeviceCmd = &cobra.Command{
	Use:   "update [deviceGroupID] [deviceID]",
	Short: "Update the information of a device",
	Long:  `Update the information of a device`,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"deviceGroupID", "deviceID"}, args, []string{deviceGroupID, deviceID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			n.UpdateDevice(cliOptions.Resolve("organizationID", organizationID),
				targetValues[0], targetValues[1], enabled, disabled)
		}
	},
}

var removeDeviceCmd = &cobra.Command{
	Use:     "delete [deviceGroupID] [deviceID]",
	Aliases: []string{"remove", "del", "rm"},
	Short:   "Remove a device",
	Long:    `Remove a device`,
	Args:    cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"deviceGroupID", "deviceID"}, args, []string{deviceGroupID, deviceID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			n.RemoveDevice(cliOptions.Resolve("organizationID", organizationID),
				targetValues[0], targetValues[1])
		}
	},
}

var deviceInfoCmd = &cobra.Command{
	Use:     "info [deviceGroupID] [deviceID]",
	Aliases: []string{"get"},
	Short:   "Show device info",
	Long:    `Show device info`,
	Args:    cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure,
			useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"deviceGroupID", "deviceID"}, args, []string{deviceGroupID, deviceID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			n.GetDeviceInfo(cliOptions.Resolve("organizationID", organizationID), targetValues[0], targetValues[1])
		}
	},
}
