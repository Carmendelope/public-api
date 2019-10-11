package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-infrastructure-manager-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-monitoring-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-provisioner-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/nalej/public-api/internal/pkg/server/agent"
	"github.com/nalej/public-api/internal/pkg/server/application-network"
	"github.com/nalej/public-api/internal/pkg/server/applications"
	"github.com/nalej/public-api/internal/pkg/server/clusters"
	"github.com/nalej/public-api/internal/pkg/server/devices"
	"github.com/nalej/public-api/internal/pkg/server/ec"
	"github.com/nalej/public-api/internal/pkg/server/inventory"
	"github.com/nalej/public-api/internal/pkg/server/monitoring"
	"github.com/nalej/public-api/internal/pkg/server/nodes"
	"github.com/nalej/public-api/internal/pkg/server/organizations"
	"github.com/nalej/public-api/internal/pkg/server/provisioner"
	"github.com/nalej/public-api/internal/pkg/server/resources"
	"github.com/nalej/public-api/internal/pkg/server/roles"
	"github.com/nalej/public-api/internal/pkg/server/unified-logging"
	"github.com/nalej/public-api/internal/pkg/server/users"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"strings"
)

type Service struct {
	Configuration Config
}

// NewService creates a new system model service.
func NewService(conf Config) *Service {
	return &Service{
		conf,
	}
}

type Clients struct {
	orgClient    grpc_organization_go.OrganizationsClient
	clusClient   grpc_infrastructure_go.ClustersClient
	nodeClient   grpc_infrastructure_go.NodesClient
	infraClient  grpc_infrastructure_manager_go.InfrastructureManagerClient
	umClient     grpc_user_manager_go.UserManagerClient
	appClient    grpc_application_manager_go.ApplicationManagerClient
	deviceClient grpc_device_manager_go.DevicesClient
	ulClient     grpc_unified_logging_go.CoordinatorClient
	mmClient     grpc_monitoring_go.MonitoringManagerClient
	amClient     grpc_monitoring_go.AssetMonitoringClient
	eicClient    grpc_inventory_manager_go.EICClient
	invClient    grpc_inventory_manager_go.InventoryClient
	agentClient  grpc_inventory_manager_go.AgentClient
	appNetClient grpc_application_manager_go.ApplicationNetworkClient
	provisionerClient grpc_provisioner_go.ProvisionClient
}

func (s *Service) GetClients() (*Clients, derrors.Error) {
	smConn, err := grpc.Dial(s.Configuration.SystemModelAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the system model")
	}
	infraConn, err := grpc.Dial(s.Configuration.InfrastructureManagerAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the applications manager")
	}
	umConn, err := grpc.Dial(s.Configuration.UserManagerAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the user manager")
	}
	appConn, err := grpc.Dial(s.Configuration.ApplicationsManagerAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the applications manager")
	}
	devConn, err := grpc.Dial(s.Configuration.DeviceManagerAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the device manager")
	}
	ulConn, err := grpc.Dial(s.Configuration.UnifiedLoggingAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with unified logging coordinator")
	}
	mmConn, err := grpc.Dial(s.Configuration.MonitoringManagerAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with infrastructure monitor coordinator")
	}
	invManagerConn, err := grpc.Dial(s.Configuration.InventoryManagerAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with inventory manager coordinator")
	}
	provConn, err := grpc.Dial(s.Configuration.ProvisionerManagerAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "canot create connection with provisioner manager address")
	}

	oClient := grpc_organization_go.NewOrganizationsClient(smConn)
	cClient := grpc_infrastructure_go.NewClustersClient(smConn)
	nClient := grpc_infrastructure_go.NewNodesClient(smConn)
	infraClient := grpc_infrastructure_manager_go.NewInfrastructureManagerClient(infraConn)
	umClient := grpc_user_manager_go.NewUserManagerClient(umConn)
	appClient := grpc_application_manager_go.NewApplicationManagerClient(appConn)
	deviceClient := grpc_device_manager_go.NewDevicesClient(devConn)
	ulClient := grpc_unified_logging_go.NewCoordinatorClient(ulConn)
	mmClient := grpc_monitoring_go.NewMonitoringManagerClient(mmConn)
	amClient := grpc_monitoring_go.NewAssetMonitoringClient(mmConn)
	eicClient := grpc_inventory_manager_go.NewEICClient(invManagerConn)
	invClient := grpc_inventory_manager_go.NewInventoryClient(invManagerConn)
	agentClient := grpc_inventory_manager_go.NewAgentClient(invManagerConn)
	appNetClient := grpc_application_manager_go.NewApplicationNetworkClient(appConn)
	provClient := grpc_provisioner_go.NewProvisionClient(provConn)

	return &Clients{oClient, cClient, nClient, infraClient, umClient,
		appClient, deviceClient, ulClient, mmClient, amClient,
		eicClient, invClient, agentClient, appNetClient,
		provClient}, nil
}

