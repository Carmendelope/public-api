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
 *
 */

package decorators

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"reflect"
	"sort"
)

var AppDescriptorListAllowedFields = []string{"Name"}

// OrderOptions represents the ordering to be applied
type OrderOptions struct {
	// Field to be ordered by
	Field string
	// ascending or descending order
	Asc bool
}

// OrderDecorator implements Decorator interface
type OrderDecorator struct {
	Options OrderOptions
}

func NewOrderDecorator(options OrderOptions) Decorator {
	orderDecorator := OrderDecorator{options}
	return &orderDecorator
}

// Validate checks if the field is a field to which decorators can be applied
func (od *OrderDecorator) Validate() derrors.Error {

	found := false

	for _, allowed := range AppDescriptorListAllowedFields {
		if allowed == od.Options.Field {
			found = true
		}
	}

	if !found {
		return derrors.NewInvalidArgumentError("unable to apply decorator in field").WithParams(od.Options.Field)
	}

	return nil
}

func (od *OrderDecorator) Apply(elements []interface{}) []interface{} {

	sort.SliceStable(elements, func(i, j int) bool {

		e1 := reflect.ValueOf(elements[i])
		e2 := reflect.ValueOf(elements[j])

		switch e1.FieldByName(od.Options.Field).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if od.Options.Asc {
				return e1.FieldByName(od.Options.Field).Int() < e2.FieldByName(od.Options.Field).Int()
			}
			return e1.FieldByName(od.Options.Field).Int() > e2.FieldByName(od.Options.Field).Int()

		case reflect.String:
			if od.Options.Asc {
				return e1.FieldByName(od.Options.Field).String() < e2.FieldByName(od.Options.Field).String()
			}
			return e1.FieldByName(od.Options.Field).String() > e2.FieldByName(od.Options.Field).String()
		}
		log.Warn().Interface("Field", e1.FieldByName(od.Options.Field).Kind()).Msg("not supported")
		return false
	})

	log.Debug().Interface("elements", elements).Msg("elements returned")

	return elements
}
