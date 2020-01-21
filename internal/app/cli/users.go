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

package cli

import (
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Users struct {
	Connection
	Credentials
}

func NewUsers(address string, port int, insecure bool, useTLS bool, caCertPath string, output string, labelLength int) *Users {
	return &Users{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (u *Users) load() {
	err := u.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (u *Users) getClient() (grpc_public_api_go.UsersClient, *grpc.ClientConn) {
	conn, err := u.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewUsersClient(conn)
	return client, conn
}

// Add a new user to the organization.
func (u *Users) Add(organizationID string, email string, password string, name string, roleName string, photoPath string, lastName string, location string, phone string, title string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if email == "" {
		log.Fatal().Msg("email cannot be empty")
	}

	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	photoBase64 := ""
	if photoPath != "" {
		tempPhotoBase64, err := PhotoPathToBase64(photoPath)
		if err != nil {
			log.Error().Err(err).Msg("error converting image to base64")
			u.ExitOnError(err, "error trying to add image")
		}
		photoBase64 = tempPhotoBase64
	}

	addRequest := &grpc_public_api_go.AddUserRequest{
		OrganizationId: organizationID,
		Email:          email,
		Password:       password,
		Name:           name,
		PhotoBase64:    photoBase64,
		LastName:       lastName,
		Location:       location,
		Phone:          phone,
		Title:          title,
		RoleName:       roleName,
	}
	log.Debug().Interface("add user request", addRequest).Msg("debugging")
	added, err := client.Add(ctx, addRequest)
	u.PrintResultOrError(added, err, "cannot add user")
}

// Info retrieves the information of a user.
func (u *Users) Info(organizationID string, email string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	userID := &grpc_user_go.UserId{
		OrganizationId: organizationID,
		Email:          email,
	}
	info, err := client.Info(ctx, userID)
	u.PrintResultOrError(info, err, "cannot obtain user info")
}

// List the users of an organization.
func (u *Users) List(organizationID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	userID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}
	users, err := client.List(ctx, userID)
	u.PrintResultOrError(users, err, "cannot obtain user list")
}

// Delete a user from an organization.
func (u *Users) Delete(organizationID string, email string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	userID := &grpc_user_go.UserId{
		OrganizationId: organizationID,
		Email:          email,
	}
	done, err := client.Delete(ctx, userID)
	u.PrintResultOrError(done, err, "cannot delete user")
}

// Reset the password of a user.
func (u *Users) ChangePassword(organizationID string, email string, password string, newPassword string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	passwordRequest := &grpc_user_manager_go.ChangePasswordRequest{
		OrganizationId: organizationID,
		Email:          email,
		Password:       password,
		NewPassword:    newPassword,
	}
	done, err := client.ChangePassword(ctx, passwordRequest)
	u.PrintResultOrError(done, err, "cannot change password")
}

// Update the user information.
func (u *Users) Update(organizationID string, email string, updateName bool, newName string, updatePhoto bool, newPhotoPath string, updateLastName bool, newLastName string, updateTitle bool, newTitle string, updatePhone bool, newPhone string, updateLocation bool, newLocation string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	updateRequest, err := ApplyUpdate(organizationID, email, updateName, newName, updatePhoto, newPhotoPath, updateLastName, newLastName, updateTitle, newTitle, updatePhone, newPhone, updateLocation, newLocation)
	log.Debug().Interface("updateRequest", updateRequest).Msg("sending update request")
	done, err := client.Update(ctx, updateRequest)
	u.PrintResultOrError(done, err, "cannot update user")
}

func ApplyUpdate(organizationID string, email string, updateName bool, newName string, updatePhoto bool, newPhotoPath string, updateLastName bool, newLastName string, updateTitle bool, newTitle string, updatePhone bool, newPhone string, updateLocation bool, newLocation string) (*grpc_user_go.UpdateUserRequest, error) {
	updateRequest := &grpc_user_go.UpdateUserRequest{
		OrganizationId: organizationID,
		Email:          email,
	}
	if updateName {
		updateRequest.UpdateName = true
		updateRequest.Name = newName
	}

	if updatePhoto {
		updateRequest.UpdatePhotoBase64 = true
		photoBase64, err := PhotoPathToBase64(newPhotoPath)
		if err != nil {
			log.Error().Err(err).Msg("error converting image to base64")
			return nil, err
		}
		updateRequest.PhotoBase64 = photoBase64

	}

	if updateLastName {
		updateRequest.UpdateLastName = true
		updateRequest.LastName = newLastName
	}

	if updateTitle {
		updateRequest.UpdateTitle = true
		updateRequest.Title = newTitle
	}

	if updatePhone {
		updateRequest.UpdatePhone = true
		updateRequest.Phone = newPhone
	}

	if updateLocation {
		updateRequest.UpdateLocation = true
		updateRequest.Location = newLocation
	}

	return updateRequest, nil
}
