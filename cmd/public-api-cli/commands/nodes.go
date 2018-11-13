/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var nodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: []string{"node"},
	Short:   "Manage nodes",
	Long:    `Manage nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(nodesCmd)
	listNodesCmd.Flags().StringVar(&clusterID, "clusterID", "", "Cluster identifier")
	nodesCmd.AddCommand(listNodesCmd)
}

var listNodesCmd = &cobra.Command{
	Use:   "list",
	Short: "List nodes",
	Long:  `List nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewNodes(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		n.List(options.Resolve("organizationID", organizationID),
			options.Resolve("clusterID", clusterID))
	},
}
