package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/nalej/grpc-utils/pkg/tools"
	"github.com/nalej/public-api/internal/pkg/server/applications"
	"github.com/nalej/public-api/internal/pkg/server/clusters"
	"github.com/nalej/public-api/internal/pkg/server/nodes"
	"github.com/nalej/public-api/internal/pkg/server/organizations"
	"github.com/nalej/public-api/internal/pkg/server/resources"
	"github.com/nalej/public-api/internal/pkg/server/roles"
	"github.com/nalej/public-api/internal/pkg/server/users"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
)

type Service struct {
	Configuration Config
	Server * tools.GenericGRPCServer
}

// NewService creates a new system model service.
func NewService(conf Config) *Service {
	return &Service{
		conf,
		tools.NewGenericGRPCServer(uint32(conf.Port)),
	}
}

type Clients struct {
	orgClient grpc_organization_go.OrganizationsClient
	clusClient grpc_infrastructure_go.ClustersClient
	nodeClient grpc_infrastructure_go.NodesClient
	umClient grpc_user_manager_go.UserManagerClient
	appClient grpc_application_manager_go.ApplicationManagerClient
}

func (s * Service) GetClients() (* Clients, derrors.Error) {
	smConn, err := grpc.Dial(s.Configuration.SystemModelAddress, grpc.WithInsecure())
	if err != nil{
		return nil, derrors.AsError(err, "cannot create connection with the system model")
	}

	umConn, err := grpc.Dial(s.Configuration.UserManagerAddress, grpc.WithInsecure())
	if err != nil{
		return nil, derrors.AsError(err, "cannot create connection with the user manager")
	}

	appConn, err := grpc.Dial(s.Configuration.ApplicationsManagerAddress, grpc.WithInsecure())
	if err != nil{
		return nil, derrors.AsError(err, "cannot create connection with the applications manager")
	}

	oClient := grpc_organization_go.NewOrganizationsClient(smConn)
	cClient := grpc_infrastructure_go.NewClustersClient(smConn)
	nClient := grpc_infrastructure_go.NewNodesClient(smConn)
	umClient := grpc_user_manager_go.NewUserManagerClient(umConn)
	appClient := grpc_application_manager_go.NewApplicationManagerClient(appConn)

	return &Clients{oClient, cClient, nClient, umClient, appClient}, nil
}

// Run the service, launch the REST service handler.
func (s *Service) Run() error {
	vErr := s.Configuration.Validate()
	if vErr != nil{
		log.Fatal().Str("err", vErr.DebugReport()).Msg("invalid configuration")
	}

	s.Configuration.Print()

	authConfig, authErr := s.Configuration.LoadAuthConfig()
	if authErr != nil {
		log.Fatal().Str("err", authErr.DebugReport()).Msg("cannot load authx config")
	}

	log.Info().Bool("AllowsAll", authConfig.AllowsAll).Int("permissions", len(authConfig.Permissions)).Msg("Auth config")

	go s.LaunchGRPC(authConfig)
	return s.LaunchHTTP()
}

func (s * Service) LaunchHTTP() error {
	addr := fmt.Sprintf(":%d", s.Configuration.HTTPPort)
	clientAddr := fmt.Sprintf(":%d", s.Configuration.Port)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	mux := runtime.NewServeMux()

	if err := grpc_public_api_go.RegisterOrganizationsHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start organizations handler")
	}
	if err := grpc_public_api_go.RegisterClustersHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start cluster handler")
	}

	log.Info().Str("address", addr).Msg("HTTP Listening")
	return http.ListenAndServe(addr, mux)
}

func (s * Service) LaunchGRPC(authConfig * interceptor.AuthorizationConfig) error {
	clients, cErr := s.GetClients()
	if cErr != nil{
		log.Fatal().Str("err", cErr.DebugReport()).Msg("cannot generate clients")
		return cErr
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
	}

	// Create handlers
	orgManager := organizations.NewManager(clients.orgClient)
	orgHandler := organizations.NewHandler(orgManager)

	clusManager := clusters.NewManager(clients.clusClient, clients.nodeClient)
	clusHandler := clusters.NewHandler(clusManager)

	nodesManager := nodes.NewManager(clients.nodeClient)
	nodesHandler := nodes.NewHandler(nodesManager)

	resManager := resources.NewManager(clients.clusClient, clients.nodeClient)
	resHandler := resources.NewHandler(resManager)

	userManager := users.NewManager(clients.umClient)
	userHandler := users.NewHandler(userManager)

	roleManager := roles.NewManager(clients.umClient)
	roleHandler := roles.NewHandler(roleManager)

	appManager := applications.NewManager(clients.appClient)
	appHandler := applications.NewHandler(appManager)

	grpcServer := grpc.NewServer(interceptor.WithServerAuthxInterceptor(
		interceptor.NewConfig(authConfig, s.Configuration.AuthSecret, s.Configuration.AuthHeader)))
	grpc_public_api_go.RegisterOrganizationsServer(grpcServer, orgHandler)
	grpc_public_api_go.RegisterClustersServer(grpcServer, clusHandler)
	grpc_public_api_go.RegisterNodesServer(grpcServer, nodesHandler)
	grpc_public_api_go.RegisterResourcesServer(grpcServer, resHandler)
	grpc_public_api_go.RegisterUsersServer(grpcServer, userHandler)
	grpc_public_api_go.RegisterRolesServer(grpcServer, roleHandler)
	grpc_public_api_go.RegisterApplicationsServer(grpcServer, appHandler)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}