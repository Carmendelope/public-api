/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
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
	installClustersCmd.Flags().StringVar(&clusterID, "clusterID", "", "Cluster identifier")
	installClustersCmd.Flags().StringVar(&kubeConfigPath, "kubeConfigPath", "", "KubeConfig path for installing an existing cluster")
	installClustersCmd.Flags().StringVar(&hostname, "ingressHostname", "", "Hostname of the application cluster ingress")
	installClustersCmd.Flags().StringVar(&username, "username", "", "Username (for clusters requiring the install of Kubernetes)")
	installClustersCmd.Flags().StringVar(&password, "password", "", "Password (for clusters requiring the install of Kubernetes)")
	installClustersCmd.Flags().StringArrayVar(&nodes, "nodes", []string{}, "Nodes (for clusters requiring the install of Kubernetes)")
	installClustersCmd.Flags().BoolVar(&useCoreDNS, "useCoreDNS", true, "Indicate if CoreDNS is going to be used. If not, kubeDNS will be set")
	installClustersCmd.Flags().StringVar(&targetPlatform, "targetPlatform", "minikube", "Indicate the target platform between minikube azure")
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
			insecure,
			options.Resolve("cacert", caCertPath))
		c.Install(
			options.Resolve("organizationID", organizationID),
			options.Resolve("clusterID", clusterID),
			kubeConfigPath,
			hostname,
			username,
			privateKeyPath,
			nodes,
			useCoreDNS,
			stringToTargetPlatform(targetPlatform),
			useStaticIPAddresses,
			ipAddressIngress)
	},
}

var infoClusterCmd = &cobra.Command{
	Use:   "info",
	Short: "Get the cluster information",
	Long:  `Get the cluster information`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		c.Info(options.Resolve("organizationID", organizationID), options.Resolve("clusterID", clusterID))
	},
}

var listClustersCmd = &cobra.Command{
	Use:   "list",
	Short: "List clusters",
	Long:  `List clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		c.List(options.Resolve("organizationID", organizationID))
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
	Use:   "add",
	Short: "Add a set of labels to a cluster",
	Long:  `Add a set of labels to a cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		c.ModifyClusterLabels(options.Resolve("organizationID", organizationID),
			clusterID, true, rawLabels)
	},
}

var removeLabelFromClusterCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a set of labels from a cluster",
	Long:  `Remove a set of labels from a cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		c.ModifyClusterLabels(options.Resolve("organizationID", organizationID),
			clusterID, false, rawLabels)
	},
}