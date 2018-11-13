/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var clustersCmd = &cobra.Command{
	Use:   "cluster",
	Aliases: []string{"clusters"},
	Short: "Manage clusters",
	Long:  `Manage clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(clustersCmd)
	clustersCmd.AddCommand(listClustersCmd)
}

var listClustersCmd = &cobra.Command{
	Use:   "list",
	Short: "List clusters",
	Long:  `List clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		c.List(options.Resolve("organizationID", organizationID))
	},
}
