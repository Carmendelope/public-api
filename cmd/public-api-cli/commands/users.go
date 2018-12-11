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

	resetPasswordCmd.Flags().StringVar(&password, "password", "", "Password")
	resetPasswordCmd.Flags().StringVar(&newPassword, "newPassword", "", "New password")
	resetPasswordCmd.MarkPersistentFlagRequired("email")
	resetPasswordCmd.MarkFlagRequired("password")
	resetPasswordCmd.MarkFlagRequired("newPassword")
	usersCmd.AddCommand(resetPasswordCmd)

	updateUserCmd.Flags().StringVar(&name, "name", "", "New name for the user")
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
	Short: "Get user info",
	Long:  `Get user info`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		u.Info(options.Resolve("organizationID", organizationID), email)
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Long:  `List users`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
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
		u := cli.NewUsers(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
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
		u := cli.NewUsers(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		u.ChangePassword(options.Resolve("organizationID", organizationID), email, password, newPassword)
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the info of a user",
	Long:  `Update the info of a user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		u.Update(options.Resolve("organizationID", organizationID), email, name)
	},
}

var addUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new user",
	Long:  `Add a new user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			options.Resolve("nalejAddress", nalejAddress),
			options.ResolveAsInt("port", nalejPort),
			insecure,
			options.Resolve("cacert", caCertPath))
		u.Add(options.Resolve("organizationID", organizationID), email, password, name, roleName)
	},
}
