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
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(userListCmd)
	usersCmd.AddCommand(deleteUserCmd)
	usersCmd.AddCommand(resetPasswordCmd)

	usersCmd.AddCommand(userInfoCmd)
	userInfoCmd.Flags().StringVar(&email, "email", "", "User email")

	updateUserCmd.Flags().StringVar(&name, "name", "", "New name for the user")
	updateUserCmd.Flags().StringVar(&photoPath, "photoPath", "", "Path to the new user photo")
	updateUserCmd.Flags().StringVar(&lastName, "lastName", "", "New last name for the user")
	updateUserCmd.Flags().StringVar(&title, "title", "", "New title for the user")
	updateUserCmd.Flags().StringVar(&phone, "phone", "", "New phone for the user")
	updateUserCmd.Flags().StringVar(&location, "location", "", "New location for the user")
	usersCmd.AddCommand(updateUserCmd)

	addUserCmd.Flags().StringVar(&photoPath, "photoPath", "", "Path to user photo")
	addUserCmd.Flags().StringVar(&phone, "phone", "", "Phone")
	addUserCmd.Flags().StringVar(&location, "location", "", "Location")
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
	Use:     "delete <email>",
	Aliases: []string{"remove", "del", "rm"},
	Args:    cobra.ExactArgs(1),
	Short:   "Delete a user",
	Long:    `Delete a user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		u.Delete(cliOptions.Resolve("organizationID", organizationID), args[0])
	},
}

var resetPasswordCmd = &cobra.Command{
	Use:     "reset-password <email> <newPassword>",
	Aliases: []string{"reset"},
	Args:    cobra.ExactArgs(3),
	Short:   "Reset the password of a user",
	Long:    `Reset the password of a user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		u.ChangePassword(cliOptions.Resolve("organizationID", organizationID), args[0], args[1])
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "update <email>",
	Args:  cobra.ExactArgs(1),
	Short: "Update the info of a user",
	Long:  `Update the info of a user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath),
			cliOptions.Resolve("output", output),
			cliOptions.ResolveAsInt("labelLength", labelLength))
		u.Update(cliOptions.Resolve("organizationID", organizationID), args[0],
			cmd.Flag("name").Changed, name,
			cmd.Flag("photoPath").Changed, photoPath,
			cmd.Flag("lastName").Changed, lastName,
			cmd.Flag("title").Changed, title,
			cmd.Flag("phone").Changed, phone,
			cmd.Flag("location").Changed, location)
	},
}

var addUserCmd = &cobra.Command{
	Use:   "add <email> <password> <name> <lastName> <roleName> <title>",
	Args:  cobra.ExactArgs(6),
	Short: "Add a new user",
	Long:  `Add a new user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		u := cli.NewUsers(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		u.Add(cliOptions.Resolve("organizationID", organizationID), args[0], args[1], args[2], args[4], photoPath, args[3], location, phone, args[5])
	},
}
