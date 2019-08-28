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
	userInfoCmd.MarkPersistentFlagRequired("email")
	usersCmd.AddCommand(userListCmd)
	usersCmd.AddCommand(deleteUserCmd)

	resetPasswordCmd.Flags().StringVar(&password, "password", "", "Password")
	resetPasswordCmd.Flags().StringVar(&newPassword, "newPassword", "", "New password")
	resetPasswordCmd.MarkPersistentFlagRequired("email")
	resetPasswordCmd.MarkFlagRequired("password")
	resetPasswordCmd.MarkFlagRequired("newPassword")
	usersCmd.AddCommand(resetPasswordCmd)

	updateUserCmd.Flags().StringVar(&name, "name", "", "New name for the user")
	updateUserCmd.MarkFlagRequired("name")
	usersCmd.AddCommand(updateUserCmd)

	addUserCmd.Flags().StringVar(&name, "name", "", "Full name")
	addUserCmd.Flags().StringVar(&roleName, "role", "", "Rol name")
	addUserCmd.Flags().StringVar(&password, "password", "", "Password")
	addUserCmd.MarkPersistentFlagRequired("email")
	addUserCmd.MarkFlagRequired("name")
	addUserCmd.MarkFlagRequired("role")
	addUserCmd.MarkFlagRequired("password")
	usersCmd.AddCommand(addUserCmd)
}

var userInfoCmd = &cobra.Command{
	Use:   "info",
	Aliases: []string{"get"},
	Short: "Get user info",
	Long:  `Get user info`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		u.Info(cliOptions.Resolve("organizationID", organizationID), cliOptions.Resolve("email", email))
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Aliases: []string{"ls"},
	Short: "List users",
	Long:  `List users`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		u.List(cliOptions.Resolve("organizationID", organizationID))
	},
}

var deleteUserCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "del", "rm"},
	Short:   "Delete a user",
	Long:    `Delete a user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		u.Delete(cliOptions.Resolve("organizationID", organizationID), email)
	},
}

var resetPasswordCmd = &cobra.Command{
	Use:     "reset-password",
	Aliases: []string{"reset"},
	Short:   "Reset the password of a user",
	Long:    `Reset the password of a user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		u.ChangePassword(cliOptions.Resolve("organizationID", organizationID), email, password, newPassword)
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the info of a user",
	Long:  `Update the info of a user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		u.Update(cliOptions.Resolve("organizationID", organizationID), email, name)
	},
}

var addUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new user",
	Long:  `Add a new user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		u.Add(cliOptions.Resolve("organizationID", organizationID), email, password, name, roleName)
	},
}
