/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
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
	installClustersCmd.Flags().StringVar(&username, "username", "", "Username (for clusters requiring the install of Kubernetes)")
	installClustersCmd.Flags().StringVar(&password, "password", "", "Password (for clusters requiring the install of Kubernetes)")
	installClustersCmd.Flags().StringArrayVar(&nodes, "nodes", []string{}, "Nodes (for clusters requiring the install of Kubernetes)")
	clustersCmd.AddCommand(installClustersCmd)
	clustersCmd.AddCommand(listClustersCmd)
}

var installClustersCmd = &cobra.Command{
	Use:   "install",
	Short: "Install an application cluster",
	Long:  `Install an application cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		c.Install(
			options.Resolve("organizationID", organizationID),
			options.Resolve("clusterID", clusterID),
			kubeConfigPath,
			username,
			privateKeyPath,
			nodes)
	},
}

var listClustersCmd = &cobra.Command{
	Use:   "list",
	Short: "List clusters",
	Long:  `List clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		c := cli.NewClusters(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		c.List(options.Resolve("organizationID", organizationID))
	},
}
