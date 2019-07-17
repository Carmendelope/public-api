/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package main

import (
	"github.com/nalej/public-api/cmd/public-api-cli/commands"
	"github.com/nalej/public-api/version"
)

var MainVersion string

var MainCommit string

func main() {
	version.AppVersion = MainVersion
	version.Commit = MainCommit
	commands.Execute()
}
