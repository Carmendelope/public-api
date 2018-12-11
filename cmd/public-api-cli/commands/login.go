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
		_, err := l.Login(email, password)
		if err != nil {
			log.Error().Str("trace", err.DebugReport()).Msg("unable to login into the platform")
		}else{
			log.Info().Msg("Successfully logged into the platform")
		}
	},
}

func init() {
	loginCmd.Flags().IntVar(&loginPort, "loginPort", 8444, "Port of the Login API (gRPC)")
	loginCmd.Flags().StringVar(&email, "email", "", "User email")
	loginCmd.Flags().StringVar(&password, "password", "", "User password")
	rootCmd.AddCommand(loginCmd)
}
