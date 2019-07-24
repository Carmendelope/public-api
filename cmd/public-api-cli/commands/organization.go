/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var orgCmd = &cobra.Command{
	Use:   "organization",
	Aliases: []string{"org"},
	Short: "Organization related operations",
	Long:  `Organization related operations`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	orgCmd.PersistentFlags().StringVar(&organizationID, "organizationID", "", "Organization identifier")
	rootCmd.AddCommand(orgCmd)
	orgCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Aliases: []string{"get"},
	Short: "Retrieve organization information",
	Long:  `Retrieve organization information`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		o := cli.NewOrganizations(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath),
			options.Resolve("output", output),
			options.ResolveAsInt("labelLength", labelLength))
		o.Info(options.Resolve("organizationID", organizationID))
	},
}
