/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var descriptorID string
var descriptorPath string

var appsCmd = &cobra.Command{
	Use:   "app",
	Aliases: []string{"applications"},
	Short: "Manage applications",
	Long:  `Manage applications`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)
	appsCmd.PersistentFlags().StringVar(&organizationID, "organizationID", "", "Organization identifier")
	appsCmd.AddCommand(descriptorCmd)
	descriptorCmd.PersistentFlags().StringVar(&descriptorID, "descriptorID", "", "Application descriptor identifier")
	addDescriptorCmd.PersistentFlags().StringVar(&descriptorPath, "descriptorPath", "", "Application descriptor path containing a JSON spec")
	descriptorCmd.AddCommand(addDescriptorCmd)
	descriptorCmd.AddCommand(getDescriptorCmd)
	descriptorCmd.AddCommand(listDescriptorsCmd)
	descriptorCmd.AddCommand(addDescriptorHelpCmd)
}

var descriptorCmd = &cobra.Command{
	Use:   "desc",
	Aliases: []string{"descriptor"},
	Short: "Manage applications descriptors",
	Long:  `Manage applications descriptors`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var addDescriptorCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new application descriptor",
	Long:  `Add a new application descriptor`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		a.AddDescriptor(options.Resolve("organizationID", organizationID), descriptorPath)
	},
}

var addDescriptorHelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help related to adding a new application descriptor",
	Long:  `Show help related to adding a new application descriptor`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(nalejAddress, nalejPort)
		a.AddDescriptorHelp()
	},
}

var listDescriptorsCmd = &cobra.Command{
	Use:   "list",
	Short: "List the application descriptors",
	Long:  `List the application descriptors`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(nalejAddress, nalejPort)
		a.ListDescriptors(options.Resolve("organizationID", organizationID))
	},
}

var getDescriptorCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an application descriptor",
	Long:  `Get an application descriptor`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(nalejAddress, nalejPort)
		a.GetDescriptor(options.Resolve("organizationID", organizationID), descriptorID)
	},
}
