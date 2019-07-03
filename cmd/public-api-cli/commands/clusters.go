/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"fmt"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var clustersCmd = &cobra.Command{
	Use:     "cluster",
	Aliases: []string{"clusters"},
	Short:   "Manage clusters",
	Long:    `Manage clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(clustersCmd)
	installClustersCmd.Flags().StringVar(&kubeConfigPath, "kubeConfigPath", "", "KubeConfig path for installing an existing cluster")
	installClustersCmd.Flags().StringVar(&hostname, "ingressHostname", "", "Hostname of the application cluster ingress")
	installClustersCmd.Flags().StringVar(&username, "username", "", "Username (for clusters requiring the install of Kubernetes)")
	installClustersCmd.Flags().StringVar(&password, "password", "", "Password (for clusters requiring the install of Kubernetes)")
	installClustersCmd.Flags().StringArrayVar(&nodes, "nodes", []string{}, "Nodes (for clusters requiring the install of Kubernetes)")
	installClustersCmd.Flags().StringVar(&targetPlatform, "targetPlatform", "MINIKUBE", "Indicate the target platform between MINIKUBE AZURE")
	installClustersCmd.Flags().BoolVar(&useStaticIPAddresses, "useStaticIPAddresses", false,
		"Use statically assigned IP Addresses for the public facing services")
	installClustersCmd.Flags().StringVar(&ipAddressIngress, "ipAddressIngress", "",
		"Public IP Address assigned to the public ingress service")
	clustersCmd.AddCommand(installClustersCmd)
	clustersCmd.AddCommand(listClustersCmd)

	clusterLabelsCmd.PersistentFlags().StringVar(&clusterID, "clusterID", "", "Cluster identifier")
	clusterLabelsCmd.PersistentFlags().StringVar(&rawLabels, "labels", "", "Labels separated by ; as in key1:value;key2:value")

	clusterLabelsCmd.AddCommand(addLabelToClusterCmd)
	clusterLabelsCmd.AddCommand(removeLabelFromClusterCmd)
	clustersCmd.AddCommand(clusterLabelsCmd)

	clustersCmd.AddCommand(infoClusterCmd)
	infoClusterCmd.Flags().StringVar(&clusterID, "clusterID", "", "Cluster identifier")

	clustersCmd.AddCommand(monitorClusterCmd)
	monitorClusterCmd.Flags().StringVar(&clusterID, "clusterID", "", "Cluster identifier")
	monitorClusterCmd.Flags().Int32Var(&rangeMinutes, "rangeMinutes", 0, "Return average values over the past <rangeMinutes> minutes")

	clustersCmd.AddCommand(cordonClusterCmd)

	clustersCmd.AddCommand(uncordonClusterCmd)

	clustersCmd.AddCommand(drainClusterCmd)

}

var installClustersCmd = &cobra.Command{
	Use:   "install",
	Short: "Install an application cluster",
	Long:  `Install an application cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		c.Install(
			options.Resolve("organizationID", organizationID),
			kubeConfigPath,
			hostname,
			username,
			privateKeyPath,
			nodes,
			stringToTargetPlatform(targetPlatform),
			useStaticIPAddresses,
			ipAddressIngress)
	},
}

var infoClusterCmd = &cobra.Command{
	Use:   "info [clusterID]",
	Aliases: []string{"get"},
	Short: "Get the cluster information",
	Long:  `Get the cluster information`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"clusterID"}, args, []string{clusterID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
		c.Info(options.Resolve("organizationID", organizationID), options.Resolve("clusterID", targetValues[0]))
		}

	},
}

var listClustersCmd = &cobra.Command{
	Use:   "list",
	Aliases: []string{"ls"},
	Short: "List clusters",
	Long:  `List clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		c.List(options.Resolve("organizationID", organizationID))
	},
}

var monitorClusterCmd = &cobra.Command{
	Use:   "monitor [clusterID]",
	Aliases: []string{"mon"},
	Short: "Monitor cluster",
	Long:  `Get summarized monitoring information for a single cluster`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"clusterID"}, args, []string{clusterID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			c.Monitor(
				options.Resolve("organizationID", organizationID),
				options.Resolve("clusterID", targetValues[0]),
				rangeMinutes,
			)
		}
	},
}

// Convert a string to the corresponding cluster platform
func stringToTargetPlatform(p string) grpc_public_api_go.Platform {
	var result grpc_public_api_go.Platform
	switch p {
	case grpc_public_api_go.Platform_AZURE.String():
		result = grpc_public_api_go.Platform_AZURE
	case grpc_public_api_go.Platform_MINIKUBE.String():
		result = grpc_public_api_go.Platform_MINIKUBE
	default:
		log.Fatal().Str("platform", p).Msg("unknown platform")
	}

	return result
}

var clusterLabelsCmd = &cobra.Command{
	Use:   "label",
	Aliases: []string{"labels", "l"},
	Short: "Manage cluster labels",
	Long:  `Manage cluster labels`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var addLabelToClusterCmd = &cobra.Command{
	Use:   "add [clusterID] [labels]",
	Short: "Add a set of labels to a cluster",
	Long:  `Add a set of labels to a cluster`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

		targetValues, err := ResolveArgument([]string{"clusterID", "labels"}, args, []string{clusterID, rawLabels})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			c.ModifyClusterLabels(options.Resolve("organizationID", organizationID),
				targetValues[0], true, targetValues[1])
		}
	},
}

var removeLabelFromClusterCmd = &cobra.Command{
	Use:   "delete [clusterID] [labels]",
	Aliases: []string{"remove", "del", "rm"},
	Short: "Remove a set of labels from a cluster",
	Long:  `Remove a set of labels from a cluster`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		targetValues, err := ResolveArgument([]string{"clusterID", "labels"}, args, []string{clusterID, rawLabels})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			c.ModifyClusterLabels(options.Resolve("organizationID", organizationID),
				targetValues[0], false, targetValues[1])
		}
	},
}

var cordonClusterCmd = &cobra.Command{
	Use: "cordon [clusterID]",
	Short: "cordon a cluster ignoring new application deployments",
	Long: `cordon a cluster ignoring new application deployments`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		c.CordonCluster(options.Resolve("organizationID", organizationID),args[0])
	},
}

var uncordonClusterCmd = &cobra.Command{
	Use: "uncordon [clusterID]",
	Short: "uncordon a cluster making possible new application deployments",
	Long: `uncordon a cluster making possible new application deployments`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		c.UncordonCluster(options.Resolve("organizationID", organizationID),args[0])
	},
}

var drainClusterCmd = &cobra.Command{
	Use: "drain [clusterID]",
	Short: "drain a cluster",
	Long: `drain a cordoned cluster and force current applications to be re-scheduled`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		c.DrainCluster(options.Resolve("organizationID", organizationID),args[0])
	},
}

