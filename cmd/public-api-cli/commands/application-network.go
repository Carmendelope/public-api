package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var appNetCmd = &cobra.Command{
	Use:     "appnet",
	Aliases: []string{"application-network"},
	Short:   "Manage ApplicationNetwork",
	Long:    `Application Network related operations`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(appNetCmd)

	appNetCmd.AddCommand(addConnectionCmd)

	appNetCmd.AddCommand(removeConnectionCmd)
	removeConnectionCmd.Flags().BoolVar(&force, "force", false, "User confirmation, force required outbound connection ")

	appNetCmd.AddCommand(listConnectionCmd)

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
	Use:   "delete [sourceInstanceID] [outboundName] [targetInstanceID] [inboundName]",
	Short: "remove a connection",
	Long:  `Remove a connection`,
	Aliases: []string{"remove", "del", "rm"},
	Args: cobra.ExactArgs(4),
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
	Use:   "list",
	Aliases: []string{"ls"},
	Short: "List the connections",
	Long:  `List the connections`,
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

