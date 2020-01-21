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

package output

import (
	"fmt"
	"github.com/nalej/grpc-common-go"
	"github.com/rs/zerolog/log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

// MinWidth with the minimum column width.
const MinWidth = 5

// TabWidth with the tabulator length.
const TabWidth = 2

// Padding with the length of the padding element.
const Padding = 3

// ResultTable structure containing a table like structure.
type ResultTable struct {
	data [][]string
}

// AsTable obtains the table structure of a given result based on its type.
func AsTable(result interface{}, labelLength int) *ResultTable {
	switch result.(type) {
	case *grpc_common_go.Success:
		return FromSuccess(result.(*grpc_common_go.Success))
	default:
		log.Fatal().Str("type", fmt.Sprintf("%T", result)).Msg("unsupported type when producing table output")
	}
	return nil
}

// Print the given table on stdout
func (t *ResultTable) Print() {
	w := tabwriter.NewWriter(os.Stdout, MinWidth, TabWidth, Padding, ' ', 0)
	for _, d := range t.data {
		toPrint := strings.Join(d, "\t")
		fmt.Fprintln(w, toPrint)
	}
	w.Flush()
}

// PrintFromValues creates a table with a set of values under a header and prints it to the stdout.
func PrintFromValues(header []string, values [][]string) {
	w := tabwriter.NewWriter(os.Stdout, MinWidth, TabWidth, Padding, ' ', 0)
	fmt.Fprintln(w, strings.Join(header, "\t"))
	for _, d := range values {
		toPrint := strings.Join(d, "\t")
		fmt.Fprintln(w, toPrint)
	}
	w.Flush()
}

// TransformLabels transforms a map of labels into a printable entity
func TransformLabels(labels map[string]string, labelLength int) string {
	r := make([]string, 0)

	sortedKeys := GetSortedKeys(labels)
	for _, k := range sortedKeys {
		label := fmt.Sprintf("%s:%s", k, labels[k])
		r = append(r, label)
	}
	labelString := strings.Join(r, ",")
	truncatedR := TruncateString(labelString, labelLength)
	return truncatedR
}

// GetSortedKeys get the map of labels sorted by key name.
func GetSortedKeys(labels map[string]string) []string {
	sortedKeys := make([]string, len(labels))
	i := 0
	for k, _ := range labels {
		sortedKeys[i] = k
		i++
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

// TruncateString truncates a string to a given length depending on the user options.
func TruncateString(text string, length int) string {
	if length <= 0 {
		return text
	}
	truncatedString := text
	if len(text) > length {
		if length > 3 {
			length -= 3
		}
		truncatedString = text[0:length] + "..."
	}
	return truncatedString
}

// ----
// Users
// ----

// ----
// Account
// ----

// ----
// Project
// ----

// ----
// Roles
// ----

// ----
// Common
// ----

func FromSuccess(result *grpc_common_go.Success) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"RESULT"})
	r = append(r, []string{"OK"})
	return &ResultTable{r}
}
