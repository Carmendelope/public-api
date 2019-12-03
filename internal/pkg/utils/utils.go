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

package utils

import (
	"github.com/rs/zerolog/log"
	"reflect"
	"strings"
)

func tagMatches(fieldName string, jsonTag string) bool {
	for _, s := range strings.Split(jsonTag, ",") {
		if s == fieldName {
			return true
		}
	}
	return false
}

// GetFieldName returns the name of a field from the name it has in the json structure
func GetFieldName(jsonName string, element interface{}) string {

	e0 := reflect.ValueOf(element)
	targetField := ""
	for i := 0; i < e0.NumField(); i++ {
		typeField := e0.Type().Field(i)
		value, exists := typeField.Tag.Lookup("json")
		if exists && tagMatches(jsonName, value) {
			targetField = typeField.Name
		}
	}
	log.Debug().Str("TargetField", targetField).Str("jsonName", jsonName).Msg("name retrieved")

	return targetField
}
