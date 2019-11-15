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
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-unified-logging-go"
)

// DecoratorResponse is a structure returned when applying decorator.
// It is composed of all the fields to which a decorator can be applied
// Keep in mind that the decorator could return an error
// and that the field wondered about could be nil
type DecoratorResponse struct {
	AppDescriptorList *grpc_application_go.AppDescriptorList
	AppInstanceList   *grpc_public_api_go.AppInstanceList
	LogResponse       *grpc_unified_logging_go.LogResponse
	Error             derrors.Error
}

type Decorator interface {
	Apply(elements []interface{}) []interface{}
	Validate() derrors.Error
}

// ApplyDecorator function in which the 'decorator apply' is called depending on the type of the input argument
func ApplyDecorator(result interface{}, decorator Decorator) *DecoratorResponse {
	switch result := result.(type) {
	case *grpc_application_go.AppDescriptorList:
		return FromAppDescriptorList(result.Descriptors, decorator)
	case *grpc_public_api_go.AppInstanceList:
		return FromAppInstanceList(result.Instances, decorator)
	}
	return nil
}

// FromAppInstanceList not implemented yet
func FromAppInstanceList(result []*grpc_public_api_go.AppInstance, decorator Decorator) *DecoratorResponse {
	return &DecoratorResponse{
		Error: derrors.NewUnimplementedError("not implemented yet"),
	}
}

// FromAppDescriptorList applies decorator to a AppDescriptorList
func FromAppDescriptorList(result []*grpc_application_go.AppDescriptor, decorator Decorator) *DecoratorResponse {

	// Validate the decorator
	vErr := decorator.Validate()
	if vErr != nil {
		return &DecoratorResponse{
			Error: vErr,
		}
	}

	// convert to []interface{}
	toGenericList := make([]interface{}, len(result))
	for i, d := range result {
		toGenericList[i] = *d
	}

	// call to apply
	ordered := decorator.Apply(toGenericList)

	// reconvert to grpc_application_go.AppDescriptor
	orderedResult := make([]*grpc_application_go.AppDescriptor, len(result))
	for i, d := range ordered {
		aux := d.(grpc_application_go.AppDescriptor)
		orderedResult[i] = &aux
	}

	appDescriptorList := &grpc_application_go.AppDescriptorList{
		Descriptors: orderedResult,
	}

	return &DecoratorResponse{
		AppDescriptorList: appDescriptorList,
	}
}
