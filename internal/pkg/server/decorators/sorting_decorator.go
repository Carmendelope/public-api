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

package decorators

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-organization-manager-go"
	"github.com/nalej/public-api/internal/pkg/utils"
	"github.com/rs/zerolog/log"
	"reflect"
	"sort"
)

// AppDescriptorListAllowedFields: name of the fields by which a list of descriptors can be sorted
// Keep in mind that the names are those defined in the json structure
var AppDescriptorListAllowedFields = []string{"name"}
var LogResponseAllowedFields = []string{"timestamp"}
var SettingsAllowedFields = []string{"key"}

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
func (od *OrderDecorator) Validate(result interface{}) derrors.Error {

	switch result.(type) {
	case []*grpc_application_go.AppDescriptor:
		return od.ValidateSortingDecorator(AppDescriptorListAllowedFields)
	case []*grpc_application_manager_go.LogEntryResponse:
		return od.ValidateSortingDecorator(LogResponseAllowedFields)
	case []*grpc_organization_manager_go.Setting:
		return od.ValidateSortingDecorator(SettingsAllowedFields)
	}

	return derrors.NewInvalidArgumentError("sorting decorator not allowed")

}

func (od *OrderDecorator) ValidateSortingDecorator(allowedFields []string) derrors.Error {
	found := false

	for _, allowed := range allowedFields {
		if allowed == od.Options.Field {
			found = true
		}
	}

	if !found {
		return derrors.NewInvalidArgumentError("unable to apply decorator in field").WithParams(od.Options.Field)
	}

	return nil
}

func (od *OrderDecorator) Apply(elements []interface{}) ([]interface{}, derrors.Error) {

	if len(elements) == 0 {
		return elements, nil
	}

	targetName := utils.GetFieldName(od.Options.Field, elements[0])
	// if targetName is not found
	if targetName == "" {
		return nil, derrors.NewInvalidArgumentError("unable to apply decorator, field not found").WithParams(od.Options.Field)
	}

	sort.SliceStable(elements, func(i, j int) bool {

		e1 := reflect.ValueOf(elements[i])
		e2 := reflect.ValueOf(elements[j])

		switch e1.FieldByName(targetName).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if od.Options.Asc {
				return e1.FieldByName(targetName).Int() < e2.FieldByName(targetName).Int()
			}
			return e1.FieldByName(targetName).Int() > e2.FieldByName(targetName).Int()

		case reflect.String:
			if od.Options.Asc {
				return e1.FieldByName(targetName).String() < e2.FieldByName(targetName).String()
			}
			return e1.FieldByName(targetName).String() > e2.FieldByName(targetName).String()
		}
		log.Warn().Interface("Field", e1.FieldByName(targetName).Kind()).Msg("not supported")
		return false
	})

	return elements, nil
}
