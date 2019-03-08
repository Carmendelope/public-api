package cli

import (
	"github.com/nalej/grpc-application-go"
)

const HeartBeatConfigContent = `heartbeat.monitors:
- type: http
  schedule: '@every 5s'
  urls: ["http://${NALEJ_SERV_SIMPLE-WORDPRESS}:80/"]
  check.request:
    method: "GET"
  check.response:
    status: 200
output.elasticsearch:
  hosts: ["${NALEJ_SERV_ELASTIC}:9200"]
`

func (a *Applications) getBasicDescriptor(sType grpc_application_go.StorageType) *grpc_application_go.AddAppDescriptorRequest {

	service1 := &grpc_application_go.Service{
		Name:        "simple-mysql",
		Type:        grpc_application_go.ServiceType_DOCKER,
		Image:       "mysql:5.6",
		Specs:       &grpc_application_go.DeploySpecs{Replicas: 1},
		Storage:     []*grpc_application_go.Storage{&grpc_application_go.Storage{MountPath: "/tmp", Type: grpc_application_go.StorageType_EPHEMERAL, Size: int64(100 * 1024 * 1024)}},
		ExposedPorts: []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "mysqlport", InternalPort: 3306, ExposedPort: 3306,
		}},
		EnvironmentVariables: map[string]string{"MYSQL_ROOT_PASSWORD": "root"},
		Labels:               map[string]string{"app": "simple-mysql", "component": "simple-app"},
	}

	service2 := &grpc_application_go.Service{
		Name:        "simple-wordpress",
		Type:        grpc_application_go.ServiceType_DOCKER,
		Image:       "wordpress:5.0.0",
		Specs:       &grpc_application_go.DeploySpecs{Replicas: 1},
		DeployAfter: []string{"simple-mysql"},
		Storage:     []*grpc_application_go.Storage{&grpc_application_go.Storage{MountPath: "/tmp", Type: grpc_application_go.StorageType_EPHEMERAL, Size: int64(100 * 1024 * 1024)}},
		ExposedPorts: []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "wordpressport", InternalPort: 80, ExposedPort: 80,
			Endpoints: []*grpc_application_go.Endpoint{
				&grpc_application_go.Endpoint{
					Type: grpc_application_go.EndpointType_WEB,
					Path: "/",
				},
			},
		}},
		EnvironmentVariables: map[string]string{"WORDPRESS_DB_HOST": "NALEJ_SERV_SIMPLE-MYSQL:3306", "WORDPRESS_DB_PASSWORD": "root"},
		Labels:               map[string]string{"app": "simple-wordpress", "component": "simple-app"},
	}

	group1 := &grpc_application_go.ServiceGroup{
		Name: "application",
		Services: []*grpc_application_go.Service{service1, service2},
		Specs: &grpc_application_go.ServiceGroupDeploymentSpecs{NumReplicas:1,MultiClusterReplica:false},
	}

	// add additional storage for persistence example
	if sType == grpc_application_go.StorageType_CLUSTER_LOCAL {
		// use persistence storage SQL and wordpress
		service1.Storage = append(service1.Storage, &grpc_application_go.Storage{MountPath: "/var/lib/mysql", Type: sType, Size: int64(1024 * 1024 * 1024)})
		service2.Storage = append(service2.Storage, &grpc_application_go.Storage{MountPath: "/var/www/html", Type: sType, Size: int64(512 * 1024 * 1024)})
	}
	secRule := grpc_application_go.SecurityRule{
		Name:            "allow access to wordpress",
		Access:          grpc_application_go.PortAccess_PUBLIC,
		RuleId:          "001",
		TargetPort:      80,
		TargetServiceName: "simple-wordpress",
		TargetServiceGroupName: "application",
	}

	return &grpc_application_go.AddAppDescriptorRequest{
		Name:        "Sample application",
		Labels:      map[string]string{"app": "simple-app"},
		Rules:       []*grpc_application_go.SecurityRule{&secRule},
		Groups:      []*grpc_application_go.ServiceGroup{group1},
	}
}

