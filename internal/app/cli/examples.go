package cli

import (
	"github.com/nalej/grpc-application-go"
)

const HeartBeatConfigContent = `heartbeat.monitors:
- type: http
  schedule: '@every 5s'
  urls: ["http://${NALEJ_SERV_WORDPRESS}:80/"]
  check.request:
    method: "GET"
  check.response:
    status: 200
output.elasticsearch:
  hosts: ["${NALEJ_SERV_ELASTIC}:9200"]
`

func (a *Applications) getComplexDescriptor(sType grpc_application_go.StorageType) *grpc_application_go.AddAppDescriptorRequest {

	service1 := &grpc_application_go.Service{
		ServiceId:   "mysql",
		Name:        "simple-mysql",
		Description: "A MySQL instance",
		Type:        grpc_application_go.ServiceType_DOCKER,
		Image:       "mysql:5.6",
		Specs:       &grpc_application_go.DeploySpecs{Replicas: 1},
		Storage:     []*grpc_application_go.Storage{&grpc_application_go.Storage{MountPath: "/tmp"}},
		ExposedPorts: []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "mysqlport", InternalPort: 3306, ExposedPort: 3306,
		}},
		EnvironmentVariables: map[string]string{"MYSQL_ROOT_PASSWORD": "root"},
		Labels:               map[string]string{"app": "simple-mysql", "component": "simple-app"},
	}

	service2 := &grpc_application_go.Service{
		ServiceId:   "wordpress",
		Name:        "simple-wordpress",
		Description: "A Wordpress instance",
		Type:        grpc_application_go.ServiceType_DOCKER,
		Image:       "wordpress:5.0.0",
		Specs:       &grpc_application_go.DeploySpecs{Replicas: 1},
		Storage:     []*grpc_application_go.Storage{&grpc_application_go.Storage{MountPath: "/tmp"}},
		ExposedPorts: []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "wordpressport", InternalPort: 80, ExposedPort: 80,
			Endpoints: []*grpc_application_go.Endpoint{
				&grpc_application_go.Endpoint{
					Type: grpc_application_go.EndpointType_WEB,
					Path: "/",
				},
			},
		}},
		EnvironmentVariables: map[string]string{"WORDPRESS_DB_HOST":"NALEJ_SERV_MYSQL:3306","WORDPRESS_DB_PASSWORD":"root"},
		Labels:               map[string]string{"app": "simple-wordpress", "component": "simple-app"},
	}

	// add additional storage for persistence example
	if sType == grpc_application_go.StorageType_CLUSTER_LOCAL {
		// use persistence storage SQL and wordpress
		service1.Storage = append(service1.Storage, &grpc_application_go.Storage{MountPath: "/var/lib/mysql", Type: sType, Size: int64(1024 * 1024 * 1024)})
		service2.Storage = append(service2.Storage, &grpc_application_go.Storage{MountPath: "/var/www/html", Type: sType, Size: int64(512 * 1024 * 1024)})
	}

	service3 := &grpc_application_go.Service{
		ServiceId:            "kibana",
		Name:                 "kibana",
		Description:          "A Kibana dashboard",
		Type:                 grpc_application_go.ServiceType_DOCKER,
		Image:                "docker.elastic.co/kibana/kibana:6.4.2",
		Specs:                &grpc_application_go.DeploySpecs{Replicas: 1},
		ExposedPorts:         []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "kibanaport", InternalPort: 5601, ExposedPort: 5601,
			Endpoints: []*grpc_application_go.Endpoint{
				&grpc_application_go.Endpoint{
					Type: grpc_application_go.EndpointType_WEB,
					Path: "/",
				},
			},
		}},
		EnvironmentVariables: map[string]string{"ELASTICSEARCH_URL":"http://NALEJ_SERV_ELASTIC:9200"},
		Labels:               map[string]string{"app": "kibana"},
	}

	service4 := &grpc_application_go.Service{
		ServiceId:            "elastic",
		Name:                 "elastic",
		Description:          "Elastic gathering metrics",
		Type:                 grpc_application_go.ServiceType_DOCKER,
		Image:                "docker.elastic.co/elasticsearch/elasticsearch:6.4.2",
		Specs:                &grpc_application_go.DeploySpecs{Replicas: 1},
		Storage:              []*grpc_application_go.Storage{&grpc_application_go.Storage{MountPath: "/usr/share/elasticsearch/data",Type:sType}},
		ExposedPorts:         []*grpc_application_go.Port{&grpc_application_go.Port{
			Name: "elasticport", InternalPort: 9200, ExposedPort: 9200,
		}},
		EnvironmentVariables: map[string]string{
			"cluster.name":"elastic-cluster",
			"bootstrap.memory_lock":"true",
			"ES_JAVA_OPTS":"-Xms512m -Xmx512m",
			"discovery.type":"single-node",
		},
		Labels:               map[string]string{"app": "elastic"},
	}

	service5 := &grpc_application_go.Service{
		ServiceId:            "heartbeat",
		Name:                 "heartbeat",
		Description:          "A tool to gather data from ",
		Type:                 grpc_application_go.ServiceType_DOCKER,
		//Image:                "docker.elastic.co/beats/heartbeat:6.4.2",
		Image:                "nalejops/heartbeat:1.0.0",
		Specs:                &grpc_application_go.DeploySpecs{Replicas: 1},
		/*
		Configs:              []*grpc_application_go.ConfigFile{
			&grpc_application_go.ConfigFile{
				ConfigFileId:         "heartbeat-config",
				Content:              []byte(HeartBeatConfigContent),
				MountPath:            "/conf/heartbeat.yml",
			},
		},
		*/
		Labels:               map[string]string{"app": "heartbeat"},
	}

	secRuleWP := grpc_application_go.SecurityRule{
		Name:            "allow access to wordpress",
		Access:          grpc_application_go.PortAccess_PUBLIC,
		RuleId:          "001",
		SourcePort:      80,
		SourceServiceId: "wordpress",
		AuthServices: []string{"heartbeat"},
	}

	secRuleMysql := grpc_application_go.SecurityRule{
		RuleId:               "002",
		Name:                 "allow access to mysql",
		SourceServiceId:      "mysql",
		SourcePort:           3306,
		Access:               grpc_application_go.PortAccess_APP_SERVICES,
		AuthServices:         []string{"wordpress"},
	}

	secRuleK := grpc_application_go.SecurityRule{
		RuleId: "003",
		Name:                 "allow access to kibana",
		SourceServiceId:      "kibana",
		SourcePort:           5601,
		Access:               grpc_application_go.PortAccess_PUBLIC,
	}

	secRuleElastic := grpc_application_go.SecurityRule{
		RuleId: "004",
		Name:                 "allow access to elastic",
		SourceServiceId:      "elastic",
		SourcePort:           9200,
		Access:               grpc_application_go.PortAccess_APP_SERVICES,
		AuthServices: []string{"kibana", "heartbeat"},
	}

	return &grpc_application_go.AddAppDescriptorRequest{
		Name:        "Sample application with 5 elements",
		Description: "Wordpress with a Kibana monitoring backend",
		Labels:      map[string]string{"app": "simple-app"},
		Rules:       []*grpc_application_go.SecurityRule{&secRuleMysql, &secRuleWP, &secRuleElastic, &secRuleK},
		Services:    []*grpc_application_go.Service{service1, service2, service3, service4, service5},
	}
}
