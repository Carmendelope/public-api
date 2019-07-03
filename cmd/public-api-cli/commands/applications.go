/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// TODO Remove descriptor NP-338
// TODO Remove Args: cobra.MaximumNArgs(1), when the flags are not longer deprecated

package commands

import (
	"fmt"
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:     "application",
	Aliases: []string{"app", "applications"},
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
	addDescriptorCmd.Flags().MarkDeprecated("descriptorPath", "Use command argument instead")
	descriptorCmd.AddCommand(addDescriptorCmd)
	// Get descriptor
	getDescriptorCmd.Flags().StringVar(&descriptorID, "descriptorID", "", "Application descriptor identifier")
	getDescriptorCmd.Flags().MarkDeprecated("descriptorID", "Use command argument instead")
	descriptorCmd.AddCommand(getDescriptorCmd)
	// List descriptors
	descriptorCmd.AddCommand(listDescriptorsCmd)
	// Help
	addDescriptorHelpCmd.Flags().StringVar(&exampleName, "exampleName", "simple", "Example to show: simple or complex or pstorage")
	addDescriptorHelpCmd.Flags().StringVar(&storageType, "storage", "ephemeral", "Type: ephemeral local replica cloud")
	descriptorCmd.AddCommand(addDescriptorHelpCmd)
	// Delete descriptor
	deleteDescriptorCmd.Flags().StringVar(&descriptorID, "descriptorID", "", "Application descriptor identifier")
	deleteDescriptorCmd.Flags().MarkDeprecated("descriptorID", "Use command argument instead")
	descriptorCmd.AddCommand(deleteDescriptorCmd)
	// Application descriptor labels
	appDescLabelsCmd.PersistentFlags().StringVar(&descriptorID, "descriptorID", "", "Application descriptor identifier")
	appDescLabelsCmd.PersistentFlags().StringVar(&rawLabels, "labels", "", "Labels separated by ; as in key1:value;key2:value")
	descriptorCmd.AddCommand(appDescLabelsCmd)
	appDescLabelsCmd.AddCommand(addLabelToAppDescriptorCmd)
	appDescLabelsCmd.AddCommand(removeLabelFromAppDescriptorCmd)
	// List descriptor Parameters
	descriptorCmd.AddCommand(getDescriptorParamsCmd)


	// Instances
	appsCmd.AddCommand(instanceCmd)
	// List
	instanceCmd.AddCommand(listInstancesCmd)
	// Deploy
	deployInstanceCmd.Flags().StringVar(&name, "name", "", "Name of the application instance")
	deployInstanceCmd.Flags().StringVar(&descriptorID, "descriptorID", "", "Application instance identifier")
	deployInstanceCmd.Flags().StringVar(&params, "params", "", "Param values to deploy (param1=value1,...,paramN=valueN)")
	deployInstanceCmd.Flags().MarkDeprecated("name", "Use command argument instead")
	deployInstanceCmd.Flags().MarkDeprecated("descriptorID", "Use command argument instead")
	instanceCmd.AddCommand(deployInstanceCmd)
	// Undeploy
	undeployInstanceCmd.Flags().StringVar(&instanceID, "instanceID", "", "Application instance identifier")
	undeployInstanceCmd.Flags().MarkDeprecated("instanceID", "Use command argument instead")
	instanceCmd.AddCommand(undeployInstanceCmd)
	// Get
	getInstanceCmd.Flags().StringVar(&instanceID, "instanceID", "", "Application instance identifier")
	getInstanceCmd.Flags().MarkDeprecated("instanceID", "Use command argument instead")
	getInstanceCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch for changes")
	instanceCmd.AddCommand(getInstanceCmd)
	// List instance params
	instanceCmd.AddCommand(getInstanceParamsCmd)
}

var descriptorCmd = &cobra.Command{
	Use:     "descriptor",
	Aliases: []string{"desc", "descriptors"},
	Short:   "Manage applications descriptors",
	Long:    `Manage applications descriptors`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var addDescriptorCmd = &cobra.Command{
	Use:   "add [descriptorPath]",
	Short: "Add a new application descriptor",
	Long:  `Add a new application descriptor`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		targetDescriptorPath, err := ResolveArgument([]string{"descriptorPath"}, args, []string{descriptorPath})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else {
			a.AddDescriptor(options.Resolve("organizationID", organizationID), targetDescriptorPath[0])
		}

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
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
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
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		a.ListDescriptors(options.Resolve("organizationID", organizationID))
	},
}

