/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package commands

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Show command tree",
	Long:  `Show a tree with the existing commands in public-api-cli`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		log.Info().Msg("Showing command tree!")
		generateTree()
	},
}

func init() {
	rootCmd.AddCommand(treeCmd)
}

func generateTree() {
	result := ""
	for i, level0 := range rootCmd.Commands() {
		aux := recursiveTreeGenerator(level0,1, (i+1)==len(rootCmd.Commands()))
		result = result + "\n" + aux
	}
	fmt.Print(result+"\n")
}

func recursiveTreeGenerator(c *cobra.Command, level int, last bool) string {
	prefix := ""
	numLevels := level


	for i:=1;i<numLevels;i++{
		prefix = prefix + "│    "
	}

	var result = ""
	var flags = ""

	cmdName := c.Use
	if len(c.Aliases) > 0 {
		cmdName = cmdName + " (" + strings.Join(c.Aliases, ", ") + ")"
	}

	if last {
		if c.HasFlags() {
			c.Flags().VisitAll(func(f *pflag.Flag) {
				flags = flags + "--" + f.Name + " "
			})

			if flags != "" {
				result = "└─── " + cmdName + " " + flags
			}
		} else {
			result = "└─── " + cmdName
		}
	} else {
		if c.HasFlags() {
			c.Flags().VisitAll(func(f *pflag.Flag) {
				flags = flags + "--" + f.Name + " "
			})

			if flags != "" {
				result = "├─── " + cmdName + " " + flags
			}
		}else {
			result = "├─── " + cmdName
		}
	}

	for i, subCommand := range c.Commands() {
		aux := recursiveTreeGenerator(subCommand, level + 1, (i+1)==len(c.Commands()))
		result = result + "\n" + aux
	}

	return prefix + result
}
