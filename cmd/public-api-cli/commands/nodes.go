/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"fmt"
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var nodesCmd = &cobra.Command{
	Use:     "node",
	Aliases: []string{"nodes"},
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
	Use:   "list [clusterID]",
	Aliases: []string{"ls"},
	Short: "List nodes",
	Long:  `List nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewNodes(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output), options.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"clusterID"}, args, []string{clusterID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			n.List(options.Resolve("organizationID", organizationID),
				options.Resolve("clusterID", targetValues[0]))
		}
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
	Use:   "add [nodeID] [labels]",
	Short: "Add a set of labels to a node",
	Long:  `Add a set of labels to a node`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewNodes(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output), options.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"nodeID", "labels"}, args, []string{nodeID, rawLabels})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			n.ModifyNodeLabels(options.Resolve("organizationID", organizationID),
				targetValues[0], true, targetValues[1])
		}
	},
}

var removeLabelFromNodeCmd = &cobra.Command{
	Use:   "delete [nodeID] [labels]",
	Aliases: []string{"remove", "del", "rm"},
	Short: "Remove a set of labels from a cluster",
	Long:  `Remove a set of labels from a cluster`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewNodes(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output), options.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"nodeID", "labels"}, args, []string{nodeID, rawLabels})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			n.ModifyNodeLabels(options.Resolve("organizationID", organizationID),
				targetValues[0], false, targetValues[1])
		}
	},
}