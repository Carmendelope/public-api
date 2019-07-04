/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"fmt"
	"github.com/nalej/public-api/internal/app/cli"
	"github.com/spf13/cobra"
)

var key string
var value string

var optionsCmd = &cobra.Command{
	Use:   "option",
	Aliases: []string{"options", "opt"},
	Short: "Manage default options",
	Long:  `Manage default values for the commands parameters`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	optionsCmd.PersistentFlags().StringVar(&key, "key", "", "Specify the key")
	optionsCmd.PersistentFlags().StringVar(&value, "value", "", "Specify the value")
	rootCmd.AddCommand(optionsCmd)
	optionsCmd.AddCommand(setOptionCmd)
	optionsCmd.AddCommand(getOptionCmd)
	optionsCmd.AddCommand(deleteOptionCmd)
	optionsCmd.AddCommand(listOptionsCmd)
}

var setOptionCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the value for a given key",
	Long:  `Set the value for a given key`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		opts := cli.NewOptions()
		opts.Set(key, value)
	},
}

var getOptionCmd = &cobra.Command{
	Use:   "info",
	Aliases: []string{"get"},
	Short: "Get the value for a given key",
	Long:  `Get the value for a given key`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		opts := cli.NewOptions()
		retrieved := opts.Get(key)
		fmt.Printf("Key: %s Value: %s\n", key, retrieved)
	},
}

var listOptionsCmd = &cobra.Command{
	Use:   "list",
	Aliases: []string{"ls"},
	Short: "List the stored values",
	Long:  `List the stored values`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		opts := cli.NewOptions()
		opts.List()
	},
}

var deleteOptionCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "del", "rm"},
	Short:   "Delete a given key",
	Long:    `Delete a given key`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		opts := cli.NewOptions()
		opts.Delete(key)
		fmt.Printf("Key: %s has been deleted\n", key)
	},
}
