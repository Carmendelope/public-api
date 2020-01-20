/*
 * Copyright 2020 Nalej
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
	_ = userInfoCmd.MarkPersistentFlagRequired("email")
	usersCmd.AddCommand(userListCmd)
	usersCmd.AddCommand(deleteUserCmd)

	resetPasswordCmd.Flags().StringVar(&password, "password", "", "Password")
	resetPasswordCmd.Flags().StringVar(&newPassword, "newPassword", "", "New password")
	_ = resetPasswordCmd.MarkPersistentFlagRequired("email")
	_ = resetPasswordCmd.MarkFlagRequired("password")
	_ = resetPasswordCmd.MarkFlagRequired("newPassword")
	usersCmd.AddCommand(resetPasswordCmd)

	updateUserCmd.Flags().StringVar(&name, "name", "", "New name for the user")
	updateUserCmd.Flags().StringVar(&photoPath, "photoPath", "", "Path to the new user photo")
	updateUserCmd.Flags().StringVar(&lastName, "lastName", "", "New last name for the user")
	updateUserCmd.Flags().StringVar(&title, "title", "", "New title for the user")
	updateUserCmd.Flags().StringVar(&phone, "phone", "", "New phone for the user")
	updateUserCmd.Flags().StringVar(&location, "location", "", "New location for the user")
	_ = updateUserCmd.MarkFlagRequired("email")
	usersCmd.AddCommand(updateUserCmd)

	addUserCmd.Flags().StringVar(&name, "name", "", "Full name")
	addUserCmd.Flags().StringVar(&roleName, "role", "", "Role name")
	addUserCmd.Flags().StringVar(&password, "password", "", "Password")
	addUserCmd.Flags().StringVar(&photoPath, "photoPath", "", "Path to user photo")
	addUserCmd.Flags().StringVar(&lastName, "lastName", "", "Last name")
	addUserCmd.Flags().StringVar(&title, "title", "", "Title")
	addUserCmd.Flags().StringVar(&phone, "phone", "", "Phone")
	addUserCmd.Flags().StringVar(&location, "location", "", "Location")
	_ = addUserCmd.MarkPersistentFlagRequired("email")
	_ = addUserCmd.MarkFlagRequired("name")
	_ = addUserCmd.MarkFlagRequired("role")
	_ = addUserCmd.MarkFlagRequired("password")
	usersCmd.AddCommand(addUserCmd)
}

var userInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"get"},
	Short:   "Get user info",
	Long:    `Get user info`,
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
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List users",
	Long:    `List users`,
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
		u.Update(cliOptions.Resolve("organizationID", organizationID), email, name, photoPath, lastName, title, phone, location)

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
		u.Add(cliOptions.Resolve("organizationID", organizationID), email, password, name, roleName, photoPath, lastName, location, phone, title)
	},
}
