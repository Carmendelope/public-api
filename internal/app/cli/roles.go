package cli

import (
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Roles struct {
	Connection
	Credentials
}

func NewRoles(address string, port int) *Roles {
	return &Roles{
		Connection:  *NewConnection(address, port),
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
	if internal{
		roles, err = client.ListInternal(ctx, orgID)
	}else{
		roles, err = client.List(ctx, orgID)
	}
	r.PrintResultOrError(roles, err, "cannot obtain role list")
}
