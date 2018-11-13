package cli

import (
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Users struct {
	Connection
	Credentials
}

func NewUsers(address string, port int) * Users{
	return &Users{
		Connection: *NewConnection(address, port),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}


func (u * Users) load() {
	err := u.LoadCredentials()
	if err != nil{
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (u * Users) getClient() (grpc_public_api_go.UsersClient, *grpc.ClientConn){
	conn, err := u.GetConnection()
	if err != nil{
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewUsersClient(conn)
	return client, conn
}

func (u * Users) Info(organizationID string, email string){
	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	userID := &grpc_user_go.UserId{
		OrganizationId:       organizationID,
		Email:                email,
	}
	info, err := client.Info(ctx, userID)
	u.PrintResultOrError(info, err, "cannot obtain user info")
}

func (u * Users) List(organizationID string) {
	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	userID := &grpc_organization_go.OrganizationId{
		OrganizationId:       organizationID,
	}
	users, err := client.List(ctx, userID)
	u.PrintResultOrError(users, err, "cannot obtain user list")
}

func (u * Users) Delete(organizationID string, email string) {
	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	userID := &grpc_user_go.UserId{
		OrganizationId:       organizationID,
		Email:                email,
	}
	done, err := client.Delete(ctx, userID)
	u.PrintResultOrError(done, err, "cannot delete user")
}

func (u * Users) ResetPassword(organizationID string, email string) {
	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	userID := &grpc_user_go.UserId{
		OrganizationId:       organizationID,
		Email:                email,
	}
	done, err := client.ResetPassword(ctx, userID)
	u.PrintResultOrError(done, err, "cannot change password")
}

func (u * Users) Update(organizationID string, email string, newName string, newRole string) {
	u.load()
	ctx, cancel := u.GetContext()
	client, conn := u.getClient()
	defer conn.Close()
	defer cancel()

	updateRequest := &grpc_user_go.UpdateUserRequest{
		OrganizationId:       organizationID,
		Email:                email,
	}
	if newName != "" {
		updateRequest.Name = newName
	}
	if newRole != "" {
		updateRequest.Role = newRole
	}
	done, err := client.Update(ctx, updateRequest)
	u.PrintResultOrError(done, err, "cannot change password")
}