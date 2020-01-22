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

var rolesCmd = &cobra.Command{
	Use:     "role",
	Aliases: []string{"rol", "roles"},
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
	assignRolesCmd.Flags().StringVar(&roleID, "roleID", "", "User Role ID")
	assignRolesCmd.MarkFlagRequired("email")
	assignRolesCmd.MarkFlagRequired("roleID")
	rolesCmd.AddCommand(assignRolesCmd)
}

var listRolesCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List roles",
	Long:    `List roles`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		r := cli.NewRoles(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		r.List(cliOptions.Resolve("organizationID", organizationID), internal)
	},
}

var assignRolesCmd = &cobra.Command{
	Use:   "assign",
	Short: "Assign new role",
	Long:  `Assign new role`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		r := cli.NewRoles(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath), cliOptions.Resolve("output", output), cliOptions.ResolveAsInt("labelLength", labelLength))
		r.Assign(cliOptions.Resolve("organizationID", organizationID), email, roleID)
	},
}
