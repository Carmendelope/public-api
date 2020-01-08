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

package cli2

import (
	"github.com/nalej/public-api/internal/app/options"
	"github.com/nalej/public-api/internal/app/output"
)

// NewLoginParameters creates a connection pointing to the login API.
func NewLoginParameters(options options.Options, loginAddress string, loginPort int, insecure bool, useTLS bool, caCertPath, outputFormat string, labelLength int) (*Connection, *output.Output) {
	conn := NewConnection(
		options.Resolve(LoginAddress, loginAddress),
		loginPort,
		insecure, useTLS,
		options.Resolve(CACert, caCertPath))
	output := output.NewOutput(
		options.Resolve(OutputFormat, outputFormat),
		options.ResolveAsInt(OutputLabelLength, labelLength))
	return conn, output
}

// NewCommandParameters creates the connection and output structures based on the options set
// by the user at system level, and the value of the provided flags.
func NewCommandParameters(options options.Options, nalejAddress string, port int, insecure bool, useTLS bool, caCertPath, outputFormat string, labelLength int) (*Connection, *output.Output) {
	conn := NewConnection(
		options.Resolve(NalejAddress, nalejAddress),
		port,
		insecure, useTLS,
		options.Resolve(CACert, caCertPath))
	output := output.NewOutput(
		options.Resolve(OutputFormat, outputFormat),
		options.ResolveAsInt(OutputLabelLength, labelLength))
	return conn, output
}
