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

	nodeLabelsCmd.PersistentFlags().StringVar(&nodeID, "nodeID", "", "Node identifier")
	nodeLabelsCmd.PersistentFlags().StringVar(&rawLabels, "labels", "", "Labels separated by ; as in key1:value;key2:value")
	nodesCmd.AddCommand(nodeLabelsCmd)
	nodeLabelsCmd.AddCommand(addLabelToNodeCmd)
	nodeLabelsCmd.AddCommand(removeLabelFromNodeCmd)
}

var listNodesCmd = &cobra.Command{
	Use:   "list",
	Short: "List nodes",
	Long:  `List nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewNodes(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		n.List(options.Resolve("organizationID", organizationID),
			options.Resolve("clusterID", clusterID))
	},
}

var nodeLabelsCmd = &cobra.Command{
	Use:   "label",
	Aliases: []string{"labels", "l"},
	Short: "Manage node labels",
	Long:  `Manage node labels`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var addLabelToNodeCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a set of labels to a node",
	Long:  `Add a set of labels to a node`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewNodes(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		n.ModifyNodeLabels(options.Resolve("organizationID", organizationID),
			nodeID, true, rawLabels)
	},
}

var removeLabelFromNodeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a set of labels from a cluster",
	Long:  `Remove a set of labels from a cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewNodes(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		n.ModifyNodeLabels(options.Resolve("organizationID", organizationID),
			nodeID, false, rawLabels)
	},
}