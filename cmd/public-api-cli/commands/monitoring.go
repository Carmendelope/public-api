/*
 * Copyright 2019 Nalej
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
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var monitoringCmd = &cobra.Command{
	Use:     "monitoring",
	Aliases: []string{"mon"},
	Short:   "Get monitoring statistics",
	Long:    `Application Network related operations`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(monitoringCmd)

	monitoringCmd.AddCommand(clusterStatsCmd)
	clusterStatsCmd.Flags().Int32Var(&rangeMinutes, "range-minutes", 0, "Return average values over the past <rangeMinutes> minutes.")
	clusterStatsCmd.Flags().StringVar(&clusterStatFields, "fields", "", "Fields of the cluster stats (SERVICES, VOLUMES, FRAGMENTS, ENDPOINTS) comma separated.")

	monitoringCmd.AddCommand(clusterSummaryCmd)
	clusterSummaryCmd.Flags().Int32Var(&rangeMinutes, "range-minutes", 0, "Return average values over the past <rangeMinutes> minutes.")
}

var clusterStatsCmd = &cobra.Command{
	Use:   "cluster-stats [clusterID]",
	Short: "Display the cluster stats",
	Long:  `Display the cluster stats`,

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		monitoring := cli.NewMonitoring(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		monitoring.GetClusterStats(cliOptions.Resolve("organizationID", organizationID), args[0], rangeMinutes, clusterStatFields)
	},
}

var clusterSummaryCmd = &cobra.Command{
	Use:   "cluster-summary [clusterID]",
	Short: "Display the cluster stats summary",
	Long:  `Display the cluster stats summary`,

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		monitoring := cli.NewMonitoring(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		monitoring.GetClusterSummary(cliOptions.Resolve("organizationID", organizationID), args[0], rangeMinutes)
	},
}
