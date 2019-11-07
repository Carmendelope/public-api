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
 *
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
	Use:     "list [clusterID]",
	Aliases: []string{"ls"},
	Short:   "List nodes",
	Long:    `List nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewNodes(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"clusterID"}, args, []string{clusterID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			n.List(cliOptions.Resolve("organizationID", organizationID),
				cliOptions.Resolve("clusterID", targetValues[0]))
		}
	},
}

var nodeLabelsCmd = &cobra.Command{
	Use:     "label",
	Aliases: []string{"labels", "l"},
	Short:   "Manage node labels",
	Long:    `Manage node labels`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var addLabelToNodeCmd = &cobra.Command{
	Use:   "add [nodeID] [labels]",
	Short: "Add a set of labels to a node",
	Long:  `Add a set of labels to a node`,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewNodes(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"nodeID", "labels"}, args, []string{nodeID, rawLabels})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			n.ModifyNodeLabels(cliOptions.Resolve("organizationID", organizationID),
				targetValues[0], true, targetValues[1])
		}
	},
}

var removeLabelFromNodeCmd = &cobra.Command{
	Use:     "delete [nodeID] [labels]",
	Aliases: []string{"remove", "del", "rm"},
	Short:   "Remove a set of labels from a cluster",
	Long:    `Remove a set of labels from a cluster`,
	Args:    cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		n := cli.NewNodes(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"nodeID", "labels"}, args, []string{nodeID, rawLabels})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			n.ModifyNodeLabels(cliOptions.Resolve("organizationID", organizationID),
				targetValues[0], false, targetValues[1])
		}
	},
}
