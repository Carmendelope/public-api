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
var description string

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
	appsCmd.PersistentFlags().StringVar(&descriptorID, "descriptorID", "", "Application descriptor identifier")
	appsCmd.AddCommand(descriptorCmd)
	addDescriptorCmd.PersistentFlags().StringVar(&descriptorPath, "descriptorPath", "", "Application descriptor path containing a JSON spec")
	descriptorCmd.AddCommand(addDescriptorCmd)
	descriptorCmd.AddCommand(getDescriptorCmd)
	descriptorCmd.AddCommand(listDescriptorsCmd)
	descriptorCmd.AddCommand(addDescriptorHelpCmd)
	appsCmd.AddCommand(instanceCmd)
	deployInstanceCmd.Flags().StringVar(&name, "name", "", "Name of the application instance")
	deployInstanceCmd.Flags().StringVar(&description, "description", "", "Description of the application instance")
	instanceCmd.AddCommand(deployInstanceCmd)
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
		a := cli.NewApplications(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		a.AddDescriptorHelp()
	},
}

var listDescriptorsCmd = &cobra.Command{
	Use:   "list",
	Short: "List the application descriptors",
	Long:  `List the application descriptors`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		a.ListDescriptors(options.Resolve("organizationID", organizationID))
	},
}

var getDescriptorCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an application descriptor",
	Long:  `Get an application descriptor`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		a.GetDescriptor(options.Resolve("organizationID", organizationID), descriptorID)
	},
}

var instanceCmd = &cobra.Command{
	Use:   "inst",
	Aliases: []string{"instance"},
	Short: "Manage applications instances",
	Long:  `Manage applications instances`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var deployInstanceCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an application instance",
	Long:  `Deploy an application instance`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		a.Deploy(options.Resolve("organizationID", organizationID), descriptorID, name, description)
	},
}

var listInstancesCmd = &cobra.Command{
	Use:   "list",
	Short: "List application instances",
	Long:  `List application intances`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		a.ListInstances(options.Resolve("organizationID", organizationID))
	},
}
