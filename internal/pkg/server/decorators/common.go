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
	"github.com/nalej/grpc-public-api-go"
)

// DecoratorResponse is a structure returned when applying decorator.
// It is composed of all the fields to which a decorator can be applied
// Keep in mind that the decorator could return an error
// and that the field wondered about could be nil
type DecoratorResponse struct {
	AppDescriptorList []*grpc_application_go.AppDescriptor
	AppInstanceList   []*grpc_public_api_go.AppInstance
	LogResponseList   []*grpc_application_manager_go.LogEntryResponse
	Error             derrors.Error
}

// ApplyDecorator function in which the 'decorator apply' is called depending on the type of the input argument
// if more structures are added to which a decorator can be applied, the switch will have to be extended to include them
func ApplyDecorator(result interface{}, decorator Decorator) *DecoratorResponse {

	// validate if the decorator can be applied
	vErr := decorator.Validate(result)
	if vErr != nil {
		return &DecoratorResponse{
			Error: vErr,
		}
	}

	switch result := result.(type) {
	case []*grpc_application_go.AppDescriptor:
		return FromAppDescriptorList(result, decorator)
	case []*grpc_public_api_go.AppInstance:
		return FromAppInstanceList(result, decorator)
	case []*grpc_application_manager_go.LogEntryResponse:
		return FromLogEntryResponse(result, decorator)
	}
	return &DecoratorResponse{
		Error: derrors.NewInvalidArgumentError("unable to apply decorator"),
	}
}

// FromAppInstanceList not implemented yet
func FromAppInstanceList(result []*grpc_public_api_go.AppInstance, decorator Decorator) *DecoratorResponse {
	return &DecoratorResponse{
		Error: derrors.NewUnimplementedError("not implemented yet"),
	}
}

// FromAppDescriptorList applies decorator to a AppDescriptorList
func FromAppDescriptorList(result []*grpc_application_go.AppDescriptor, decorator Decorator) *DecoratorResponse {
	// convert to []interface{}
	toGenericList := make([]interface{}, len(result))
	for i, d := range result {
		toGenericList[i] = *d
	}

	// call to apply
	ordered, err := decorator.Apply(toGenericList)
	if err != nil {
		return &DecoratorResponse{
			Error: err,
		}
	}

	// reconvert to grpc_application_go.AppDescriptor
	orderedResult := make([]*grpc_application_go.AppDescriptor, len(result))
	for i, d := range ordered {
		aux := d.(grpc_application_go.AppDescriptor)
		orderedResult[i] = &aux
	}

	return &DecoratorResponse{
		AppDescriptorList: orderedResult,
	}
}

// FromLogEntryResponse applies decorator to a FromLogEntryResponse
func FromLogEntryResponse(result []*grpc_application_manager_go.LogEntryResponse, decorator Decorator) *DecoratorResponse {
	// convert to []interface{}
	toGenericList := make([]interface{}, len(result))
	for i, d := range result {
		toGenericList[i] = *d
	}

	// call to apply
	ordered, err := decorator.Apply(toGenericList)
	if err != nil {
		return &DecoratorResponse{
			Error: err,
		}
	}

	// reconvert to grpc_public_api_go.LogEntryResponse
	orderedResult := make([]*grpc_application_manager_go.LogEntryResponse, len(result))
	for i, d := range ordered {
		aux := d.(grpc_application_manager_go.LogEntryResponse)
		orderedResult[i] = &aux
	}

	return &DecoratorResponse{
		LogResponseList: orderedResult,
	}
}
