/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var organizationID string

var orgCmd = &cobra.Command{
	Use:   "org",
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
	Short: "Retrieve organization information",
	Long:  `Retrieve organization information`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		o := cli.NewOrganizations(nalejAddress, nalejPort)
		o.Info(options.Resolve("organizationID", organizationID))
	},
}