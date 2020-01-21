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

var orgCmd = &cobra.Command{
	Use:     "organization",
	Aliases: []string{"org"},
	Short:   "Organization related operations",
	Long:    `Organization related operations`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	orgCmd.PersistentFlags().StringVar(&organizationID, "organizationID", "", "Organization identifier")
	rootCmd.AddCommand(orgCmd)

	// organization commands
	orgCmd.AddCommand(infoCmd)

	updateOrgCmd.Flags().StringVar(&name, "name", "", "new organization name")
	updateOrgCmd.Flags().StringVar(&address, "address", "", "new organization address")
	updateOrgCmd.Flags().StringVar(&city, "city", "", "new organization city")
	updateOrgCmd.Flags().StringVar(&state, "state", "", "new state")
	updateOrgCmd.Flags().StringVar(&country, "country", "", "new organization country")
	updateOrgCmd.Flags().StringVar(&zipCode, "zipCode", "", "new zipCode")
	updateOrgCmd.Flags().StringVar(&photoPath, "photoPath", "", "Organization logo path")
	orgCmd.AddCommand(updateOrgCmd)

	// Setting commands
	orgCmd.AddCommand(setCmd)
	// update
	setCmd.AddCommand(updateSetCmd)
	// list
	setCmd.AddCommand(listSetCmd)
	listSetCmd.Flags().BoolVar(&desc, "desc", false, "Sort settings in descending order")

}

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"get"},
	Short:   "Retrieve organization information",
	Long:    `Retrieve organization information`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		o := cli.NewOrganizations(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath),
			cliOptions.Resolve("output", output),
			cliOptions.ResolveAsInt("labelLength", labelLength))
		o.Info(cliOptions.Resolve("organizationID", organizationID))
	},
}

var updateOrgCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the info of an organization",
	Long:  `Update the info of an organization`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		o := cli.NewOrganizations(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath),
			cliOptions.Resolve("output", output),
			cliOptions.ResolveAsInt("labelLength", labelLength))
		o.Update(cliOptions.Resolve("organizationID", organizationID),
			cmd.Flag("name").Changed, name,
			cmd.Flag("address").Changed, address,
			cmd.Flag("city").Changed, city,
			cmd.Flag("state").Changed, state,
			cmd.Flag("country").Changed, country,
			cmd.Flag("zipCode").Changed, zipCode,
			cmd.Flag("photoPath").Changed, photoPath)
	},
}

var setCmd = &cobra.Command{
	Use:     "setting",
	Aliases: []string{"set"},
	Short:   "Settings related operations",
	Long:    `Settings related operations`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

var updateSetCmd = &cobra.Command{
	Use:   "update [key] [value]",
	Short: "Update a setting",
	Long:  `Update the setting value of a organization`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		o := cli.NewOrganizations(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath),
			cliOptions.Resolve("output", output),
			cliOptions.ResolveAsInt("labelLength", labelLength))

		// Key argument
		key = args[0]
		value = args[1]

		o.UpdateSetting(cliOptions.Resolve("organizationID", organizationID), key, value)

	},
}

var listSetCmd = &cobra.Command{
	Use:   "list",
	Short: "List the settings",
	Long:  `List the settings of a organization`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		o := cli.NewOrganizations(
			cliOptions.Resolve("nalejAddress", nalejAddress),
			cliOptions.ResolveAsInt("port", nalejPort),
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath),
			cliOptions.Resolve("output", output),
			cliOptions.ResolveAsInt("labelLength", labelLength))
		o.ListSettings(cliOptions.Resolve("organizationID", organizationID), desc)
	},
}