// Run the service, launch the REST service handler.
func (s *Service) Run() error {
	vErr := s.Configuration.Validate()
	if vErr != nil {
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

// allowCORS allows Cross Origin Resource Sharing from any origin.
// Don't do this without consideration in production systems.
func (s *Service) allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}

func (s *Service) LaunchHTTP() error {
	addr := fmt.Sprintf(":%d", s.Configuration.HTTPPort)
	clientAddr := fmt.Sprintf(":%d", s.Configuration.Port)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	mux := runtime.NewServeMux()

	if err := grpc_public_api_go.RegisterApplicationsHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start applications handler")
	}
	if err := grpc_public_api_go.RegisterClustersHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start cluster handler")
	}
	if err := grpc_public_api_go.RegisterNodesHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start nodes handler")
	}
	if err := grpc_public_api_go.RegisterOrganizationsHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start organizations handler")
	}
	if err := grpc_public_api_go.RegisterResourcesHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start applications handler")
	}
	if err := grpc_public_api_go.RegisterRolesHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start applications handler")
	}
	if err := grpc_public_api_go.RegisterUsersHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start applications handler")
	}
	if err := grpc_public_api_go.RegisterDevicesHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start device handler")
	}
	if err := grpc_public_api_go.RegisterUnifiedLoggingHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start unified logging handler")
	}
	if err := grpc_public_api_go.RegisterEdgeControllersHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start edge controller handler")
	}
	if err := grpc_public_api_go.RegisterInventoryHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start inventory handler")
	}
	if err := grpc_public_api_go.RegisterInventoryMonitoringHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start inventory monitoring handler")
	}
	if err := grpc_public_api_go.RegisterAgentHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start agent handler")
	}
	if err := grpc_public_api_go.RegisterApplicationNetworkHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start application-network handler")
	}
	if err := grpc_public_api_go.RegisterProvisionHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start provision handler")
	}

	server := &http.Server{
		Addr:    addr,
		Handler: s.allowCORS(mux),
	}
	log.Info().Str("address", addr).Msg("HTTP Listening")
	return server.ListenAndServe()
}

func (s *Service) LaunchGRPC(authConfig *interceptor.AuthorizationConfig) error {
	clients, cErr := s.GetClients()
	if cErr != nil {
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

	clusManager := clusters.NewManager(clients.clusClient, clients.nodeClient, clients.infraClient, clients.mmClient)
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

	devManager := devices.NewManager(clients.deviceClient)
	devHandler := devices.NewHandler(devManager)

	ulManager := unified_logging.NewManager(clients.ulClient)
	ulHandler := unified_logging.NewHandler(ulManager)

	ecManager := ec.NewManager(clients.eicClient, clients.agentClient)
	ecHandler := ec.NewHandler(ecManager)

	invManager := inventory.NewManager(clients.invClient, clients.eicClient)
	invHandler := inventory.NewHandler(invManager)

	amManager := monitoring.NewManager(clients.amClient)
	amHandler := monitoring.NewHandler(amManager)

	agentManager := agent.NewManager(clients.agentClient)
	agentHandler := agent.NewHandler(agentManager)

	appNetManager := application_network.NewManager(clients.appNetClient, clients.appClient)
	appNetHandler := application_network.NewHandler(appNetManager)

	provManager := provisioner.NewManager(clients.provisionerClient)
	provHandler := provisioner.NewHandler(provManager)

	grpcServer := grpc.NewServer(interceptor.WithServerAuthxInterceptor(
		interceptor.NewConfig(authConfig, s.Configuration.AuthSecret, s.Configuration.AuthHeader)))
	grpc_public_api_go.RegisterOrganizationsServer(grpcServer, orgHandler)
	grpc_public_api_go.RegisterClustersServer(grpcServer, clusHandler)
	grpc_public_api_go.RegisterNodesServer(grpcServer, nodesHandler)
	grpc_public_api_go.RegisterResourcesServer(grpcServer, resHandler)
	grpc_public_api_go.RegisterUsersServer(grpcServer, userHandler)
	grpc_public_api_go.RegisterRolesServer(grpcServer, roleHandler)
	grpc_public_api_go.RegisterApplicationsServer(grpcServer, appHandler)
	grpc_public_api_go.RegisterDevicesServer(grpcServer, devHandler)
	grpc_public_api_go.RegisterUnifiedLoggingServer(grpcServer, ulHandler)
	grpc_public_api_go.RegisterEdgeControllersServer(grpcServer, ecHandler)
	grpc_public_api_go.RegisterInventoryServer(grpcServer, invHandler)
	grpc_public_api_go.RegisterInventoryMonitoringServer(grpcServer, amHandler)
	grpc_public_api_go.RegisterAgentServer(grpcServer, agentHandler)
	grpc_public_api_go.RegisterApplicationNetworkServer(grpcServer, appNetHandler)
	grpc_public_api_go.RegisterProvisionServer(grpcServer, provHandler)

	if s.Configuration.Debug {
		log.Info().Msg("Enabling gRPC server reflection")
		// Register reflection service on gRPC server.
		reflection.Register(grpcServer)
	}
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}
