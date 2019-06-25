/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
	"strings"
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

var (
	queryAssetSelector = &cli.AssetSelector{}
	queryTimeRange = &cli.TimeRange{}
	queryMetrics = []string{}
	queryAggr = ""
	listAssetSelector = &cli.AssetSelector{}
)

func addAssetSelector(cmd *cobra.Command, selector *cli.AssetSelector) {
	cmd.Flags().StringVar(&selector.EdgeControllerId, "edgeControllerId", "", "Select assets for this Edge Controller")
	cmd.Flags().StringSliceVar(&selector.AssetIds, "assetId", []string{}, "Select one or multiple assets")
	cmd.Flags().StringSliceVar(&selector.GroupIds, "groupId", []string{}, "Select assets in one or multiple groups")
	cmd.Flags().StringToStringVar(&selector.Labels, "label", map[string]string{}, "Select assets with the intersection of a set of labels")
}

func addTimeRange(cmd *cobra.Command, timeRange *cli.TimeRange) {
	cmd.Flags().StringVar(&timeRange.Timestamp, "timestamp", "", "Timestamp for point query")
	cmd.Flags().StringVar(&timeRange.Start, "start", "", "Start time for range query")
	cmd.Flags().StringVar(&timeRange.End, "end", "", "End time for range query")
	cmd.Flags().DurationVar(&timeRange.Resolution, "resolution", 0, "Range interval resolution - 0 to aggregate to single value")
}

func init() {
	rootCmd.AddCommand(inventoryCmd)
	inventoryCmd.AddCommand(inventoryListCmd)
	inventoryCmd.AddCommand(invControllerCommand)
	inventoryCmd.AddCommand(invAssetCommand)
	inventoryCmd.AddCommand(invDeviceCommand)
	inventoryCmd.AddCommand(invMonitoringCmd)

	invControllerCommand.AddCommand(invControllerExtInfoCmd)
	invControllerCommand.AddCommand(invEdgeControllerUpdateLocationCmd)
	invControllerCommand.AddCommand(edgeControllerLabelsCmd)
	edgeControllerLabelsCmd.AddCommand(addLabelToECCmd)
	edgeControllerLabelsCmd.AddCommand(removeLabelFromECCmd)

	invAssetCommand.AddCommand(invAssetInfoCmd)
	invAssetCommand.AddCommand(invAssetUpdateLocationCmd)
	invAssetCommand.AddCommand(assetLabelsCmd)
	assetLabelsCmd.AddCommand(addLabelToAssetCmd)
	assetLabelsCmd.AddCommand(removeLabelFromAssetCmd)

	invDeviceCommand.AddCommand(invDeviceInfoCmd)
	invDeviceCommand.AddCommand(invDeviceUpdateLocationCmd)
	invDeviceCommand.AddCommand(invDeviceLabelsCmd)
	invDeviceLabelsCmd.AddCommand(addLabelToInvDeviceCmd)
	invDeviceLabelsCmd.AddCommand(removeLabelFromInvDeviceCmd)

	invMonitoringCmd.AddCommand(invMonitoringListCmd)

	addAssetSelector(invMonitoringCmd, queryAssetSelector)
	addAssetSelector(invMonitoringListCmd, listAssetSelector)
	addTimeRange(invMonitoringCmd, queryTimeRange)

	invMonitoringCmd.Flags().StringSliceVar(&queryMetrics, "metric", []string{}, "Metrics to query; all if empty")
	invMonitoringCmd.Flags().StringVar(&queryAggr, "aggregation", "NONE", "Aggregation type")

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

var invDeviceLabelsCmd = &cobra.Command{
	Use:   "label",
	Aliases: []string{"labels", "l"},
	Short: "Manage device labels",
	Long:  `Manage device labels`,
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
			insecure,
			useTLS,
			options.Resolve("cacert", caCertPath),
			options.Resolve("output", output))
		ec.GetAssetInfo(options.Resolve("organizationID", organizationID), args[0])
	},
}

var invAssetUpdateLocationCmd = &cobra.Command{
	Use:   "location-update [assetID] [location]",
	Short: "Update asset location",
	Long:  `Update asset location`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewAsset(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			useTLS,
			options.Resolve("cacert", caCertPath),
			options.Resolve("output", output))

			assetID = args[0]
			newLocation := args[1]
		a.UpdateLocation(options.Resolve("organizationID", organizationID), assetID, newLocation)
	},
}

