/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"encoding/json"
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/nalej/public-api/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"fmt"
	"os"
)

var debugLevel bool
var consoleLogging bool

var nalejAddress string
var nalejPort int

var options cli.Options

var rootCmd = &cobra.Command{
	Use:     "public-api-cli",
	Short:   "CLI for the public-api",
	Long:    `A command line client for the public-api`,
	Version: "unknown-version",
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debugLevel, "debug", false, "Set debug level")
	rootCmd.PersistentFlags().BoolVar(&consoleLogging, "consoleLogging", false, "Pretty print logging")
	rootCmd.PersistentFlags().StringVar(&nalejAddress, "nalejAddress", "", "Address (host) of the Nalej platform")
	rootCmd.PersistentFlags().IntVar(&nalejPort, "port", -1, "Port of the Nalej platform Public API")
}

func Execute() {
	rootCmd.SetVersionTemplate(version.GetVersionInfo())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// SetupLogging sets the debug level and console logging if required.
func SetupLogging() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debugLevel {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if consoleLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func PrintResult(result interface{}) error {
	//Print descriptors
	res, err := json.MarshalIndent(result, "", "  ")
	if err == nil {
		fmt.Println(string(res))
	}
	return err
}
