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
	invAssetCommand.AddCommand(invAssetInfoCmd)
	invDeviceCommand.AddCommand(invDeviceInfoCmd)
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

var invControllerExtInfoCmd = &cobra.Command{
	Use:   "info [edgeControllerID]",
	Short: "Get extended information on an edge controller",
	Long:  `Get extended information on an edge controller`,
	Args: cobra.ExactArgs(1),
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
	Args: cobra.ExactArgs(1),
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
	Args: cobra.ExactArgs(1),
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
