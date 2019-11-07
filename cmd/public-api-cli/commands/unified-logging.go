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
 *
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Manage application logs",
	Long:  `Search logs of specific application and service group instances`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.PersistentFlags().StringVar(&organizationID, "organizationID", "", "Organization identifier")
	_ = logCmd.MarkFlagRequired("organizationID")

	logCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVar(&instanceID, "instanceID", "", "Application instance identifier")
	_ = searchCmd.MarkFlagRequired("instanceID")
	searchCmd.Flags().StringVar(&sgInstanceID, "sgInstanceID", "", "Service group instance identifier")
	searchCmd.Flags().StringVar(&from, "from", "", "Start time of logs")
	searchCmd.Flags().StringVar(&to, "to", "", "End time of logs")
	searchCmd.Flags().BoolVar(&asc, "asc", false, "Sort results in ascending time order")
	searchCmd.Flags().BoolVar(&desc, "desc", false, "Sort results in descending time order")
	searchCmd.Flags().BoolVar(&redirectLog, "redirectResultAsLog", false, "Redirect the result to the CLI log")
}

var searchCmd = &cobra.Command{
	Use:   "search [filter string]",
	Short: "Search application logs",
	Long:  `Search application logs based on application and service group instance`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		// Message filter argument
		if len(args) > 0 {
			message = args[0]
		}

		// Verify sort order
		// Default is ascending, so "--asc" is only used to make the command line
		// very explicit. If "--desc" is set, we know "--asc" is not because of
		// this check.
		if asc && desc {
			log.Fatal().Msg("Specify sort order as ascending OR descending - not both")
		}

		l := cli.NewUnifiedLogging(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		l.Search(cliOptions.Resolve("organizationID", organizationID), instanceID, sgInstanceID, message, from, to, desc, redirectLog)
	},
}