var assetLabelsCmd = &cobra.Command{
	Use:   "label",
	Aliases: []string{"labels", "l"},
	Short: "Manage asset labels",
	Long:  `Manage asset labels`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
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
			insecure,
			useTLS,
			options.Resolve("cacert", caCertPath),
			options.Resolve("output", output))
		ec.GetDeviceInfo(options.Resolve("organizationID", organizationID), args[0])
	},
}

var invDeviceUpdateLocationCmd = &cobra.Command{
	Use:   "location-update [assetDeviceID] [location]",
	Short: "update the location of a device",
	Long:  `Update the location of a device`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		device := cli.NewInventory(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			useTLS,
			options.Resolve("cacert", caCertPath),
			options.Resolve("output", output))

			assetDeviceId = args[0]
			newLocation := args[1]
		device.UpdateDeviceLocation(options.Resolve("organizationID", organizationID), assetDeviceId, newLocation)
	},
}

var invEdgeControllerUpdateLocationCmd = &cobra.Command{
	Use:   "location-update [edgeControllerID] [location]",
	Short: "update the location of a device",
	Long:  `Update the location of a device`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		ec := cli.NewEdgeController(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

			edgeControllerID = args[0]
			newLocation := args[1]

		ec.UpdateGeolocation(options.Resolve("organizationID", organizationID), edgeControllerID, newLocation)
	},
}

var invMonitoringCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Asset metrics retrieval",
	Long:  `Metrics for an asset or aggregated metrics for a group of assets`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		inv := cli.NewInventoryMonitoring(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		queryAssetSelector.OrganizationId = options.Resolve("organizationID", organizationID)
		inv.QueryMetrics(queryAssetSelector, queryMetrics, queryTimeRange, queryAggr)
	},
}

var invMonitoringListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available metrics for assets",
	Long:  `List available metrics for assets`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		inv := cli.NewInventoryMonitoring(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		listAssetSelector.OrganizationId = options.Resolve("organizationID", organizationID)
		inv.ListMetrics(listAssetSelector)
	},
}

var edgeControllerLabelsCmd = &cobra.Command{
	Use:   "label",
	Aliases: []string{"labels", "l"},
	Short: "Manage ec labels",
	Long:  `Manage ec labels`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var addLabelToAssetCmd = &cobra.Command{
	Use:   "add [assetID] [labels]",
	Short: "Add a set of labels to an asset",
	Long:  `Add a set of labels to an asset`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewAsset(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		n.AddLabelToAsset(options.Resolve("organizationID", organizationID), args[0], args[1])
	},
}

var removeLabelFromAssetCmd = &cobra.Command{
	Use:   "delete [assetID] [labels]",
	Aliases: []string{"remove", "del"},
	Short: "Remove a set of labels from an asset",
	Long:  `Remove a set of labels from an asset`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewAsset(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		n.RemoveLabelFromAsset(options.Resolve("organizationID", organizationID), args[0], args[1])
	},
}

var addLabelToECCmd = &cobra.Command{
	Use:   "add [edgeControllerID] [labels]",
	Short: "Add a set of labels to an EC",
	Long:  `Add a set of labels to an EC`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewEdgeController(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		n.AddLabelToEC(options.Resolve("organizationID", organizationID), args[0], args[1])
	},
}

var removeLabelFromECCmd = &cobra.Command{
	Use:   "delete [edgeControllerID] [labels]",
	Aliases: []string{"remove", "del"},
	Short: "Remove a set of labels from an EC",
	Long:  `Remove a set of labels from an EC`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewEdgeController(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		n.RemoveLabelFromEC(options.Resolve("organizationID", organizationID),	args[0], args[1])
	},
}

var addLabelToInvDeviceCmd = &cobra.Command{
	Use:   "add [assetDeviceId] [labels]",
	Short: "Add a set of labels to a device",
	Long:  `Add a set of labels to a device`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		device := strings.Split(args[0],"#")
		deviceGroupID := device[0]
		deviceID := device[1]

		n.AddLabelToDevice(options.Resolve("organizationID", organizationID), deviceGroupID, deviceID, args[1])
	},
}

var removeLabelFromInvDeviceCmd = &cobra.Command{
	Use:   "delete [assetDeviceId] [labels]",
	Aliases: []string{"remove", "del"},
	Short: "Remove a set of labels from a device",
	Long:  `Remove a set of labels from a device`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewDevices(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		device := strings.Split(args[0],"#")
		deviceGroupID := device[0]
		deviceID := device[1]

		n.RemoveLabelFromDevice(options.Resolve("organizationID", organizationID),	deviceGroupID, deviceID, args[1])
	},
}