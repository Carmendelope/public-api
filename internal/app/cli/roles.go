package cli

import (
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Roles struct {
	Connection
	Credentials
}

func NewRoles(address string, port int, insecure bool, caCertPath string, output string) *Roles {
	return &Roles{
		Connection:  *NewConnection(address, port, insecure, caCertPath, output),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (r *Roles) load() {
	err := r.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (r *Roles) getClient() (grpc_public_api_go.RolesClient, *grpc.ClientConn) {
	conn, err := r.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewRolesClient(conn)
	return client, conn
}

func (r *Roles) List(organizationID string, internal bool) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	r.load()
	ctx, cancel := r.GetContext()
	client, conn := r.getClient()
	defer conn.Close()
	defer cancel()
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}
	var roles *grpc_public_api_go.RoleList
	var err error
	if internal {
		roles, err = client.ListInternal(ctx, orgID)
	} else {
		roles, err = client.List(ctx, orgID)
	}
	r.PrintResultOrError(roles, err, "cannot obtain role list")
}

func (r *Roles) Assign(organizationID string, email string, roleID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	r.load()
	ctx, cancel := r.GetContext()
	client, conn := r.getClient()
	defer conn.Close()
	defer cancel()
	request := &grpc_user_manager_go.AssignRoleRequest{
		OrganizationId: organizationID,
		Email:          email,
		RoleId:         roleID,
	}

	user, err := client.AssignRole(ctx, request)

	r.PrintResultOrError(user, err, "cannot obtain assign the new role")
}
