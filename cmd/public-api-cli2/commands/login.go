/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	tableOutput "github.com/nalej/public-api/internal/app/output"
	"github.com/nalej/public-api/internal/app/cli2"
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

		conn, output := cli2.NewLoginParameters(cliOptions, loginAddress, loginPort, insecure, useTLS, caCertPath, output, labelLength)
		l := cli2.NewLogin(conn, output)

		creds, err := l.Login(email, password)
		if err != nil {
			if debugLevel {
				log.Fatal().Str("trace", err.DebugReport()).Msg("unable to login into the platform")
			}else{
				log.Fatal().Str("trace", err.Error()).Msg("unable to login into the platform")
			}
		}
		claims, err := l.GetPersonalClaims(creds)
		if err != nil {
			if debugLevel{
				log.Fatal().Str("trace", err.DebugReport()).Msg("unable to login into the platform")
			}else{
				log.Fatal().Str("trace", err.Error()).Msg("unable to login into the platform")
			}
		}
		// TODO Update with the new fields in the newer claim
		cliOptions.Set("organizationID", claims.OrganizationID)
		cliOptions.Set("email", claims.UserID)
		expiration := time.Unix(claims.ExpiresAt, 0).String()
		printLoginResult(claims.UserID, claims.RoleName, claims.OrganizationID, expiration)
	},
}

func printLoginResult(email string, role string, organizationID string, expiration string){
	header := []string{"EMAIL", "ROLE", "ORG_ID", "EXPIRES"}
	values := [][]string{[]string{email, role, organizationID, expiration}}
	tableOutput.PrintFromValues(header, values)
}

func init() {
	loginCmd.Flags().StringVar(&loginAddress, "loginAddress", "", "Address (host) of the login endpoint of the Nalej platform")
	loginCmd.Flags().IntVar(&loginPort, "loginPort", 443, "Port of the Login API (gRPC)")
	loginCmd.Flags().MarkHidden("loginPort")
	loginCmd.Flags().StringVar(&email, "email", "", "User email")
	loginCmd.Flags().StringVar(&password, "password", "", "User password")
	loginCmd.MarkFlagRequired("email")
	loginCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(loginCmd)
}
