/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
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
	agentCmd.AddCommand(createAgentJoinTokenCmd)
	createAgentJoinTokenCmd.Flags().StringVar(&edgeControllerID, "edgeControllerID", "", "edge controller id to attach the agent")

}

var createAgentJoinTokenCmd = &cobra.Command{
	Use:   "create-join-token to attach new agent to an edge controller",
	Short: "Create a join token",
	Long:  `Create a join token for being able to attach new agent to an edge controller`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		agent := cli.NewAgent(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"edgeControllerID"}, args, []string{edgeControllerID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			agent.CreateAgentJoinToken(options.Resolve("organizationID", organizationID),
				                       options.Resolve("edgeControllerID", targetValues[0]),
				                       outputPath)

		}
	},
}