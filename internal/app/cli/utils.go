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
 */

package cli

import (
	"github.com/rs/zerolog/log"
	"strings"
)

func GetLabels(rawLabels string) map[string]string {
	labels := make(map[string]string, 0)

	split := strings.Split(rawLabels, ";")
	for _, l := range split {
		ls := strings.Split(l, ":")
		if len(ls) != 2 {
			log.Fatal().Str("label", l).Msg("malformed label, expecting key:value")
		}
		labels[ls[0]] = ls[1]
	}
	return labels
}
