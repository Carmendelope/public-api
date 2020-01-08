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
	"github.com/nalej/public-api/internal/app/options"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"time"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into the Nalej platform",
	Long:  `Login into the Nalej platform`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()

		targetAddress := cliOptions.Resolve("loginAddress", loginAddress)
		if targetAddress == "" {
			log.Fatal().Msg("loginAddress is required")
		}

		l := cli.NewLogin(
			targetAddress,
			loginPort,
			insecure, useTLS,
			cliOptions.Resolve("cacert", caCertPath),
			cliOptions.Resolve("output", output),
			cliOptions.ResolveAsInt("labelLength", labelLength))
		creds, err := l.Login(email, password)
		if err != nil {
			if debugLevel {
				log.Fatal().Str("trace", err.DebugReport()).Msg("unable to login into the platform")
			} else {
				log.Fatal().Str("trace", err.Error()).Msg("unable to login into the platform")
			}
		}
		claims, err := l.GetPersonalClaims(creds)
		if err != nil {
			if debugLevel {
				log.Fatal().Str("trace", err.DebugReport()).Msg("unable to login into the platform")
			} else {
				log.Fatal().Str("trace", err.Error()).Msg("unable to login into the platform")
			}
		}
		opts := options.NewOptions()
		opts.Set("organizationID", claims.OrganizationID)
		opts.Set("email", claims.UserID)
		expiration := time.Unix(claims.ExpiresAt, 0).String()
		printLoginResult(claims.UserID, claims.RoleName, claims.OrganizationID, expiration)
	},
}

func printLoginResult(email string, role string, organizationID string, expiration string) {
	header := []string{"EMAIL", "ROLE", "ORG_ID", "EXPIRES"}
	values := [][]string{[]string{email, role, organizationID, expiration}}
	cli.PrintFromValues(header, values)
}

func init() {
	loginCmd.Flags().StringVar(&loginAddress, "loginAddress", "", "Address (host) of the login endpoint of the Nalej platform")
	loginCmd.Flags().IntVar(&loginPort, "loginPort", 443, "Port of the Login API (gRPC)")
	_ = loginCmd.Flags().MarkHidden("loginPort")
	loginCmd.Flags().StringVar(&email, "email", "", "User email")
	loginCmd.Flags().StringVar(&password, "password", "", "User password")
	_ = loginCmd.MarkFlagRequired("email")
	_ = loginCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(loginCmd)
}
