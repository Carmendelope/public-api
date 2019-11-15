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
	"github.com/spf13/cobra"
)

var appNetCmd = &cobra.Command{
	Use:     "appnet",
	Aliases: []string{"application-network", "app-net"},
	Short:   "Manage ApplicationNetwork",
	Long:    `Application Network related operations`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(appNetCmd)

	appNetCmd.AddCommand(addConnectionCmd)

	appNetCmd.AddCommand(removeConnectionCmd)
	removeConnectionCmd.Flags().BoolVar(&force, "force", false, "User confirmation, force required outbound connection ")

	appNetCmd.AddCommand(listConnectionCmd)

	appNetCmd.AddCommand(inboundCmd)
	inboundCmd.AddCommand(inboundAvailableCmd)

	appNetCmd.AddCommand(outboundCmd)
	outboundCmd.AddCommand(outboundAvailableCmd)
}

var addConnectionCmd = &cobra.Command{
	Use:   "add [sourceInstanceID] [outboundName] [targetInstanceID] [inboundName]",
	Short: "Add a new connection",
	Long:  `Add a new connection`,

	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		appNet := cli.NewApplicationNetwork(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		appNet.AddConnection(cliOptions.Resolve("organizationID", organizationID), args[0], args[1], args[2], args[3])

	},
}

var removeConnectionCmd = &cobra.Command{
	Use:     "delete [sourceInstanceID] [outboundName] [targetInstanceID] [inboundName]",
	Short:   "remove a connection",
	Long:    `Remove a connection`,
	Aliases: []string{"remove", "del", "rm"},
	Args:    cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		appNet := cli.NewApplicationNetwork(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		appNet.RemoveConnection(cliOptions.Resolve("organizationID", organizationID), args[0], args[1], args[2], args[3], force)

	},
}

var listConnectionCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the connections",
	Long:    `List the connections`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		appNet := cli.NewApplicationNetwork(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		appNet.ListConnection(cliOptions.Resolve("organizationID", organizationID))
	},
}

var inboundCmd = &cobra.Command{
	Use:   "inbound",
	Short: "Manage inbound interfaces",
	Long:  `Inbound Interfaces related operations`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		_ = cmd.Help()
	},
}

var inboundAvailableCmd = &cobra.Command{
	Use:        "available",
	SuggestFor: nil,
	Short:      "List the network available inbound interfaces",
	Long:       `List the nerwork available inbound interfaces in the organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		appNet := cli.NewApplicationNetwork(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		appNet.ListAvailableInbounds(cliOptions.Resolve("organizationID", organizationID))
	},
}

var outboundCmd = &cobra.Command{
	Use:   "outbound",
	Short: "Manage outbound interfaces",
	Long:  `Outbound Interfaces related operations`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		_ = cmd.Help()
	},
}

var outboundAvailableCmd = &cobra.Command{
	Use:        "available",
	SuggestFor: nil,
	Short:      "List the network available outbound interfaces",
	Long:       `List the nerwork available outbound interfaces in the organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		appNet := cli.NewApplicationNetwork(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		appNet.ListAvailableOutbounds(cliOptions.Resolve("organizationID", organizationID))
	},
}
