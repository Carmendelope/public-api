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
 *
 */

package entities

// Definition of the Nalej application descriptors JSON schema

const APP_DESC_SCHEMA = `
{
  "definitions": {
    "labels": {
      "$id": "#/definitions/labels",
      "type": "object",
      "title": "The Labels Schema",
      "additionalProperties": {
        "type": "string",
        "minItems": 1,
        "minLength": 1,
        "maxLength": 63
      }
    },
    "host_port": {
      "$id": "#/definitions/host_port",
      "type": "integer",
      "title": "Host port",
      "minimum": 1,
      "maximum": 65535
    },
    "security_rule": {
      "$id": "#/definitions/security_rule",
      "type": "object",
      "title": "Security connectivity rules",
      "required": [
        "name",
        "target_service_group_name",
        "target_service_name",
        "target_port",
        "access"
      ],
      "properties": {
        "name": {
          "title": "Rule name",
          "type": "string",
          "minLength": 1,
          "maxLength": 63
        },
        "target_service_group_name": {
          "title": "Name of the target service group",
          "type": "string",
          "minLength": 1,
          "maxLength": 63
        },
        "target_service_name": {
          "title": "Name of the target service contained by the service group",
          "type": "string",
          "minLength": 1,
          "maxLength": 63
        },
        "target_port": {
          "title": "Access port",
          "$ref": "#/definitions/host_port"
        },
        "access": {
          "title": "Port this rule refers to",
          "type": "integer",
          "$comment": "ALL_APP_SERVICES,APP_SERVICES,PUBLIC,DEVICE_GROUP",
          "enum": [0, 1, 2, 3]
        },
        "auth_service_group_name": {
          "title": "Name of the group with permission granted to access the target_service_name",
          "type": "string",
          "minLength": 1,
          "maxLength": 63
        },
        "auth_services": {
          "type": "array",
          "title": "List of services authenticated to access",
          "minLength": 1,
          "items": {
            "type": "string",
            "minLength": 1,
            "maxLength": 63
          }
        },
        "device_group_names": {
          "type": "array",
          "minLength": 1,
          "title": "List of device group names with access granted",
          "items": {
            "type": "string",
            "minLength": 1,
            "maxLength": 63
          }
        }
      }
    },
    "service_group_deployment_specs": {
      "$id": "#/definitions/service_group_deployment_specs",
      "title": "Definition of deployment specs for a service group",
      "properties": {
        "multi_cluster_replica": {
          "title": "Set the multiple cluster replication policy",
          "type": "boolean"
        },
        "replicas": {
          "title": "Number of replicas for the service group",
          "type": "integer",
          "minimum": 0
        },
        "deployment_selectors": {
          "title": "Set of labels to be matched by target application clusters",
          "$ref": "#/definitions/labels"
        }
      },
      "oneOf": [
        {
          "required": ["multi_cluster_replica"],
          "not": {"required": ["replicas"]}
        },
        {
          "required": ["replicas"],
          "not": {"required": ["multi_cluster_replica"]}
        },
        {
          "oneOf": [
            {},
            {
              "required": ["multi_cluster_replica"]
            },
            {
              "required": ["replicas"]
            }
          ]
        }
      ]
    },
    "image_credentials": {
      "$id": "#/definitions/image_credentials",
      "title": "Credentials for an image",
      "properties": {
        "username": {
          "type": "string"
        },
        "password": {
          "type": "string"
        },
        "email": {
          "type": "string",
          "pattern": "^([a-zA-Z0-9_\\-\\.]+)@([a-zA-Z0-9_\\-\\.]+)\\.([a-zA-Z]{2,5})$"
        },
        "docker_repository":{
          "type": "string"
        }
      }
    },
    "deploy_specs": {
      "$id": "#/definitions/deploy_specs",
      "title": "Deployment specifications for a service",
      "properties": {
        "cpu": {
          "type": "number",
          "title": "Ratio of reserved cpu",
          "minimum": 0.1
        },
        "memory": {
          "type": "integer",
          "title": "Amount of memory required",
          "minimum": 16
        },
        "replicas": {
          "type": "integer",
          "title": "Number of replicas of this service",
          "minimum": 0,
          "maximum": 255
        }
      }
    },
    "storage": {
      "$id": "#/definitions/storage",
      "type": "object",
      "title": "Storage service definition",
      "required": [
        "size",
        "mount_path"
      ],
      "properties": {
        "size": {
          "title": "Size of the storage volume",
          "type": "integer",
          "minimum": 100
        },
        "mount_path": {
          "title": "Path to mount the volume in the service instance",
          "type": "string"
        }
      }
    },
    "port": {
      "$id": "#/definitions/port",
      "type": "object",
      "title": "Definition of an exposed port",
      "required": [
        "name",
        "internal_port",
        "exposed_port"
      ],
      "properties": {
        "name": {
          "type": "string",
          "title": "Name of the port",
          "minLength": 2,
          "maxLength": 63
        },
        "internal_port": {
          "$ref": "#/definitions/host_port",
          "title": "Internal image port"
        },
        "exposed_port": {
          "$ref": "#/definitions/host_port",
          "title": "Exposed image port"
        },
        "endpoints": {
          "type": "array",
          "title": "List of endpoints for the service",
          "minLength": 1,
          "items": {
            "$ref": "#/definitions/endpoint"
          }
        },
        "environment_variables": {
          "$ref": "#/definitions/labels",
          "title": "Map of environment variables for the service"
        },
        "configs": {
          "$ref": "#/definitions/labels",
          "title": "Map of configuration options for the application"
        },
        "labels": {
          "$ref": "#/definitions/labels",
          "title": "Labels for this service"
        },
        "deploy_after": {
          "type": "array",
          "title": "Name of services that have to be deployed before this",
          "minLength": 1,
          "items": {
            "type": "string"
          }
        },
        "run_arguments": {
          "type": "array",
          "title": "List of running arguments for the service",
          "minLength": 1,
          "items": {
            "type": "string"
          }
        }
      }
    },
    "endpoint": {
      "$id": "#/definitions/endpoint",
      "type": "object",
      "title": "Endpoint definition",
      "required": [
        "path",
        "type"
      ],
      "properties": {
        "path": {
          "type": "string",
          "minLength": 1
        },
        "type": {
          "type": "integer",
          "$comment": "IS_ALIVE=0; REST=1; WEB=2; PROMETHEUS=3; INGESTION=4;",
          "enum": [0,1,2,3,4]
        }
      }
    },
    "service": {
      "$id": "#/definitions/service",
      "title": "Definition of service",
      "required": [
        "name",
        "image"
      ],
      "properties": {
        "name": {
          "type": "string",
          "title": "Name of a service"
        },
        "image": {
          "type": "string",
          "title": "Name of the image to download"
        },
        "image_credentials": {
          "title": "Definition of credentials to download the image",
          "$ref": "#/definitions/image_credentials"
        },
        "specs": {
          "title": "Service deployment specs",
          "$ref": "#/definitions/deploy_specs"
        },
        "storage": {
          "type": "array",
          "title": "Storage definition for this service",
          "minLength": 1,
          "items": {
            "$ref": "#/definitions/storage"
          }
        },
        "exposed_ports": {
          "type": "array",
          "title": "List of exposed ports",
          "minLength": 1,
          "items": {
            "$ref": "#/definitions/port"
          }
        }
      }
    },
    "service_group": {
      "$id": "#/definitions/service_group",
      "type": "object",
      "title": "Group of services to be allocated together",
      "required": [
        "name",
        "services"
      ],
      "properties": {
        "name": {
          "type": "string",
          "title": "Name of the service group",
          "minLength": 4,
          "maxLength": 63
        },
        "specs": {
          "title": "Deployment specifications for this service group",
          "$ref": "#/definitions/service_group_deployment_specs"
        },
        "services": {
          "title": "Array of defined services",
          "type": "array",
          "minLength": 1,
          "items": {
            "$ref": "#/definitions/service"
          }
        }
      }
    }
  },
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "http://nalej.com/app_descriptor.json",
  "type": "object",
  "title": "Nalej application descriptor",
  "required": [
    "name",
    "groups"
  ],
  "properties": {
    "name": {
      "$id": "#/properties/name",
      "type": "string",
      "minLength": 4,
      "maxLength": 63,
      "title": "Name of the application descriptor",
      "pattern": "^(.*)$"
    },
    "labels": {
      "$id": "#/properties/labels",
      "title": "Labels for this app",
      "$ref": "#/definitions/labels"
    },
    "rules": {
      "$id": "#/properties/rules",
      "title": "Connectivity rules",
      "type": "array",
      "items": {
        "$ref": "#/definitions/security_rule",
        "minLength": 1
      }
    },
    "groups": {
      "$id": "#/properties/groups",
      "type": "array",
      "minLength": 1,
      "items": {
        "$ref": "#/definitions/service_group"
      }
    }
  }
}
`
