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

// This file contains the variables that are used through the commands.

package commands

var loginPort int
var email string
var password string
var publicKeyPath string
var title string
var phone string
var location string
var lastName string
var updateName bool
var updateLocation bool
var updatePhone bool
var updateTitle bool
var updateLastName bool

var newPassword string

var name string
var roleName string
var roleID string

var organizationID string
var clusterID string
var username string
var privateKeyPath string
var kubeConfigPath string
var hostname string
var nodes []string
var targetPlatform string
var useStaticIPAddresses bool
var ipAddressIngress string

var descriptorID string
var descriptorPath string
var params string
var connections string

var instanceID string
var sgInstanceID string
var sgID string
var serviceID string
var serviceInstanceID string

var internal bool

var exampleName string

var storageType string

var enabled bool
var disabled bool
var enabledDefaultConnectivity bool
var disabledDefaultConnectivity bool

var deviceGroupID string
var deviceID string
var rawLabels string
var nodeID string

var message string
var from string
var to string
var redirectLog bool
var desc bool
var follow bool
var nFirst bool
var metadata bool

var rangeMinutes int32
var clusterStatFields string

var watch bool

// cluster update
var millicoresConversionFactor float64

var outputPath string
var edgeControllerID string
var assetID string
var activate bool
var geolocation string
var assetDeviceId string

var force bool

var agentTypeRaw string
var sudoer bool

var sourceInstanceID string
var targetInstanceID string
var inbound string
var outbound string

var requestId string

var provisionClusterName string
var provisionAzureCredentialsPath string
var provisionAzureDnsZoneName string
var provisionAzureResourceGroup string
var provisionClusterType string
var provisionIsProductionCluster bool
var provisionKubernetesVersion string
var provisionNodeType string
var provisionNumNodes int
var provisionTargetPlatform string
var provisionZone string
var provisionKubeConfigOutputPath string