func (a *Applications) getComplexDescriptor(sType grpc_application_go.StorageType) *grpc_application_go.AddAppDescriptorRequest {

	service1 := &grpc_application_go.Service{
		ServiceGroupId: "application",
		Name:        "simple-mysql",
		Type:        grpc_application_go.ServiceType_DOCKER,
		Image:       "mysql:5.6",
		Specs:       &grpc_application_go.DeploySpecs{Replicas: 1},
		Storage:     []*grpc_application_go.Storage{&grpc_application_go.Storage{MountPath: "/tmp", Type: grpc_application_go.StorageType_EPHEMERAL, Size: int64(100 * 1024 * 1024)}},
		ExposedPorts: []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "mysqlport", InternalPort: 3306, ExposedPort: 3306,
		}},
		EnvironmentVariables: map[string]string{"MYSQL_ROOT_PASSWORD": "root"},
		Labels:               map[string]string{"app": "simple-mysql", "component": "simple-app"},
	}

	service2 := &grpc_application_go.Service{
		ServiceGroupId: "application",
		Name:        "simple-wordpress",
		Type:        grpc_application_go.ServiceType_DOCKER,
		Image:       "wordpress:5.0.0",
		Specs:       &grpc_application_go.DeploySpecs{Replicas: 1},
		Storage:     []*grpc_application_go.Storage{&grpc_application_go.Storage{MountPath: "/tmp", Type: grpc_application_go.StorageType_EPHEMERAL, Size: int64(100 * 1024 * 1024)}},
		ExposedPorts: []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "wordpressport", InternalPort: 80, ExposedPort: 80,
			Endpoints: []*grpc_application_go.Endpoint{
				&grpc_application_go.Endpoint{
					Type: grpc_application_go.EndpointType_WEB,
					Path: "/",
				},
			},
		}},
		EnvironmentVariables: map[string]string{"WORDPRESS_DB_HOST": "NALEJ_SERV_SIMPLE-MYSQL:3306", "WORDPRESS_DB_PASSWORD": "root"},
		Labels:               map[string]string{"app": "simple-wordpress", "component": "simple-app"},
	}

	// add additional storage for persistence example
	if sType == grpc_application_go.StorageType_CLUSTER_LOCAL {
		// use persistence storage SQL and wordpress
		service1.Storage = append(service1.Storage, &grpc_application_go.Storage{MountPath: "/var/lib/mysql", Type: sType, Size: int64(1024 * 1024 * 1024)})
		service2.Storage = append(service2.Storage, &grpc_application_go.Storage{MountPath: "/var/www/html", Type: sType, Size: int64(512 * 1024 * 1024)})
	}

	service3 := &grpc_application_go.Service{
		ServiceGroupId: "application",
		ServiceId:   "kibana",
		Name:        "kibana",
		Type:        grpc_application_go.ServiceType_DOCKER,
		Image:       "docker.elastic.co/kibana/kibana:6.4.2",
		Specs:       &grpc_application_go.DeploySpecs{Replicas: 1},
		ExposedPorts: []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "kibanaport", InternalPort: 5601, ExposedPort: 5601,
			Endpoints: []*grpc_application_go.Endpoint{
				&grpc_application_go.Endpoint{
					Type: grpc_application_go.EndpointType_WEB,
					Path: "/",
				},
			},
		}},
		EnvironmentVariables: map[string]string{"ELASTICSEARCH_URL": "http://NALEJ_SERV_ELASTIC:9200"},
		Labels:               map[string]string{"app": "kibana"},
	}

	service4 := &grpc_application_go.Service{
		ServiceGroupId: "application",
		ServiceId:   "elastic",
		Name:        "elastic",
		Type:        grpc_application_go.ServiceType_DOCKER,
		Image:       "docker.elastic.co/elasticsearch/elasticsearch:6.4.2",
		Specs:       &grpc_application_go.DeploySpecs{Replicas: 1},
		Storage:     []*grpc_application_go.Storage{&grpc_application_go.Storage{MountPath: "/usr/share/elasticsearch/data", Type: sType}},
		ExposedPorts: []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "elasticport", InternalPort: 9200, ExposedPort: 9200,
		}},
		EnvironmentVariables: map[string]string{
			"cluster.name":          "elastic-cluster",
			"bootstrap.memory_lock": "true",
			"ES_JAVA_OPTS":          "-Xms512m -Xmx512m",
			"discovery.type":        "single-node",
		},
		Labels: map[string]string{"app": "elastic"},
	}

	service5 := &grpc_application_go.Service{
		ServiceGroupId: "application",
		ServiceId:   "heartbeat",
		Name:        "heartbeat",
		Type:        grpc_application_go.ServiceType_DOCKER,
		Image:       "docker.elastic.co/beats/heartbeat:6.4.2",
		Specs: &grpc_application_go.DeploySpecs{Replicas: 1},
			Configs:              []*grpc_application_go.ConfigFile{
				&grpc_application_go.ConfigFile{
					ConfigFileId:         "heartbeat-config",
					Content:              []byte(HeartBeatConfigContent),
					MountPath:            "/conf/heartbeat.yml",
				},
			},
		Labels: map[string]string{"app": "heartbeat"},
		RunArguments: []string{"--path.config=/conf/"},
	}

	secRuleWP := grpc_application_go.SecurityRule{
		Name:            "allow access to wordpress",
		Access:          grpc_application_go.PortAccess_PUBLIC,
		RuleId:          "001",
		TargetPort:      80,
		TargetServiceName: "simple-wordpress",
		TargetServiceGroupName: "application",
	}

	secRuleWP2 := grpc_application_go.SecurityRule{
		Name:            "allow access to wordpress to heartbeat",
		Access:          grpc_application_go.PortAccess_APP_SERVICES,
		RuleId:          "001",
		TargetPort:      80,
		TargetServiceName: "simple-wordpress",
		TargetServiceGroupName: "application",
		AuthServiceGroupName: "application",
		AuthServices:    []string{"heartbeat"},
	}

	secRuleMysql := grpc_application_go.SecurityRule{
		RuleId:          "002",
		Name:            "allow access to mysql",
		TargetPort:      3306,
		TargetServiceName: "simple-mysql",
		TargetServiceGroupName: "application",
		Access:          grpc_application_go.PortAccess_APP_SERVICES,
		AuthServiceGroupName: "application",
		AuthServices:    []string{"simple-wordpress"},
	}
	secRuleElastic := grpc_application_go.SecurityRule{
		RuleId:          "004",
		Name:            "allow access to elastic",
		TargetPort:      9200,
		TargetServiceName: "elastic",
		TargetServiceGroupName: "application",
		Access:          grpc_application_go.PortAccess_APP_SERVICES,
		AuthServiceGroupName: "application",
		AuthServices:    []string{"kibana", "heartbeat"},
	}

	secRuleK := grpc_application_go.SecurityRule{
		RuleId:          "003",
		Name:            "allow access to kibana",
		TargetPort:      5601,
		TargetServiceName: "kibana",
		TargetServiceGroupName: "application",
		Access:          grpc_application_go.PortAccess_PUBLIC,
	}



	group1 := &grpc_application_go.ServiceGroup{
		Name: "application",
		Services: []*grpc_application_go.Service{
			service1, service2, service3, service4, service5,
		},
		Specs: &grpc_application_go.ServiceGroupDeploymentSpecs{NumReplicas:1,MultiClusterReplica:false},
	}

	return &grpc_application_go.AddAppDescriptorRequest{
		Name:        "Sample application with 5 elements",
		Labels:      map[string]string{"app": "simple-app"},
		Rules:       []*grpc_application_go.SecurityRule{&secRuleWP, &secRuleK, &secRuleMysql, &secRuleElastic, &secRuleWP2},
		Groups:      []*grpc_application_go.ServiceGroup{group1},
	}
}

