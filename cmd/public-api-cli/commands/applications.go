/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// TODO Remove descriptor NP-338

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:     "app",
	Aliases: []string{"applications"},
	Short:   "Manage applications",
	Long:    `Manage applications`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)
	appsCmd.PersistentFlags().StringVar(&organizationID, "organizationID", "", "Organization identifier")
	// Descriptor commands
	appsCmd.AddCommand(descriptorCmd)
	// Add descriptor
	addDescriptorCmd.Flags().StringVar(&descriptorPath, "descriptorPath", "", "Application descriptor path containing a JSON spec")
	addDescriptorCmd.MarkFlagRequired("descriptorPath")
	descriptorCmd.AddCommand(addDescriptorCmd)
	// Get descriptor
	getDescriptorCmd.Flags().StringVar(&descriptorID, "descriptorID", "", "Application descriptor identifier")
	getDescriptorCmd.MarkFlagRequired("descriptorID")
	descriptorCmd.AddCommand(getDescriptorCmd)
	// List descriptors
	descriptorCmd.AddCommand(listDescriptorsCmd)
	// Help
	addDescriptorHelpCmd.Flags().StringVar(&exampleName, "exampleName", "simple", "Example to show: simple or complex or pstorage")
	addDescriptorHelpCmd.Flags().StringVar(&storageType, "storage", "ephemeral", "Type: ephemeral local replica cloud")
	descriptorCmd.AddCommand(addDescriptorHelpCmd)
	// Delete descriptor
	deleteDescriptorCmd.Flags().StringVar(&descriptorID, "descriptorID", "", "Application descriptor identifier")
	deleteDescriptorCmd.MarkFlagRequired("descriptorID")
	descriptorCmd.AddCommand(deleteDescriptorCmd)
	// Application descriptor labels
	appDescLabelsCmd.PersistentFlags().StringVar(&descriptorID, "descriptorID", "", "Application descriptor identifier")
	appDescLabelsCmd.PersistentFlags().StringVar(&rawLabels, "labels", "", "Labels separated by ; as in key1:value;key2:value")
	descriptorCmd.AddCommand(appDescLabelsCmd)
	appDescLabelsCmd.AddCommand(addLabelToAppDescriptorCmd)
	addLabelToAppDescriptorCmd.MarkPersistentFlagRequired("descriptorID")
	appDescLabelsCmd.AddCommand(removeLabelFromAppDescriptorCmd)
	removeLabelFromAppDescriptorCmd.MarkPersistentFlagRequired("descriptorID")

	// Instances
	appsCmd.AddCommand(instanceCmd)
	// List
	instanceCmd.AddCommand(listInstancesCmd)
	// Deploy
	deployInstanceCmd.Flags().StringVar(&name, "name", "", "Name of the application instance")
	deployInstanceCmd.Flags().StringVar(&instanceID, "instanceID", "", "Application instance identifier")
	deployInstanceCmd.MarkFlagRequired("name")
	deployInstanceCmd.MarkFlagRequired("instanceID")
	instanceCmd.AddCommand(deployInstanceCmd)
	// Undeploy
	undeployInstanceCmd.Flags().StringVar(&instanceID, "instanceID", "", "Application instance identifier")
	undeployInstanceCmd.MarkFlagRequired("instanceID")
	instanceCmd.AddCommand(undeployInstanceCmd)
	// Get
	getInstanceCmd.Flags().StringVar(&instanceID, "instanceID", "", "Application instance identifier")
	getInstanceCmd.MarkFlagRequired("instanceID")
	instanceCmd.AddCommand(getInstanceCmd)
}

var descriptorCmd = &cobra.Command{
	Use:     "desc",
	Aliases: []string{"descriptor"},
	Short:   "Manage applications descriptors",
	Long:    `Manage applications descriptors`,
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
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		a.AddDescriptor(options.Resolve("organizationID", organizationID), descriptorPath)
	},
}

var addDescriptorHelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help related to adding a new application descriptor",
	Long:  `Show help related to adding a new application descriptor`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			"",
			0,
			insecure,
			options.Resolve("cacert", caCertPath))
		a.ShowDescriptorHelp(exampleName, storageType)
	},
}

var listDescriptorsCmd = &cobra.Command{
	Use:   "list",
	Short: "List the application descriptors",
	Long:  `List the application descriptors`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		a.ListDescriptors(options.Resolve("organizationID", organizationID))
	},
}

var getDescriptorCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an application descriptor",
	Long:  `Get an application descriptor`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		a.GetDescriptor(options.Resolve("organizationID", organizationID), descriptorID)
	},
}

var appDescLabelsCmd = &cobra.Command{
	Use:   "label",
	Aliases: []string{"labels", "l"},
	Short: "Manage application descriptor labels",
	Long:  `Manage application descriptor labels`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var addLabelToAppDescriptorCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a set of labels to an application descriptor",
	Long:  `Add a set of labels to an application descriptor`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		a.ModifyAppDescriptorLabels(options.Resolve("organizationID", organizationID),
			descriptorID, true, rawLabels)
	},
}

var removeLabelFromAppDescriptorCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a set of labels from an application descriptor",
	Long:  `Remove a set of labels from an application descriptor`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		a.ModifyAppDescriptorLabels(options.Resolve("organizationID", organizationID),
			descriptorID, false, rawLabels)
	},
}

var deleteDescriptorCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an application descriptor",
	Long:  `Delete an application descriptor`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		a.DeleteDescriptor(options.Resolve("organizationID", organizationID), descriptorID)
	},
}

var instanceCmd = &cobra.Command{
	Use:     "inst",
	Aliases: []string{"instance"},
	Short:   "Manage applications instances",
	Long:    `Manage applications instances`,
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
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		a.Deploy(options.Resolve("organizationID", organizationID), descriptorID, name, description)
	},
}

var undeployInstanceCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Undeploy an application instance",
	Long:  `Undeploy an application instance`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		a.Undeploy(options.Resolve("organizationID", organizationID), instanceID)
	},
}

var listInstancesCmd = &cobra.Command{
	Use:   "list",
	Short: "List application instances",
	Long:  `List application intances`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		a.ListInstances(options.Resolve("organizationID", organizationID))
	},
}

var getInstanceCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an application instance",
	Long:  `Get and application instance`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		a.GetInstance(options.Resolve("organizationID", organizationID), instanceID)
	},
}
