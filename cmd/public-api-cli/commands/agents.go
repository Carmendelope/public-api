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

var agentCmd = &cobra.Command{
	Use:     "agent",
	Aliases: []string{"ag"},
	Short:   "Manage agents",
	Long:    `Manage agents`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)

	// CreateAgentJoinToken
	createAgentJoinTokenCmd.Flags().StringVar(&outputPath, "outputPath", "", "Path to store the resulting token")
	agentCmd.AddCommand(createAgentJoinTokenCmd)

	// ActivateAgentMonitoring
	activateAgentMontoringCmd.Flags().BoolVar(&activate, "activate", true, "Activate/Deactivate monitoring")
	agentCmd.AddCommand(activateAgentMontoringCmd)

	// UninstallAgentCmd
	uninstallAgentCmd.Flags().BoolVar(&force, "force", false, "force the agent uninstall")
	agentCmd.AddCommand(uninstallAgentCmd)

}

var createAgentJoinTokenCmd = &cobra.Command{
	Use:   "create-join-token [edgeControllerID]",
	Short: "Create a join token to attach new agent to an edge controller",
	Long:  `Create a join token for being able to attach new agent to an edge controller`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		agent := cli.NewAgent(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure,
			useTLS,
			cliOptions.Resolve("cacert", caCertPath),
			cliOptions.Resolve("output", output),
			cliOptions.ResolveAsInt("labelLength", labelLength))
		agent.CreateAgentJoinToken(cliOptions.Resolve("organizationID", organizationID),
			args[0],
			outputPath)
	},
}

var activateAgentMontoringCmd = &cobra.Command{
	Use:     "monitoring [edgeControllerID] [assetID]",
	Aliases: []string{"mon"},
	Short:   "Activate agent monitoring",
	Long:    `Activate agent monitoring`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		agent := cli.NewAgent(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"edgeControllerId", "assetID"}, args, []string{edgeControllerID, assetID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			agent.ActivateAgentMonitoring(cliOptions.Resolve("organizationID", organizationID),
				targetValues[0], targetValues[1], activate)
		}
	},
}

var uninstallAgentCmd = &cobra.Command{
	Use:   "uninstall [assetID]",
	Short: "Uninstall agent",
	Long:  `Uninstall agent from edge-controller`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		agent := cli.NewAgent(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath),
			cliOptions.Resolve("output", output),
			cliOptions.ResolveAsInt("labelLength", labelLength))
		agent.UninstallAgent(cliOptions.Resolve("organizationID", organizationID), args[0], force)
	},
}
