/*
 * Copyright 2020 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package commands

import (
	"fmt"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"math"
	"strconv"
)

var clustersCmd = &cobra.Command{
	Use:     "cluster",
	Aliases: []string{"clusters"},
	Short:   "Manage clusters",
	Long:    `Manage clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(clustersCmd)
	installClustersCmd.Flags().StringVar(&kubeConfigPath, "kubeConfigPath", "", "KubeConfig path for installing an existing cluster")
	installClustersCmd.Flags().StringVar(&hostname, "ingressHostname", "", "Hostname of the application cluster ingress")
	installClustersCmd.Flags().StringVar(&username, "username", "", "Username (for clusters requiring the install of Kubernetes)")
	installClustersCmd.Flags().StringVar(&password, "password", "", "Password (for clusters requiring the install of Kubernetes)")
	installClustersCmd.Flags().StringArrayVar(&nodes, "nodes", []string{}, "Nodes (for clusters requiring the install of Kubernetes)")
	installClustersCmd.Flags().StringVar(&targetPlatform, "targetPlatform", "", "Indicate the target platform between MINIKUBE AZURE")
	installClustersCmd.Flags().BoolVar(&useStaticIPAddresses, "useStaticIPAddresses", false,
		"Use statically assigned IP Addresses for the public facing services")
	installClustersCmd.Flags().StringVar(&ipAddressIngress, "ipAddressIngress", "",
		"Public IP Address assigned to the public ingress service")
	clustersCmd.AddCommand(installClustersCmd)

	listClustersCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch for changes")
	clustersCmd.AddCommand(listClustersCmd)

	updateClusterCmd.Flags().Float64Var(&millicoresConversionFactor, "millicoresConversionFactor", math.NaN(), "Modify the millicoresConversionFactor assigned to the cluster")
	clustersCmd.AddCommand(updateClusterCmd)

	clusterLabelsCmd.PersistentFlags().StringVar(&clusterID, "clusterID", "", "Cluster identifier")
	clusterLabelsCmd.PersistentFlags().StringVar(&rawLabels, "labels", "", "Labels separated by ; as in key1:value;key2:value")

	clusterLabelsCmd.AddCommand(addLabelToClusterCmd)
	clusterLabelsCmd.AddCommand(removeLabelFromClusterCmd)
	clustersCmd.AddCommand(clusterLabelsCmd)

	clustersCmd.AddCommand(infoClusterCmd)

	clustersCmd.AddCommand(cordonClusterCmd)
	clustersCmd.AddCommand(uncordonClusterCmd)
	clustersCmd.AddCommand(drainClusterCmd)

	provAndInstCmd.PersistentFlags().StringVar(&organizationID, "organizationId", "", "Organization Id")
	provAndInstCmd.PersistentFlags().StringVar(&provisionClusterName, "clusterName", "", "Cluster name")
	provAndInstCmd.PersistentFlags().StringVar(&provisionAzureCredentialsPath, "azureCredentialsPath", "", "Path for the azure credentials file")
	provAndInstCmd.PersistentFlags().StringVar(&provisionAzureDnsZoneName, "azureDnsZoneName", "", "DNS zone for azure")
	provAndInstCmd.PersistentFlags().StringVar(&provisionAzureResourceGroup, "azureResourceGroup", "", "Azure resource group")
	provAndInstCmd.PersistentFlags().StringVar(&provisionClusterType, "clusterType", "KUBERNETES", "Cluster type")
	provAndInstCmd.PersistentFlags().BoolVar(&provisionIsProductionCluster, "isProductionCluster", false, "Indicate the provisioning of a cluster in a production environment")
	provAndInstCmd.PersistentFlags().StringVar(&provisionKubernetesVersion, "kubernetesVersion", "", "Kubernetes version to be used")
	provAndInstCmd.PersistentFlags().StringVar(&provisionNodeType, "nodeType", "", "Type of node to use")
	provAndInstCmd.PersistentFlags().IntVar(&provisionNumNodes, "numNodes", 1, "Number of nodes")
	provAndInstCmd.PersistentFlags().StringVar(&provisionTargetPlatform, "targetPlatform", "", "Target platform")
	provAndInstCmd.PersistentFlags().StringVar(&provisionZone, "zone", "", "Deployment zone")
	provAndInstCmd.PersistentFlags().StringVar(&provisionKubeConfigOutputPath, "kubeConfigOutputPath", "/tmp", "Path where the kubeconfig will be stored")
	clustersCmd.AddCommand(provAndInstCmd)

	scaleClusterCmd.PersistentFlags().StringVar(&provisionClusterType, "clusterType", "KUBERNETES", "Cluster type")
	scaleClusterCmd.PersistentFlags().StringVar(&provisionAzureCredentialsPath, "azureCredentialsPath", "", "Path for the azure credentials file")
	scaleClusterCmd.PersistentFlags().StringVar(&provisionAzureResourceGroup, "azureResourceGroup", "", "Azure resource group")
	scaleClusterCmd.PersistentFlags().StringVar(&provisionTargetPlatform, "targetPlatform", "", "Target platform")
	clustersCmd.AddCommand(scaleClusterCmd)

	uninstallClusterCmd.Flags().StringVar(&provisionTargetPlatform, "targetPlatform", "AZURE", "Target platform")
	clustersCmd.AddCommand(uninstallClusterCmd)

	decomissionClusterCmd.Flags().StringVar(&provisionClusterType, "clusterType", "KUBERNETES", "Cluster type")
	decomissionClusterCmd.Flags().StringVar(&provisionAzureCredentialsPath, "azureCredentialsPath", "", "Path for the azure credentials file")
	decomissionClusterCmd.Flags().StringVar(&provisionAzureResourceGroup, "azureResourceGroup", "", "Azure resource group")
	decomissionClusterCmd.Flags().StringVar(&provisionTargetPlatform, "targetPlatform", "", "Target platform")
	clustersCmd.AddCommand(decomissionClusterCmd)
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
	Use:     "info <clusterID>",
	Aliases: []string{"get"},
	Short:   "Get the cluster information",
	Long:    `Get the cluster information`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		c.Info(cliOptions.Resolve("organizationID", organizationID), args[0])
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
		c.List(cliOptions.Resolve("organizationID", organizationID), watch)
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

var updateClusterCmd = &cobra.Command{
	Use:   "update <clusterID>",
	Short: "Update cluster params",
	Long:  `Update cluster params`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		c.Update(cliOptions.Resolve("organizationID", organizationID), args[0], "", millicoresConversionFactor)
	},
}

var clusterLabelsCmd = &cobra.Command{
	Use:     "label",
	Aliases: []string{"labels", "l"},
	Short:   "Manage cluster labels",
	Long:    `Manage cluster labels`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		_ = cmd.Help()
	},
}

var addLabelToClusterCmd = &cobra.Command{
	Use:   "add <clusterID> <labels>",
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
			_ = cmd.Help()
		} else {

			targetValues, err := ResolveArgument([]string{"clusterID", "labels"}, args, []string{clusterID, rawLabels})
			if err != nil {
				fmt.Println(err.Error())
				_ = cmd.Help()
			} else {
				c.ModifyClusterLabels(cliOptions.Resolve("organizationID", organizationID),
					targetValues[0], true, targetValues[1])
			}
		}
	},
}

var removeLabelFromClusterCmd = &cobra.Command{
	Use:     "delete <clusterID> <labels>",
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
			_ = cmd.Help()
		} else {

			targetValues, err := ResolveArgument([]string{"clusterID", "labels"}, args, []string{clusterID, rawLabels})
			if err != nil {
				fmt.Println(err.Error())
				_ = cmd.Help()
			} else {
				c.ModifyClusterLabels(cliOptions.Resolve("organizationID", organizationID),
					targetValues[0], false, targetValues[1])
			}
		}
	},
}

var cordonClusterCmd = &cobra.Command{
	Use:   "cordon <clusterID>",
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
	Use:   "uncordon <clusterID>",
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
	Use:   "drain <clusterID>",
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
	Use:     "provision",
	Aliases: []string{"provision-and-install", "pai"},
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

		p.ProvisionAndInstall(cliOptions.Resolve("organizationId", organizationID),
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

// scaleClusterCmd with the cmd definition of a scale cluster operation.
var scaleClusterCmd = &cobra.Command{
	Use:   "scale <clusterID> <numNodes>",
	Short: "Scale an application cluster",
	Long:  `Scale an application cluster`,
	Args:  cobra.ExactArgs(2),
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

		numNodes, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal().Msg("numNodes must be a number")
		}

		p.Scale(cliOptions.Resolve("organizationId", organizationID),
			args[0],
			clusterType,
			int64(numNodes),
			targetPlatform,
			provisionAzureCredentialsPath,
			provisionAzureResourceGroup,
		)
	},
}

var uninstallClusterLongHelp = `
Uninstall the Nalej components deployed in the target cluster.

This command will remove the Nalej components deployed in an application
cluster as created by the installing process. Notice that this operation
does not free the computing resources associated with the cluster. To
completelly uninstall the platform and free the associated computing resources
use the decomission command as:

public-api-cli cluster decomission ...

`

var uninstallClusterExamples = `
# Uninstall an application cluster
public-api-cli cluster uninstall 00630f9c-59fe-408a-829c-6dc67c2b98e7 nalej/appCluster.yaml
`

// uninstallClusterCmd with the cmd definition of a uninstall cluster operation.
var uninstallClusterCmd = &cobra.Command{
	Use:     "uninstall <clusterID> <kubeConfigPath>",
	Short:   "Uninstall an application cluster",
	Long:    uninstallClusterLongHelp,
	Example: uninstallClusterExamples,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()

		c := cli.NewClusters(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetPlatform := stringToTargetPlatform(provisionTargetPlatform)
		c.Uninstall(cliOptions.Resolve("organizationId", organizationID),
			args[0], args[1], targetPlatform)
	},
}

var decomissionClusterLongHelp = `
Decomission an application cluster.

This command will perform an uninstall of the Nalej components deployed in
the application cluster. Once all components are uninstalled, it will trigger
the decomissioning process so that computing resources are freed.

The decomissioning process depends on the infrastructure provider used to
host the cluster, and valid credentials for that provider must be
available to execute this operation. Once the cluster is decomissioned, it
will be removed from the list of application clusters.
`

var decomissionClusterExamples = `
# Decomission an Azure application cluster
public-api-cli cluster decomission 00630f9c-59fe-408a-829c-6dc67c2b98e7 --targetPlatform AZURE --azureCredentialsPath azure/credentials.json --azureResourceGroup dev
`

// decomissionClusterCmd with the cmd definition of a decomission cluster operation.
var decomissionClusterCmd = &cobra.Command{
	Use:     "decomission <clusterID>",
	Short:   "decomission an application cluster",
	Long:    decomissionClusterLongHelp,
	Example: decomissionClusterExamples,
	Args:    cobra.ExactArgs(1),
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

		p.Decomission(cliOptions.Resolve("organizationId", organizationID),
			args[0], clusterType, targetPlatform,
			provisionAzureCredentialsPath, provisionAzureResourceGroup)

	},
}
