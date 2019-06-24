/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/derrors"
	grpc_inventory_manager_go "github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

var edgeControllerCmd = &cobra.Command{
	Use:     "edgecontroller",
	Aliases: []string{"ec"},
	Short:   "Manage edge controllers",
	Long:    `Manage edge controllers`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	// Device groups
	rootCmd.AddCommand(edgeControllerCmd)
	edgeControllerCmd.AddCommand(createJoinTokenECCmd)
	createJoinTokenECCmd.Flags().StringVar(&outputPath, "outputPath", "", "Path to store the resulting token")

	// Unlink
	edgeControllerCmd.AddCommand(unlinkECCmd)

	// Geolocation Update
	updateGeoCmd.Flags().StringVar(&geolocation, "geolocation", "", "Edge Controller geolocation")
	edgeControllerCmd.AddCommand(updateGeoCmd)

	installAgentCmd.Flags().StringVar(&password, "password", "", "SSH password")
	installAgentCmd.Flags().StringVar(&publicKeyPath, "publicKeyPath", "", "SSH public key path")

	installAgentCmd.Flags().StringVar(&agentTypeRaw, "agentType", "LINUX_X86", "Agent type: LINUX_X86, LINUX_ARM32, LINUX_ARM64 or WINDOWS")
	edgeControllerCmd.AddCommand(installAgentCmd)
}


var createJoinTokenECCmd = &cobra.Command{
	Use:   "create-join-token to attach new edge controllers to the platform",
	Short: "Create a join token",
	Long:  `Create a join token for being able to attach new edge controllers to the platform`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		ec := cli.NewEdgeController(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
			ec.CreateJoinToken(options.Resolve("organizationID", organizationID), outputPath)
	},
}

var unlinkECCmd = &cobra.Command{
	Use:   "unlink [edgeControllerID]",
	Short: "Unlink an EIC",
	Long:  `Unlink an EIC from the platform`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		ec := cli.NewEdgeController(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		if len(args) > 0{
			edgeControllerID = args[0]
		}
		ec.Unlink(options.Resolve("organizationID", organizationID), edgeControllerID)
	},
}

var updateGeoCmd = &cobra.Command{
	Use:   "location-update [edgeControllerID]",
	Short: "update edge-controller location",
	Long:  `update edge-controller location`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		ec := cli.NewEdgeController(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))
		if len(args) > 0{
			edgeControllerID = args[0]
		}
		ec.UpdateGeolocation(options.Resolve("organizationID", organizationID), edgeControllerID, geolocation)
	},
}

// Get the agent type enum from raw.
//     // Linux agent on x86
//    LINUX_X86 = 0;
//    // Linux agent on ARM 32 bits
//    LINUX_ARM32 = 1;
//    // Linux agent on ARM 64 bits
//    LINUX_ARM64 = 2;
//    // Windows agent
//    WINDOWS = 3;
func getAgentType(agentTypeRaw string) (*grpc_inventory_manager_go.AgentType, derrors.Error){
	types := map[string]grpc_inventory_manager_go.AgentType{
		"linux_x86":grpc_inventory_manager_go.AgentType_LINUX_X86,
		"linux_arm32":grpc_inventory_manager_go.AgentType_LINUX_ARM32,
		"linux_arm64":grpc_inventory_manager_go.AgentType_LINUX_ARM64,
		"windows":grpc_inventory_manager_go.AgentType_WINDOWS,
	}

	agentType, exists := types[strings.ToLower(agentTypeRaw)]
	if !exists{
		return nil, derrors.NewInvalidArgumentError("specified agent type not suppoted")
	}
	return &agentType, nil
}

var installAgentCmd = &cobra.Command{
	Use:   "install-agent [edgeControllerID] [targetHost] [username]",
	Short: "Install an agent on a given host",
	Long:  `Install an agent through an edge controller on a given host`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		ec := cli.NewEdgeController(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			options.Resolve("cacert", caCertPath), options.Resolve("output", output))

			edgeControllerID = args[0]
			targetHost := args[1]
			username = args[2]
			agentType, err := getAgentType(agentTypeRaw)
			if err != nil{
				log.Fatal().Err(err).Msg("invalid agent type")
			}

		ec.InstallAgent(options.Resolve("organizationID", organizationID), edgeControllerID, *agentType, targetHost, username, password, publicKeyPath)
	},
}