var getDescriptorCmd = &cobra.Command{
	Use:   "get [descriptorID]",
	Short: "Get an application descriptor",
	Long:  `Get an application descriptor`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		targetValues, err := ResolveArgument([]string{"descriptorID"}, args, []string{descriptorID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			a.GetDescriptor(options.Resolve("organizationID", organizationID), targetValues[0])
		}
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
	Use:   "add [descriptorID] [labels]",
	Short: "Add a set of labels to an application descriptor",
	Long:  `Add a set of labels to an application descriptor`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"descriptorID", "labels"}, args, []string{descriptorID, rawLabels})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			a.ModifyAppDescriptorLabels(options.Resolve("organizationID", organizationID),
				targetValues[0], true, targetValues[1])
		}
	},
}

var removeLabelFromAppDescriptorCmd = &cobra.Command{
	Use:   "delete [descriptorID] [labels]",
	Aliases: []string{"remove", "del", "rm"},
	Short: "Remove a set of labels from an application descriptor",
	Long:  `Remove a set of labels from an application descriptor`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"descriptorID", "labels"}, args, []string{descriptorID, rawLabels})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			a.ModifyAppDescriptorLabels(options.Resolve("organizationID", organizationID),
				targetValues[0], false, targetValues[1])
		}
	},
}

var deleteDescriptorCmd = &cobra.Command{
	Use:   "delete [descriptorID]",
	Aliases: []string{"remove", "del", "rm"},
	Short: "Delete an application descriptor",
	Long:  `Delete an application descriptor`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		targetValues, err := ResolveArgument([]string{"descriptorID"}, args, []string{descriptorID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			a.DeleteDescriptor(options.Resolve("organizationID", organizationID), targetValues[0])
		}

	},
}

var getDescriptorParamsCmd = &cobra.Command{
	Use:   "parameters [descriptorID]",
	Aliases: []string{"params", "param", "parameter"},
	Short: "list parameters of a descriptor",
	Long:  "list parameters of a descriptor",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		targetDescriptorID, err := ResolveArgument([]string{"descriptorID"}, args, []string{descriptorID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			a.GetDescriptorParameters(options.Resolve("organizationID", organizationID),targetDescriptorID[0])
		}
	},
}

var instanceCmd = &cobra.Command{
	Use:     "instance",
	Aliases: []string{"inst"},
	Short:   "Manage applications instances",
	Long:    `Manage applications instances`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var deployInstanceCmd = &cobra.Command{
	Use:   "deploy [descriptorID] [name]",
	Short: "Deploy an application instance",
	Long:  `Deploy an application instance`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"descriptorID", "name"}, args, []string{descriptorID, name})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			a.Deploy(options.Resolve("organizationID", organizationID), targetValues[0], targetValues[1], params)
		}

	},
}

var undeployInstanceCmd = &cobra.Command{
	Use:   "undeploy [instanceID]",
	Short: "Undeploy an application instance",
	Long:  `Undeploy an application instance`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		targetInstanceID, err := ResolveArgument([]string{"instanceID"}, args, []string{instanceID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			a.Undeploy(options.Resolve("organizationID", organizationID), targetInstanceID[0])
		}
	},
}

var listInstancesCmd = &cobra.Command{
	Use:   "list",
	Aliases: []string{"ls"},
	Short: "List application instances",
	Long:  `List application intances`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		a.ListInstances(options.Resolve("organizationID", organizationID))
	},
}

var getInstanceCmd = &cobra.Command{
	Use:   "get [instanceID]",
	Short: "Get an application instance",
	Long:  `Get and application instance`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		targetInstanceID, err := ResolveArgument([]string{"instanceID"}, args, []string{instanceID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			a.GetInstance(options.Resolve("organizationID", organizationID), targetInstanceID[0], watch)
		}
	},
}

var getInstanceParamsCmd = &cobra.Command{
	Use:   "parameterss [instanceID]",
	Aliases: []string{"params", "param", "parameter"},
	Short: "list parameters of an instance",
	Long:  "list parameters of an instance",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		a := cli.NewApplications(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		targetInstanceID, err := ResolveArgument([]string{"instanceID"}, args, []string{instanceID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			a.GetInstanceParameters(options.Resolve("organizationID", organizationID),targetInstanceID[0])
		}
	},
}