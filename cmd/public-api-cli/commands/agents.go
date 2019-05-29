/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package commands

import (
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

}

var createAgentJoinTokenCmd = &cobra.Command{
	Use:   "create-join-token [edgeControllerId]",
	Short: "Create a join token to attach new agent to an edge controller",
	Long:  `Create a join token for being able to attach new agent to an edge controller`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		agent := cli.NewAgent(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
			agent.CreateAgentJoinToken(options.Resolve("organizationID", organizationID),
				                       args[0],
				                       outputPath)
	},
}