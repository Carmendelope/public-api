/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var rolesCmd = &cobra.Command{
	Use:     "roles",
	Aliases: []string{"rol", "role"},
	Short:   "Manage roles",
	Long:    `Manage roles`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(rolesCmd)

	listRolesCmd.Flags().BoolVar(&internal, "internal", false, "List internal services")
	rolesCmd.AddCommand(listRolesCmd)

	assignRolesCmd.Flags().StringVar(&email, "email", "", "User email")
	assignRolesCmd.Flags().StringVar(&email, "roleID", "", "User Role ID")
	assignRolesCmd.MarkFlagRequired("email")
	assignRolesCmd.MarkFlagRequired("roleID")
	rolesCmd.AddCommand(assignRolesCmd)
}

var listRolesCmd = &cobra.Command{
	Use:   "list",
	Short: "List roles",
	Long:  `List roles`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		r := cli.NewRoles(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		r.List(options.Resolve("organizationID", organizationID), internal)
	},
}

var assignRolesCmd = &cobra.Command{
	Use:   "assign",
	Short: "Assign new role",
	Long:  `Assign new role`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		r := cli.NewRoles(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		r.Assign(options.Resolve("organizationID", organizationID), email, roleID)
	},
}

