/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"fmt"
	"github.com/nalej/grpc-infrastructure-go"
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

	provAndInstCmd.PersistentFlags().StringVar(&organizationID, "organizationId", "", "Organization Id")
	provAndInstCmd.PersistentFlags().StringVar(&provisionClusterName, "clusterName", "", "Cluster name")
	provAndInstCmd.PersistentFlags().StringVar(&provisionAzureCredentialsPath, "azureCredentialsPath", "", "Path for the azure credentials file")
	provAndInstCmd.PersistentFlags().StringVar(&provisionAzureDnsZoneName, "azureDnsZoneName", "", "DNS zone for azure")
	provAndInstCmd.PersistentFlags().StringVar(&provisionAzureResourceGroup, "azureResourceGroup", "", "Azure resource group")
	provAndInstCmd.PersistentFlags().StringVar(&provisionClusterType, "clusterType", "kubernetes", "Cluster type")
	provAndInstCmd.PersistentFlags().BoolVar(&provisionIsProductionCluster, "isProductionCluster", false, "Indicate the provisioning of a cluster in a production environment")
	provAndInstCmd.PersistentFlags().StringVar(&provisionKubernetesVersion, "kubernetesVersion", "", "Kubernetes version to be used")
	provAndInstCmd.PersistentFlags().StringVar(&provisionNodeType, "nodeType", "", "Type of node to use")
	provAndInstCmd.PersistentFlags().IntVar(&provisionNumNodes, "numNodes", 1, "Number of nodes")
	provAndInstCmd.PersistentFlags().StringVar(&provisionTargetPlatform, "targetPlatform", "", "Target platform")
	provAndInstCmd.PersistentFlags().StringVar(&provisionZone, "zone", "", "Deployment zone")
	provAndInstCmd.PersistentFlags().StringVar(&provisionKubeConfigOutputPath, "kubeConfigOutputPath", "/tmp", "Path where the kubeconfig will be stored")
	clustersCmd.AddCommand(provAndInstCmd)
}

var installClustersCmd = &cobra.Command{
	Use:   "install",
	Short: "Install an application cluster",
	Long:  `Install an application cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		c.Install(
			cliOptions.Resolve("organizationID", organizationID),
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
	Use:     "info [clusterID]",
	Aliases: []string{"get"},
	Short:   "Get the cluster information",
	Long:    `Get the cluster information`,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"clusterID"}, args, []string{clusterID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			c.Info(cliOptions.Resolve("organizationID", organizationID), cliOptions.Resolve("clusterID", targetValues[0]))
		}

	},
}

var listClustersCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List clusters",
	Long:    `List clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		c.List(cliOptions.Resolve("organizationID", organizationID))
	},
}

var monitorClusterCmd = &cobra.Command{
	Use:     "monitor [clusterID]",
	Aliases: []string{"mon"},
	Short:   "Monitor cluster",
	Long:    `Get summarized monitoring information for a single cluster`,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"clusterID"}, args, []string{clusterID})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		} else {
			c.Monitor(
				cliOptions.Resolve("organizationID", organizationID),
				cliOptions.Resolve("clusterID", targetValues[0]),
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

// Convert a string to the corresponding cluster type
func stringToClusterType(ct string) grpc_infrastructure_go.ClusterType {
	var result grpc_infrastructure_go.ClusterType
	switch ct {
	case grpc_infrastructure_go.ClusterType_KUBERNETES.String():
		result = grpc_infrastructure_go.ClusterType_KUBERNETES
	case grpc_infrastructure_go.ClusterType_DOCKER_NODE.String():
		result = grpc_infrastructure_go.ClusterType_DOCKER_NODE
	default:
		log.Fatal().Str("cluster_type", ct).Msg("unknown cluster type")
	}

	return result
}

var clusterLabelsCmd = &cobra.Command{
	Use:     "label",
	Aliases: []string{"labels", "l"},
	Short:   "Manage cluster labels",
	Long:    `Manage cluster labels`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var addLabelToClusterCmd = &cobra.Command{
	Use:   "add [clusterID] [labels]",
	Short: "Add a set of labels to a cluster",
	Long:  `Add a set of labels to a cluster`,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		if len(args) > 0 && len(args) < 2 {
			fmt.Println("[clusterID] and [labels] must be flags or arguments, both the same type")
			cmd.Help()
		} else {

			targetValues, err := ResolveArgument([]string{"clusterID", "labels"}, args, []string{clusterID, rawLabels})
			if err != nil {
				fmt.Println(err.Error())
				cmd.Help()
			} else {
				c.ModifyClusterLabels(cliOptions.Resolve("organizationID", organizationID),
					targetValues[0], true, targetValues[1])
			}
		}
	},
}

var removeLabelFromClusterCmd = &cobra.Command{
	Use:     "delete [clusterID] [labels]",
	Aliases: []string{"remove", "del", "rm"},
	Short:   "Remove a set of labels from a cluster",
	Long:    `Remove a set of labels from a cluster`,
	Args:    cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		if len(args) > 0 && len(args) < 2 {
			fmt.Println("[clusterID] and [labels] must be flags or arguments, both the same type")
			cmd.Help()
		} else {

			targetValues, err := ResolveArgument([]string{"clusterID", "labels"}, args, []string{clusterID, rawLabels})
			if err != nil {
				fmt.Println(err.Error())
				cmd.Help()
			} else {
				c.ModifyClusterLabels(cliOptions.Resolve("organizationID", organizationID),
					targetValues[0], false, targetValues[1])
			}
		}
	},
}

var cordonClusterCmd = &cobra.Command{
	Use:   "cordon [clusterID]",
	Short: "cordon a cluster ignoring new application deployments",
	Long:  `cordon a cluster ignoring new application deployments`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		c.CordonCluster(cliOptions.Resolve("organizationID", organizationID), args[0])
	},
}

var uncordonClusterCmd = &cobra.Command{
	Use:   "uncordon [clusterID]",
	Short: "uncordon a cluster making possible new application deployments",
	Long:  `uncordon a cluster making possible new application deployments`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		c.UncordonCluster(cliOptions.Resolve("organizationID", organizationID), args[0])
	},
}

var drainClusterCmd = &cobra.Command{
	Use:   "drain [clusterID]",
	Short: "drain a cluster",
	Long:  `drain a cordoned cluster and force current applications to be re-scheduled`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		c.DrainCluster(cliOptions.Resolve("organizationID", organizationID), args[0])
	},
}

var provAndInstCmd = &cobra.Command{
	Use:     "provision-and-install",
	Aliases: []string{"pai"},
	Short:   "Provision and install a new cluster",
	Long:    `Provision and install a new cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		p := cli.NewProvision(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output),
			cliOptions.ResolveAsInt("labelLength", labelLength),
			provisionKubeConfigOutputPath)

		clusterType := stringToClusterType(provisionClusterType)
		targetPlatform := stringToTargetPlatform(provisionTargetPlatform)

		p.Cluster(cliOptions.Resolve("organizationId", organizationID),
			provisionClusterName,
			provisionAzureCredentialsPath,
			provisionAzureDnsZoneName,
			provisionAzureResourceGroup,
			clusterType,
			false, // management clusters cannot be installed from the public-api
			provisionIsProductionCluster,
			provisionKubernetesVersion,
			provisionNodeType,
			int64(provisionNumNodes),
			targetPlatform,
			provisionZone,
		)
	},
}