func (a *Applications) getMultiReplicaDescriptor(sType grpc_application_go.StorageType) *grpc_application_go.AddAppDescriptorRequest {

	service1 := &grpc_application_go.Service{
		Name:        "simple-mysql",
		Type:        grpc_application_go.ServiceType_DOCKER,
		Image:       "mysql:5.6",
		Specs:       &grpc_application_go.DeploySpecs{Replicas: 1},
		Storage:     []*grpc_application_go.Storage{{MountPath: "/tmp", Type: grpc_application_go.StorageType_EPHEMERAL, Size: int64(100 * 1024 * 1024)}},
		ExposedPorts: []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "mysqlport", InternalPort: 3306, ExposedPort: 3306,
		}},
		EnvironmentVariables: map[string]string{"MYSQL_ROOT_PASSWORD": "root"},
		Labels:               map[string]string{"app": "simple-mysql", "component": "simple-app"},
	}


	group1 := &grpc_application_go.ServiceGroup{
		Name: "database",
		Services: []*grpc_application_go.Service{service1},
		Specs: &grpc_application_go.ServiceGroupDeploymentSpecs{NumReplicas:1,MultiClusterReplica:false},
	}

	service2 := &grpc_application_go.Service{
		Name:        "simple-wordpress",
		Type:        grpc_application_go.ServiceType_DOCKER,
		Image:       "wordpress:5.0.0",
		Specs:       &grpc_application_go.DeploySpecs{Replicas: 1},
		DeployAfter: []string{"simple-mysql"},
		Storage:     []*grpc_application_go.Storage{&grpc_application_go.Storage{MountPath: "/tmp", Type: grpc_application_go.StorageType_EPHEMERAL, Size: int64(100 * 1024 * 1024)}},
		ExposedPorts: []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "wordpressport", InternalPort: 80, ExposedPort: 80,
			Endpoints: []*grpc_application_go.Endpoint{
				&grpc_application_go.Endpoint{
					Type: grpc_application_go.EndpointType_WEB,
					Path: "/",
				},
			},
		}},
		EnvironmentVariables: map[string]string{"WORDPRESS_DB_HOST": "NALEJ_SERV_SIMPLE-MYSQL:3306", "WORDPRESS_DB_PASSWORD": "root"},
		Labels:               map[string]string{"app": "simple-wordpress", "component": "simple-app"},
	}

	group2 := &grpc_application_go.ServiceGroup{
		Name: "front",
		Services: []*grpc_application_go.Service{service2},
		Specs: &grpc_application_go.ServiceGroupDeploymentSpecs{NumReplicas:0,MultiClusterReplica:true},
	}

	// add additional storage for persistence example
	if sType == grpc_application_go.StorageType_CLUSTER_LOCAL {
		// use persistence storage SQL and wordpress
		service1.Storage = append(service1.Storage, &grpc_application_go.Storage{MountPath: "/var/lib/mysql", Type: sType, Size: int64(1024 * 1024 * 1024)})
		service2.Storage = append(service2.Storage, &grpc_application_go.Storage{MountPath: "/var/www/html", Type: sType, Size: int64(512 * 1024 * 1024)})
	}
	secRule := grpc_application_go.SecurityRule{
		Name:            "allow access to wordpress",
		Access:          grpc_application_go.PortAccess_PUBLIC,
		RuleId:          "001",
		TargetPort:      80,
		TargetServiceName: "simple-wordpress",
		TargetServiceGroupName: "front",
	}

	return &grpc_application_go.AddAppDescriptorRequest{
		Name:        "Multireplica Sample application",
		Labels:      map[string]string{"app": "simple-app"},
		Rules:       []*grpc_application_go.SecurityRule{&secRule},
		Groups:      []*grpc_application_go.ServiceGroup{group1, group2},
	}
}