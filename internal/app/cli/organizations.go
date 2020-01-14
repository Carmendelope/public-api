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
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-organization-manager-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
)

type Organizations struct {
	Connection
	Credentials
}

func NewOrganizations(address string, port int, insecure bool, useTLS bool, caCertPath string, output string, labelLength int) *Organizations {
	return &Organizations{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (o *Organizations) Info(organizationID string) *grpc_organization_manager_go.Organization {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	err := o.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}

	c, err := o.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	defer c.Close()
	ctx, cancel := o.GetContext()
	defer cancel()

	orgClient := grpc_public_api_go.NewOrganizationsClient(c)
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}
	info, iErr := orgClient.Info(ctx, orgID)

	o.PrintResultOrError(info, iErr, "cannot obtain organization info")
	return info
}

/*
OrganizationId       string   `protobuf:"bytes,1,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	UpdateName           bool     `protobuf:"varint,2,opt,name=update_name,json=updateName,proto3" json:"update_name,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	UpdateFullAddress    bool     `protobuf:"varint,4,opt,name=update_full_address,json=updateFullAddress,proto3" json:"update_full_address,omitempty"`
	FullAddress          string   `protobuf:"bytes,5,opt,name=full_address,json=fullAddress,proto3" json:"full_address,omitempty"`
	UpdateCity           bool     `protobuf:"varint,6,opt,name=update_city,json=updateCity,proto3" json:"update_city,omitempty"`
	City                 string   `protobuf:"bytes,7,opt,name=city,proto3" json:"city,omitempty"`
	UpdateState          bool     `protobuf:"varint,8,opt,name=update_state,json=updateState,proto3" json:"update_state,omitempty"`
	State                string   `protobuf:"bytes,9,opt,name=state,proto3" json:"state,omitempty"`
	UpdateCountry        bool     `protobuf:"varint,10,opt,name=update_country,json=updateCountry,proto3" json:"update_country,omitempty"`
	Country              string   `protobuf:"bytes,11,opt,name=country,proto3" json:"country,omitempty"`
	UpdateZipCode        bool     `protobuf:"varint,12,opt,name=update_zip_code,json=updateZipCode,proto3" json:"update_zip_code,omitempty"`
	ZipCode              string   `protobuf:"bytes,13,opt,name=zip_code,json=zipCode,proto3" json:"zip_code,omitempty"`
	UpdatePhoto          bool     `protobuf:"varint,14,opt,name=update_photo,json=updatePhoto,proto3" json:"update_photo,omitempty"`
	PhotoBase64          string   `protobuf:"bytes,15,opt,name=photo_base64,json=photoBase64,proto3" json:"photo_base64,omitempty"`
*/
func (o *Organizations) Update(organizationID string, updateName bool, name string, updateFullAddress bool, fullAddress string,
	updateCity bool, city string, updateState bool, state string, updateCountry bool, country string, updateZipCode bool, zipCode string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	err := o.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}

	c, err := o.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	defer c.Close()
	ctx, cancel := o.GetContext()
	defer cancel()

	orgClient := grpc_public_api_go.NewOrganizationsClient(c)
	updateRequest := &grpc_organization_go.UpdateOrganizationRequest{
		OrganizationId:    organizationID,
		UpdateName:        updateName,
		Name:              name,
		UpdateFullAddress: updateFullAddress,
		FullAddress:       fullAddress,
		UpdateCity:        updateCity,
		City:              city,
		UpdateState:       updateState,
		State:             state,
		UpdateCountry:     updateCountry,
		Country:           country,
		UpdateZipCode:     updateZipCode,
		ZipCode:           zipCode,
	}
	success, iErr := orgClient.Update(ctx, updateRequest)
	o.PrintResultOrError(success, iErr, "cannot update organization info")
	return
}
