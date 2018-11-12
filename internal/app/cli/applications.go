package cli

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"io/ioutil"
)

type Applications struct {
	Connection
	Credentials
}

func NewApplications(address string, port int) * Applications {
	return &Applications{
		Connection: *NewConnection(address, port),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (a * Applications) load() {
	err := a.LoadCredentials()
	if err != nil{
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (a * Applications) getClient() (grpc_public_api_go.ApplicationsClient, * grpc.ClientConn) {
	conn, err := a.GetConnection()
	if err != nil{
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	appsClient := grpc_public_api_go.NewApplicationsClient(conn)
	return appsClient, conn
}

func (a * Applications) createAddDescriptorRequest(organizationID string, descriptorPath string) (*grpc_application_go.AddAppDescriptorRequest, derrors.Error){

	content, err := ioutil.ReadFile(descriptorPath)
	if err != nil {
		return nil, derrors.AsError(err, "cannot read descriptor")
	}

	addDescriptorRequest := &grpc_application_go.AddAppDescriptorRequest{}
	err = json.Unmarshal(content, &addDescriptorRequest)
	if err != nil {
		return nil, derrors.AsError(err, "cannot unmarshal structure")
	}

	addDescriptorRequest.OrganizationId = organizationID
	for _, s := range addDescriptorRequest.Services {
		s.OrganizationId = organizationID
	}

	return addDescriptorRequest, nil
}

func (a * Applications) getBasicDescriptor() *grpc_application_go.AddAppDescriptorRequest {

	service := &grpc_application_go.Service{
		ServiceId:            "1",
		Name:                 "Simple MySQL service",
		Description:          "A MySQL instance",
		Type:                 grpc_application_go.ServiceType_DOCKER,
		Image:                "mysql:5.6",
		Specs:                &grpc_application_go.DeploySpecs{Replicas: 1},
		Storage:              []*grpc_application_go.Storage{&grpc_application_go.Storage{MountPath: "/tmp",}},
		ExposedPorts:         []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "mysqlport", InternalPort: 3306, ExposedPort: 3306,
		}},
		EnvironmentVariables: map[string]string{"MYSQL_ROOT_PASSWORD":"root"},
		Configs:              []*grpc_application_go.ConfigFile{&grpc_application_go.ConfigFile{MountPath:"/tmp"}},
		Labels:                map[string]string { "app":"simple-app", "component":"mysql"},
	}

	secRule := grpc_application_go.SecurityRule{
		Name:"all open",
		Access: grpc_application_go.PortAccess_PUBLIC,
		RuleId: "001",
		SourcePort: 3306,
		SourceServiceId: "1",
	}

	return &grpc_application_go.AddAppDescriptorRequest{
		Name:                 "Sample application",
		Description:          "This is a basic descriptor of an application",
		Labels:               map[string]string{"app":"simple-app"},
		Rules:                []*grpc_application_go.SecurityRule{&secRule},
		Services:             []*grpc_application_go.Service{service},
	}
}

func (a * Applications) AddDescriptorHelp() {
	toAdd := a.getBasicDescriptor()
	fmt.Println(`To add a new descriptor, write a JSON file with the descriptor and pass that path as
parameter to the add command.`)
	err := a.PrintResult(toAdd)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load sample application descriptor")
	}
}

func (a * Applications) AddDescriptor(organizationID string, descriptorPath string) {
	a.load()
	ctx, cancel := a.GetContext()
	client, conn := a.getClient()
	defer conn.Close()
	defer cancel()

	addDescriptorRequest, aErr := a.createAddDescriptorRequest(organizationID, descriptorPath)
	if aErr != nil {
		log.Fatal().Str("trace", aErr.DebugReport()).Msg("cannot load application descriptor")
	}
	added, err := client.AddAppDescriptor(ctx, addDescriptorRequest)

	if err != nil{
		log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("cannot add a new application descriptor")
	}
	fmt.Println(added.String())
}

func (a * Applications) GetDescriptor(organizationID string, descriptorID string) {
	a.load()
	ctx, cancel := a.GetContext()
	client, conn := a.getClient()
	defer conn.Close()
	defer cancel()

	appDescriptorID := &grpc_application_go.AppDescriptorId{
		OrganizationId:       organizationID,
		AppDescriptorId:      descriptorID,
	}
	descriptor, err := client.GetAppDescriptor(ctx, appDescriptorID)
	if err != nil{
		log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("cannot obtain descriptor")
	}
	fmt.Println(descriptor.String())
}

func (a * Applications) ListDescriptors(organizationID string) {
	a.load()
	ctx, cancel := a.GetContext()
	client, conn := a.getClient()
	defer conn.Close()
	defer cancel()

	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId:       organizationID,
	}
	list, err := client.ListAppDescriptors(ctx, orgID)
	if err != nil{
		log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("cannot obtain descriptor list")
	}
	fmt.Println(list.String())
}