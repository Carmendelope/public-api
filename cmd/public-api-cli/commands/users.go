/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "Manage user",
	Long:    `Manage user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.PersistentFlags().StringVar(&email, "email", "", "User email")
	usersCmd.AddCommand(userInfoCmd)
	usersCmd.AddCommand(userListCmd)
	usersCmd.AddCommand(deleteUserCmd)
	usersCmd.AddCommand(resetPasswordCmd)
	updateUserCmd.Flags().StringVar(&name, "name", "", "New name for the user")
	updateUserCmd.Flags().StringVar(&roleName, "role", "", "New role for the user")
	usersCmd.AddCommand(updateUserCmd)
}

var userInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get user info",
	Long:  `Get user info`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		u.Info(options.Resolve("organizationID", organizationID), email)
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Long:  `List users`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		u.List(options.Resolve("organizationID", organizationID))
	},
}

var deleteUserCmd = &cobra.Command{
	Use:     "del",
	Aliases: []string{"delete"},
	Short:   "Delete a user",
	Long:    `Delete a user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		u.Delete(options.Resolve("organizationID", organizationID), email)
	},
}

var resetPasswordCmd = &cobra.Command{
	Use:     "reset-password",
	Aliases: []string{"reset"},
	Short:   "Reset the password of a user",
	Long:    `Reset the password of a user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		u.ResetPassword(options.Resolve("organizationID", organizationID), email)
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the info of a user",
	Long:  `Update the info of a user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(options.Resolve("nalejAddress", nalejAddress), options.ResolveAsInt("port", nalejPort))
		u.Update(options.Resolve("organizationID", organizationID), email, name, roleName)
	},
}
