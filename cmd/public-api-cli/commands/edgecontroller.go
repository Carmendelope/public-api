/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var edgeControllerCmd = &cobra.Command{
	Use:     "edgecontroller",
	Aliases: []string{"ec"},
	Short:   "Manage edge controllers",
	Long:    `Manage edge controllers`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	// Device groups
	rootCmd.AddCommand(edgeControllerCmd)
	edgeControllerCmd.AddCommand(createJoinTokenECCmd)
	createJoinTokenECCmd.Flags().StringVar(&outputPath, "outputPath", "", "Path to store the resulting token")
}

var createJoinTokenECCmd = &cobra.Command{
	Use:   "create-join-token",
	Short: "Create a join token",
	Long:  `Create a join token for being able to attach new edge controllers to the platform`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		ec := cli.NewEdgeController(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
			ec.CreateJoinToken(options.Resolve("organizationID", organizationID), outputPath)
	},
}

