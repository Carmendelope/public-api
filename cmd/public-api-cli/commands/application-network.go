package commands

import (
	"fmt"
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
	Use:   "add [sourceInstanceID] [targetInstanceID] [outboundName] [inboundName]",
	Short: "Add a new connection",
	Long:  `Add a new connection`,

	Args: cobra.MaximumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		 appNet := cli.NewApplicationNetwork(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"sourceInstanceID", "targetInstanceID", "outboundName", "inboundName"}, args,
			[]string{sourceInstanceID, targetInstanceID, outbound, inbound})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			appNet.AddConnection(cliOptions.Resolve("organizationID", organizationID), targetValues[0],
				targetValues[1], targetValues[2], targetValues[3])
		}
	},
}

var removeConnectionCmd = &cobra.Command{
	Use:   "delete [sourceInstanceID] [targetInstanceID] [outboundName] [inboundName]",
	Short: "remove a connection",
	Long:  `Remove a connection`,
	Aliases: []string{"remove", "del", "rm"},
	Args: cobra.MaximumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		appNet := cli.NewApplicationNetwork(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))

		targetValues, err := ResolveArgument([]string{"sourceInstanceID", "targetInstanceID", "outboundName", "inboundName"}, args,
			[]string{sourceInstanceID, targetInstanceID, outbound, inbound})
		if err != nil {
			fmt.Println(err.Error())
			cmd.Help()
		}else{
			appNet.RemoveConnection(cliOptions.Resolve("organizationID", organizationID), targetValues[0],
				targetValues[1], targetValues[2], targetValues[3], force)
		}
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

