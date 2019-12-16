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
	searchCmd.Flags().StringVar(&descriptorID, "descriptorID", "", "Application descriptor identifier")
	searchCmd.Flags().StringVar(&instanceID, "instanceID", "", "Application instance identifier")
	searchCmd.Flags().StringVar(&sgID, "sgID", "", "Service group identifier")
	searchCmd.Flags().StringVar(&sgInstanceID, "sgInstanceID", "", "Service group instance identifier")
	searchCmd.Flags().StringVar(&serviceID, "serviceID", "", "Service identifier")
	searchCmd.Flags().StringVar(&serviceInstanceID, "serviceInstanceID", "", "Service instance identifier")
	searchCmd.Flags().StringVar(&from, "from", "", "Start time of logs")
	searchCmd.Flags().StringVar(&to, "to", "", "End time of logs")
	searchCmd.Flags().BoolVar(&desc, "desc", false, "Sort results in descending time order")
	searchCmd.Flags().BoolVar(&redirectLog, "redirectResultAsLog", false, "Redirect the result to the CLI log")
	searchCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Specify if the logs should be streamed")
	searchCmd.Flags().BoolVar(&nFirst, "nFirst", false, "Specify if the user expects to receive the first n results or not")

	logCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVar(&descriptorID, "descriptorID", "", "Application descriptor identifier")
	downloadCmd.Flags().StringVar(&instanceID, "instanceID", "", "Application instance identifier")
	downloadCmd.Flags().StringVar(&sgID, "sgID", "", "Service group identifier")
	downloadCmd.Flags().StringVar(&sgInstanceID, "sgInstanceID", "", "Service group instance identifier")
	downloadCmd.Flags().StringVar(&serviceID, "serviceID", "", "Service identifier")
	downloadCmd.Flags().StringVar(&serviceInstanceID, "serviceInstanceID", "", "Service instance identifier")
	downloadCmd.Flags().StringVar(&from, "from", "", "Start time of logs")
	downloadCmd.Flags().StringVar(&to, "to", "", "End time of logs")
	downloadCmd.Flags().BoolVar(&desc, "desc", false, "Sort results in descending time order")
	downloadCmd.Flags().BoolVar(&metadata, "metadata", false, "Specify if the user expects to receive log entries metadata")
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

		l := cli.NewUnifiedLogging(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		l.Search(cliOptions.Resolve("organizationID", organizationID), descriptorID, instanceID, sgID, sgInstanceID, serviceID, serviceInstanceID, message, from, to, desc, redirectLog, follow, nFirst)

	},
}

var downloadCmd = &cobra.Command{
	Use:   "download [filter string]",
	Short: "Download application logs",
	Long:  `Download application logs based on application and service group instance`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()

		// Message filter argument
		if len(args) > 0 {
			message = args[0]
		}

		l := cli.NewUnifiedLogging(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		l.Download(cliOptions.Resolve("organizationID", organizationID), descriptorID, instanceID, sgID, sgInstanceID, serviceID, serviceInstanceID, message, from, to, desc, metadata)

	},
}