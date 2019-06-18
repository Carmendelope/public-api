/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var inventoryCmd = &cobra.Command{
	Use:     "inventory",
	Aliases: []string{"inv"},
	Short:   "Manage inventory",
	Long:    `Manage inventory`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(inventoryCmd)
	inventoryCmd.AddCommand(inventoryListCmd)
	inventoryCmd.AddCommand(invControllerCommand)
	inventoryCmd.AddCommand(invAssetCommand)
	inventoryCmd.AddCommand(invDeviceCommand)

	invControllerCommand.AddCommand(invControllerExtInfoCmd)
	invAssetCommand.AddCommand(invAssetInfoCmd)
	invDeviceCommand.AddCommand(invDeviceInfoCmd)
	invDeviceCommand.AddCommand(invDeviceUpdateLocationCmd)
}

var inventoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the inventory",
	Long:  `List the inventory in a given organization`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		ec := cli.NewInventory(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		ec.List(options.Resolve("organizationID", organizationID))
	},
}

var invControllerCommand = &cobra.Command{
	Use:     "controller",
	Aliases: []string{"ec"},
	Short:   "Controller commands",
	Long:    `Controller commands`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var invAssetCommand = &cobra.Command{
	Use:   "asset",
	Short: "Asset commands",
	Long:  `Asset commands`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var invDeviceCommand = &cobra.Command{
	Use:     "device",
	Aliases: []string{"dev"},
	Short:   "Device commands",
	Long:    `Device commands`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var invControllerExtInfoCmd = &cobra.Command{
	Use:   "info [edgeControllerID]",
	Short: "Get extended information on an edge controller",
	Long:  `Get extended information on an edge controller`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		ec := cli.NewInventory(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		ec.GetControllerExtendedInfo(options.Resolve("organizationID", organizationID), args[0])
	},
}

var invAssetInfoCmd = &cobra.Command{
	Use:   "info [assetID]",
	Short: "Get extended information on an asset",
	Long:  `Get extended information on an asset`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		ec := cli.NewInventory(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		ec.GetAssetInfo(options.Resolve("organizationID", organizationID), args[0])
	},
}

var invDeviceInfoCmd = &cobra.Command{
	Use:   "info [deviceID]",
	Short: "Get extended information of a device",
	Long:  `Get extended information of a device`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		ec := cli.NewInventory(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		ec.GetDeviceInfo(options.Resolve("organizationID", organizationID), args[0])
	},
}

var invAssetUpdateLocationCmd = &cobra.Command{
	Use:   "location-update [assetID]",
	Short: "Update asset location",
	Long:  `Update asset location`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewAsset(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		if len(args) > 0 {
			assetID = args[0]
		}
		a.UpdateLocation(options.Resolve("organizationID", organizationID), assetID, assetLocation)
	},
}

var invDeviceUpdateLocationCmd = &cobra.Command{
	Use:   "location-update [assetDeviceID] ",
	Short: "update the location of a device",
	Long:  `Update the location of a device`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		device := cli.NewInventory(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		device.UpdateDeviceLocation(options.Resolve("organizationID", organizationID), args[0], deviceLocation)
	},
}
