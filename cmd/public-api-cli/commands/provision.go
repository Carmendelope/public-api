/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	provisionClusterName string
	provisionAzureCredentialsPath string
	provisionAzureDnsZoneName string
	provisionAzureResourceGroup string
	provisionClusterType string
	provisionIsProductionCluster bool
	provisionKubernetesVersion string
	provisionNodeType string
	provisionNumNodes int
	provisionTargetPlatform string
	provisionZone string


)

// Conversion map for RawClusterTypes
var ClusterTypeFromRaw map[string]grpc_infrastructure_go.ClusterType = map[string]grpc_infrastructure_go.ClusterType{
	"kubernetes": grpc_infrastructure_go.ClusterType_KUBERNETES,
	"docker": grpc_infrastructure_go.ClusterType_DOCKER_NODE,
}

// Conversion map for Installation target platforms
var TargetPlatformFromRaw map[string]grpc_installer_go.Platform = map[string]grpc_installer_go.Platform{
	"azure": grpc_installer_go.Platform_AZURE,
}

var provisionCmd = &cobra.Command{
	Use:     "provision",
	Aliases: []string{"provision"},
	Short:   "Provision resources",
	Long:    `Provision resources`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(provisionCmd)

	provisionClusterCmd.PersistentFlags().StringVar(&organizationID, "organizationId", "", "Organization Id")
	provisionClusterCmd.PersistentFlags().StringVar(&provisionClusterName, "clusterName", "", "Cluster name")
	provisionClusterCmd.PersistentFlags().StringVar(&provisionAzureCredentialsPath, "azureCredentialsPath", "", "Path for the azure credentials file")
	provisionClusterCmd.PersistentFlags().StringVar(&provisionAzureDnsZoneName, "azureDnsZoneName", "", "DNS zone for azure")
	provisionClusterCmd.PersistentFlags().StringVar(&provisionAzureResourceGroup, "azureResourceGroup", "", "Azure resource group")
	provisionClusterCmd.PersistentFlags().StringVar(&provisionClusterType, "clusterType", "kubernetes", "Cluster type")
	provisionClusterCmd.PersistentFlags().BoolVar(&provisionIsProductionCluster, "isProductionCluster", false, "Indicate the provisioning of a cluster in a production environment")
	provisionClusterCmd.PersistentFlags().StringVar(&provisionKubernetesVersion, "kubernetesVersion", "", "Kubernetes version to be used")
	provisionClusterCmd.PersistentFlags().StringVar(&provisionNodeType, "nodeType", "", "Type of node to use")
	provisionClusterCmd.PersistentFlags().IntVar(&provisionNumNodes, "numNodes", 1, "Number of nodes")
	provisionClusterCmd.PersistentFlags().StringVar(&provisionTargetPlatform, "targetPlatform", "", "Target platform")
	provisionClusterCmd.PersistentFlags().StringVar(&provisionZone, "zone", "", "Deployment zone")

	provisionClusterCmd.MarkFlagRequired("organizationId")
	provisionClusterCmd.MarkFlagRequired("clusterName")
	provisionClusterCmd.MarkFlagRequired("azureCredentialsPath")
	provisionClusterCmd.MarkFlagRequired("azureDnsZoneName")
	provisionClusterCmd.MarkFlagRequired("azureResourceGroup")
	provisionClusterCmd.MarkFlagRequired("kubernetesVersion")
	provisionClusterCmd.MarkFlagRequired("numNodes")
	provisionClusterCmd.MarkFlagRequired("targetPlatform")
	provisionClusterCmd.MarkFlagRequired("nodeType")
	provisionClusterCmd.MarkFlagRequired("provisionZone")

	provisionCmd.AddCommand(provisionClusterCmd)

	checkProgressCmd.PersistentFlags().StringVar(&requestId, "requestId", "", "Request ID")
	checkProgressCmd.MarkFlagRequired("requestId")
	provisionCmd.AddCommand(checkProgressCmd)

	removeProvisionCmd.PersistentFlags().StringVar(&requestId, "requestId", "", "Request ID")
	removeProvisionCmd.MarkFlagRequired("requestId")
	provisionCmd.AddCommand(removeProvisionCmd)


}


var provisionClusterCmd = &cobra.Command{
	Use:   "cluster",
	Aliases: []string{"cluster"},
	Short: "Provision a new cluster",
	Long:  `Provision a new cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		p := cli.NewProvision(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output),
			cliOptions.ResolveAsInt("labelLength", labelLength))

		clusterType, found := ClusterTypeFromRaw[provisionClusterType]
		if !found {
			log.Fatal().Msg("invalid cluster type")
		}

		targetPlatform, found := TargetPlatformFromRaw[provisionTargetPlatform]
		if !found {
			log.Fatal().Msg("invalid target platform")
		}

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

var checkProgressCmd = &cobra.Command{
	Use: "check",
	Aliases: []string{"check"},
	Short: "Check cluster provision",
	Long: "Check cluster provision",
	Run: func(cmd *cobra.Command, args[] string) {
		SetupLogging()
		p := cli.NewProvision(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output),
			cliOptions.ResolveAsInt("labelLength", labelLength))
		p.CheckProgress(cliOptions.Resolve("requestId", requestId))
	},
}

var removeProvisionCmd = &cobra.Command{
	Use: "remove",
	Aliases: []string{"remove"},
	Short: "Remove cluster provision request",
	Long: "Remove cluster provision request",
	Run: func(cmd *cobra.Command, args[] string) {
		SetupLogging()
		p := cli.NewProvision(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output),
			cliOptions.ResolveAsInt("labelLength", labelLength))
		p.RemoveProvision(cliOptions.Resolve("requestId", requestId))
	},
}

