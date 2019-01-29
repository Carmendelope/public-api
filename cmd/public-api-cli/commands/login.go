/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into the Nalej platform",
	Long:  `Login into the Nalej platform`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		l := cli.NewLogin(
			options.Resolve("nalejAddress", nalejAddress),
			loginPort,
			insecure,
			options.Resolve("cacert", caCertPath))
		creds, err := l.Login(email, password)
		if err != nil {
			log.Fatal().Str("trace", err.DebugReport()).Msg("unable to login into the platform")
		}
		log.Info().Msg("Successfully logged into the platform")
		claims, err := l.GetPersonalClaims(creds)
		if err != nil {
			log.Fatal().Str("trace", err.DebugReport()).Msg("unable to login into the platform")
		}
		opts := cli.NewOptions()
		opts.Set("organizationID", claims.OrganizationID)
		opts.Set("email", claims.UserID)
	},
}

func init() {
	loginCmd.Flags().IntVar(&loginPort, "loginPort", 8444, "Port of the Login API (gRPC)")
	loginCmd.Flags().StringVar(&email, "email", "", "User email")
	loginCmd.Flags().StringVar(&password, "password", "", "User password")
	loginCmd.MarkFlagRequired("email")
	loginCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(loginCmd)
}
