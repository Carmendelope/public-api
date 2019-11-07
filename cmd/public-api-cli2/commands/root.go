/*
 * Copyright 2019 Nalej
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
 *
 */

package commands

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/public-api/internal/app/cli2"
	"github.com/nalej/public-api/internal/app/options"
	"github.com/nalej/public-api/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var debugLevel bool
var consoleLogging bool

var loginAddress string
var nalejAddress string
var nalejPort int

var insecure bool
var useTLS bool
var caCertPath string
var output string
var labelLength int

var cliOptions options.Options

var rootCmd = &cobra.Command{
	Use:     "public-api-cli2",
	Short:   "CLI for the new version of public-api",
	Long:    `A command line client for the public-api with improved entities`,
	Version: "unknown-version",
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debugLevel, "debug", false, "Set debug level")
	rootCmd.PersistentFlags().BoolVar(&consoleLogging, "consoleLogging", false, "Pretty print logging")
	rootCmd.PersistentFlags().StringVar(&nalejAddress, cli2.NalejAddress, "", "Address (host) of the Nalej platform")
	rootCmd.PersistentFlags().IntVar(&nalejPort, "port", 443, "Port of the Nalej platform Public API")
	rootCmd.PersistentFlags().MarkHidden("port")
	rootCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "Skip CA validation when connecting to a secure TLS server")
	rootCmd.PersistentFlags().BoolVar(&useTLS, "useTLS", true, "Connect to a TLS server")
	rootCmd.PersistentFlags().StringVar(&caCertPath, "cacert", "", "Path of the CA certificate to validate the server connection")
	rootCmd.PersistentFlags().StringVar(&output, cli2.OutputFormat, "table", "Output format: JSON or table (default)")
	rootCmd.PersistentFlags().MarkHidden(cli2.OutputFormat)
	rootCmd.PersistentFlags().IntVar(&labelLength, cli2.OutputLabelLength, 0, "Maximum labels length")
	rootCmd.PersistentFlags().MarkHidden(cli2.OutputLabelLength)
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

func ResolveArgument(attributeName []string, args []string, flagValue []string) ([]string, derrors.Error) {
	result := make([]string, 0)

	if len(args) < len(attributeName) {
		for index := 0; index < len(attributeName); index++ {
			if flagValue[index] == "" {
				return nil, derrors.NewNotFoundError(fmt.Sprintf("argument %s or flag value --%s not found", attributeName[index], attributeName[index]))
			}
		}
		return flagValue, nil
	}

	if len(attributeName) != len(flagValue) {
		return nil, derrors.NewInternalError("length mismatch")
	}

	for index := 0; index < len(attributeName); index++ {
		found := false
		if flagValue[index] != "" {
			result = append(result, flagValue[index])
			found = true
		}
		if args[index] != "" {
			result = append(result, args[index])
			found = true
		}
		if !found {
			return nil, derrors.NewNotFoundError(attributeName[index])
		}
	}

	return result, nil
}
