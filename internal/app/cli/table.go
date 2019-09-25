package cli

import (
	"fmt"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-application-network-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-infrastructure-manager-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-monitoring-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/rs/zerolog/log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const MinWidth = 5
const TabWidth = 2
const Padding = 3

const AppInstanceHeader = ""

var Headers = map[string]string{}

type ResultTable struct {
	data [][]string
}

func AsTable(result interface{}, labelLength int) *ResultTable {
	log.Debug().Int("labelLength", labelLength).Msg("Label length")
	switch result.(type) {
	case *grpc_public_api_go.OrganizationInfo:
		return FromOrganizationInfo(result.(*grpc_public_api_go.OrganizationInfo))
	case *grpc_public_api_go.User:
		return FromUser(result.(*grpc_public_api_go.User))
	case *grpc_user_manager_go.User:
		return FromUserManagerUser(result.(*grpc_user_manager_go.User))
	case *grpc_public_api_go.UserList:
		return FromUserList(result.(*grpc_public_api_go.UserList))
	case *grpc_public_api_go.Cluster:
		return FromCluster(result.(*grpc_public_api_go.Cluster), labelLength)
	case *grpc_monitoring_go.ClusterSummary:
		return FromClusterSummary(result.(*grpc_monitoring_go.ClusterSummary))
	case *grpc_public_api_go.ClusterList:
		return FromClusterList(result.(*grpc_public_api_go.ClusterList), labelLength)
	case *grpc_infrastructure_manager_go.InstallResponse:
		return FromInstallResponse(result.(*grpc_infrastructure_manager_go.InstallResponse))
	case *grpc_public_api_go.AppInstanceList:
		return FromAppInstanceList(result.(*grpc_public_api_go.AppInstanceList), labelLength)
	case *grpc_public_api_go.AppInstance:
		return FromAppInstance(result.(*grpc_public_api_go.AppInstance), labelLength)
	case *grpc_application_go.InstanceParameterList:
		return FromInstanceParameterList(result.(*grpc_application_go.InstanceParameterList))
	case *grpc_application_manager_go.DeploymentResponse:
		return FromDeploymentResponse(result.(*grpc_application_manager_go.DeploymentResponse))
	case *grpc_application_go.AppDescriptorList:
		return FromAppDescriptorList(result.(*grpc_application_go.AppDescriptorList), labelLength)
	case *grpc_application_go.AppDescriptor:
		return FromAppDescriptor(result.(*grpc_application_go.AppDescriptor), labelLength)
	case *grpc_public_api_go.AppParameterList:
		return FromAppParameterList(result.(*grpc_public_api_go.AppParameterList))
	case *grpc_device_manager_go.DeviceGroup:
		return FromDeviceGroup(result.(*grpc_device_manager_go.DeviceGroup))
	case *grpc_device_manager_go.DeviceGroupList:
		return FromDeviceGroupList(result.(*grpc_device_manager_go.DeviceGroupList))
	case *grpc_public_api_go.Device:
		return FromDevice(result.(*grpc_public_api_go.Device), labelLength)
	case *grpc_public_api_go.DeviceList:
		return FromDeviceList(result.(*grpc_public_api_go.DeviceList), labelLength)
	case *grpc_unified_logging_go.LogResponse:
		return FromLogResponse(result.(*grpc_unified_logging_go.LogResponse))
	case *grpc_public_api_go.Node:
		return FromNode(result.(*grpc_public_api_go.Node), labelLength)
	case *grpc_public_api_go.NodeList:
		return FromNodeList(result.(*grpc_public_api_go.NodeList), labelLength)
	case *grpc_public_api_go.Role:
		return FromRole(result.(*grpc_public_api_go.Role))
	case *grpc_public_api_go.RoleList:
		return FromRoleList(result.(*grpc_public_api_go.RoleList))
	case *grpc_inventory_manager_go.EICJoinToken:
		return FromEICJoinToken(result.(*grpc_inventory_manager_go.EICJoinToken))
	case *grpc_public_api_go.InventoryList:
		return FromInventoryList(result.(*grpc_public_api_go.InventoryList), labelLength)
	case *grpc_inventory_manager_go.AgentJoinToken:
		return FromAgentJoinToken(result.(*grpc_inventory_manager_go.AgentJoinToken))
	case *grpc_public_api_go.EdgeControllerExtendedInfo:
		return FromEdgeControllerExtendedInfo(result.(*grpc_public_api_go.EdgeControllerExtendedInfo), labelLength)
	case *grpc_public_api_go.Asset:
		return FromAsset(result.(*grpc_public_api_go.Asset))
	case *grpc_public_api_go.AgentOpResponse:
		return FromAgentOpResponse(result.(*grpc_public_api_go.AgentOpResponse))
	case *grpc_public_api_go.ECOpResponse:
		return FromECOpResponse(result.(*grpc_public_api_go.ECOpResponse))
	case *grpc_common_go.Success:
		return FromSuccess(result.(*grpc_common_go.Success))
	case *grpc_inventory_go.Asset:
		return FromIAsset(result.(*grpc_inventory_go.Asset), labelLength)
	case *grpc_inventory_go.EdgeController:
		return FromIEdgeController(result.(*grpc_inventory_go.EdgeController), labelLength)
	case *grpc_inventory_manager_go.InventorySummary:
		return FromInventorySummary(result.(*grpc_inventory_manager_go.InventorySummary))
	case *grpc_monitoring_go.QueryMetricsResult:
		return FromQueryMetricsResult(result.(*grpc_monitoring_go.QueryMetricsResult))
	case *grpc_monitoring_go.MetricsList:
		return FromMetricsList(result.(*grpc_monitoring_go.MetricsList))
	case *grpc_application_manager_go.AvailableInstanceInboundList:
		return FromAvailableInboundList(result.(*grpc_application_manager_go.AvailableInstanceInboundList))
	case *grpc_application_manager_go.AvailableInstanceOutboundList:
		return FromAvailableOutboundList(result.(*grpc_application_manager_go.AvailableInstanceOutboundList))
	case *grpc_application_network_go.ConnectionInstanceList:
		return FromConnectionInstanceListResult(result.(*grpc_application_network_go.ConnectionInstanceList))
	case *grpc_public_api_go.OpResponse:
		return FromOpResponse(result.(*grpc_public_api_go.OpResponse))
	default:
		log.Fatal().Str("type", fmt.Sprintf("%T", result)).Msg("unsupported")
	}
	return nil
}

func (t *ResultTable) Print() {
	w := tabwriter.NewWriter(os.Stdout, MinWidth, TabWidth, Padding, ' ', 0)
	for _, d := range t.data {
		toPrint := strings.Join(d, "\t")
		fmt.Fprintln(w, toPrint)
	}
	w.Flush()
}

func PrintFromValues(header []string, values [][]string) {
	w := tabwriter.NewWriter(os.Stdout, MinWidth, TabWidth, Padding, ' ', 0)
	fmt.Fprintln(w, strings.Join(header, "\t"))
	for _, d := range values {
		toPrint := strings.Join(d, "\t")
		fmt.Fprintln(w, toPrint)
	}
	w.Flush()
}

func TransformLabels(labels map[string]string, labelLength int) string {
	r := make([]string, 0)

	sortedKeys := GetSortedKeys(labels)
	for _, k := range sortedKeys {
		label := fmt.Sprintf("%s:%s", k, labels[k])
		r = append(r, label)
	}
	labelString := strings.Join(r, ",")
	truncatedR := TruncateString(labelString, labelLength)
	return truncatedR
}

func GetSortedKeys(labels map[string]string) []string {
	sortedKeys := make([]string, len(labels))
	i := 0
	for k := range labels {
		sortedKeys[i] = k
		i++
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

func TruncateString(text string, length int) string {
	log.Debug().Int("length", length).Str("text", text).Msg("truncate")
	if length <= 0 {
		return text
	}
	truncatedString := text
	if len(text) > length {
		if length > 3 {
			length -= 3
		}
		truncatedString = text[0:length] + "..."
	}
	return truncatedString
}

func FromOrganizationInfo(info *grpc_public_api_go.OrganizationInfo) *ResultTable {
	result := make([][]string, 0)
	result = append(result, []string{"ID", "NAME"})
	result = append(result, []string{info.OrganizationId, info.Name})
	return &ResultTable{result}
}

// ----
// Users
// ----

func FromUserManagerUser(user *grpc_user_manager_go.User) *ResultTable {
	result := make([][]string, 0)
	result = append(result, []string{"NAME", "ROLE", "EMAIL"})
	result = append(result, []string{user.Name, user.RoleName, user.Email})
	return &ResultTable{result}
}

func FromUser(user *grpc_public_api_go.User) *ResultTable {
	result := make([][]string, 0)
	result = append(result, []string{"NAME", "ROLE", "EMAIL"})
	result = append(result, []string{user.Name, user.RoleName, user.Email})
	return &ResultTable{result}
}

func FromUserList(user *grpc_public_api_go.UserList) *ResultTable {
	result := make([][]string, 0)
	result = append(result, []string{"NAME", "ROLE", "EMAIL"})
	for _, u := range user.Users {
		result = append(result, []string{u.Name, u.RoleName, u.Email})
	}
	return &ResultTable{result}
}

// ----
// Clusters
// ----

func FromCluster(result *grpc_public_api_go.Cluster, labelLength int) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"NAME", "ID", "NODES", "LABELS", "STATUS"})
	r = append(r, []string{result.Name, result.ClusterId, fmt.Sprintf("%d", result.TotalNodes), TransformLabels(result.Labels, labelLength), result.StatusName})
	return &ResultTable{r}
}

func FromClusterList(result *grpc_public_api_go.ClusterList, labelLength int) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"NAME", "ID", "NODES", "LABELS", "STATUS"})
	for _, c := range result.Clusters {
		r = append(r, []string{c.Name, c.ClusterId, fmt.Sprintf("%d", c.TotalNodes), TransformLabels(c.Labels, labelLength), c.StatusName})
	}
	return &ResultTable{r}
}

func FromClusterSummary(result *grpc_monitoring_go.ClusterSummary) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"CPU", "MEM", "STORAGE"})

	cpuPercentage := int((float64(result.CpuMillicores.Available) / float64(result.CpuMillicores.Total)) * 100)
	cpu := fmt.Sprintf("%d/%d (%d%%)", result.CpuMillicores.Available, result.CpuMillicores.Total, cpuPercentage)
	memPercentage := int((float64(result.MemoryBytes.Available) / float64(result.MemoryBytes.Total)) * 100)
	mem := fmt.Sprintf("%d/%d (%d%%)", result.MemoryBytes.Available, result.MemoryBytes.Total, memPercentage)
	storagePercentage := int((float64(result.StorageBytes.Available) / float64(result.StorageBytes.Total)) * 100)
	storage := fmt.Sprintf("%d/%d (%d%%)", result.StorageBytes.Available, result.StorageBytes.Total, storagePercentage)

	r = append(r, []string{cpu, mem, storage})
	return &ResultTable{r}
}

func FromInstallResponse(result *grpc_infrastructure_manager_go.InstallResponse) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"ID", "STATE", "ERROR"})
	r = append(r, []string{result.ClusterId, result.State.String(), result.Error})
	return &ResultTable{r}
}

// ----
// Nodes
// ----

func FromNode(result *grpc_public_api_go.Node, labelLength int) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"ID", "IP", "STATE", "LABELS", "STATUS"})
	r = append(r, []string{result.NodeId, result.Ip, result.StateName, TransformLabels(result.Labels, labelLength), result.StatusName})
	return &ResultTable{r}
}

func FromNodeList(result *grpc_public_api_go.NodeList, labelLength int) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"ID", "IP", "STATE", "LABELS", "STATUS"})
	for _, n := range result.Nodes {
		r = append(r, []string{n.NodeId, n.Ip, n.StateName, TransformLabels(n.Labels, labelLength), n.StatusName})
	}
	return &ResultTable{r}
}

// ----
// Applications
// ----

func FromAppParameterList(result *grpc_public_api_go.AppParameterList) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"PARAM NAME", "DESCRIPTION", "PATH", "TYPE", "DEFAULT_VALUE", "B/A"})
	for _, p := range result.Parameters {
		r = append(r, []string{p.Name, p.Description, p.Path, p.Type, p.DefaultValue, p.Category})
	}
	return &ResultTable{r}
}

func FromInstanceParameterList(result *grpc_application_go.InstanceParameterList) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"PARAM NAME", "VALUE"})
	for _, p := range result.Parameters {
		r = append(r, []string{p.ParameterName, p.Value})
	}
	return &ResultTable{r}
}

func FromAppInstanceList(result *grpc_public_api_go.AppInstanceList, labelLength int) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"NAME", "ID", "LABELS", "STATUS"})
	for _, i := range result.Instances {
		r = append(r, []string{i.Name, i.AppInstanceId, TransformLabels(i.Labels, labelLength), i.StatusName})
	}
	return &ResultTable{r}
}

func FromAppInstance(result *grpc_public_api_go.AppInstance, labelLength int) *ResultTable {
	r := make([][]string, 0)

	r = append(r, []string{"NAME", "LABELS"})
	r = append(r, []string{result.Name, TransformLabels(result.Labels, labelLength)})
	r = append(r, []string{""})

	if result.StatusName == grpc_application_go.ApplicationStatus_DEPLOYMENT_ERROR.String() {
		r = append(r, []string{"STATUS", "INFO"})
		r = append(r, []string{result.StatusName, result.Info})
	} else {
		r = append(r, []string{"SERVICE_NAME", "REPLICAS", "STATUS", "ENDPOINTS"})
		for _, g := range result.Groups {
			groupReplicas := "NA"
			if g.Specs != nil {
				groupReplicas = strconv.Itoa(int(g.Specs.Replicas))
				if g.Specs.MultiClusterReplica {
					groupReplicas = "MULTI_CLUSTER"
				}
			}

			r = append(r, []string{fmt.Sprintf("[Group] %s", g.Name), groupReplicas, g.StatusName, strings.Join(g.GlobalFqdn, ", ")})
			for _, s := range g.ServiceInstances {
				r = append(r, []string{s.Name, strconv.Itoa(int(s.Specs.Replicas)), s.StatusName, strings.Join(s.Endpoints, ", ")})
			}
		}
		r = append(r, []string{"", "", "", ""})
		if (result.OutboundConnections != nil && len(result.OutboundConnections) > 0) ||
			(result.InboundConnections != nil && len(result.InboundConnections) > 0) {
			r = append(r, []string{"SOURCE", "OUTBOUND", "TARGET", "INBOUND", "REQUIRED"})
			if result.OutboundConnections != nil && len(result.OutboundConnections) > 0 {
				for _, out := range result.OutboundConnections {
					required := "FALSE"
					if out.OutboundRequired {
						required = "TRUE"
					}
					r = append(r, []string{out.SourceInstanceName, out.OutboundName, out.TargetInstanceName, out.InboundName, required})
				}
			}
			if result.InboundConnections != nil && len(result.InboundConnections) > 0 {
				for _, in := range result.InboundConnections {
					required := "FALSE"
					if in.OutboundRequired {
						required = "TRUE"
					}
					r = append(r, []string{in.SourceInstanceName, in.OutboundName, in.TargetInstanceName, in.InboundName, required})
				}
			}
		}
	}
	return &ResultTable{r}
}

func FromDeploymentResponse(result *grpc_application_manager_go.DeploymentResponse) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"REQUEST", "ID", "STATUS"})
	r = append(r, []string{result.RequestId, result.AppInstanceId, result.Status.String()})
	return &ResultTable{r}
}

func FromAppDescriptorList(result *grpc_application_go.AppDescriptorList, labelLength int) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"NAME", "ID", "LABELS", "SERVICES"})
	for _, d := range result.Descriptors {
		numServices := 0
		for _, g := range d.Groups {
			numServices = numServices + len(g.Services)
		}
		r = append(r, []string{d.Name, d.AppDescriptorId, TransformLabels(d.Labels, labelLength), strconv.Itoa(numServices)})
	}
	return &ResultTable{r}
}

func FromAppDescriptor(result *grpc_application_go.AppDescriptor, labelLength int) *ResultTable {
	r := make([][]string, 0)

	r = append(r, []string{"DESCRIPTOR", "ID", "LABELS"})
	r = append(r, []string{result.Name, result.AppDescriptorId, TransformLabels(result.Labels, labelLength)})
	r = append(r, []string{"", "", ""})

	if len(result.Parameters) > 0 {
		r = append(r, []string{"PARAM NAME", "DESCRIPTION", "DEFAULT VALUE"})
		for _, p := range result.Parameters {
			r = append(r, []string{p.Name, p.Description, p.DefaultValue})
		}
		r = append(r, []string{"", "", ""})
	}

	r = append(r, []string{"NAME", "IMAGE", "LABELS"})
	for _, g := range result.Groups {
		r = append(r, []string{fmt.Sprintf("[Group] %s", g.Name), "===", TransformLabels(g.Labels, labelLength)})
		for _, s := range g.Services {
			r = append(r, []string{s.Name, s.Image, TransformLabels(s.Labels, labelLength)})
		}
	}

	return &ResultTable{r}
}

// ----
// Devices
// ----

func FromDeviceGroup(result *grpc_device_manager_go.DeviceGroup) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"ID", "NAME", "API_KEY", "ENABLED", "DEV_ENABLED"})
	r = append(r, []string{result.DeviceGroupId, result.Name, result.DeviceGroupApiKey, strconv.FormatBool(result.Enabled), strconv.FormatBool(result.DefaultDeviceConnectivity)})
	return &ResultTable{r}
}

func FromDeviceGroupList(result *grpc_device_manager_go.DeviceGroupList) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"ID", "NAME", "API_KEY", "ENABLED", "DEV_ENABLED"})

	for _, dg := range result.Groups {
		r = append(r, []string{dg.DeviceGroupId, dg.Name, dg.DeviceGroupApiKey, strconv.FormatBool(dg.Enabled), strconv.FormatBool(dg.DefaultDeviceConnectivity)})
	}

	return &ResultTable{r}
}

func FromDevice(result *grpc_public_api_go.Device, labelLength int) *ResultTable {
	id := result.AssetDeviceId
	if id == "" {
		id = result.DeviceId
	}
	r := make([][]string, 0)
	r = append(r, []string{"ID", "DATE", "STATUS", "LABELS", "ENABLED"})
	r = append(r, []string{id, time.Unix(result.RegisterSince, 0).String(), result.DeviceStatusName, TransformLabels(result.Labels, labelLength), strconv.FormatBool(result.Enabled)})
	r = append(r, []string{""})
	r = append(r, []string{"GEOLOCATION"})
	location := "NA"
	if result.Location != nil {
		location = result.Location.Geolocation
	}
	r = append(r, []string{location})
	r = append(r, []string{""})
	r = append(r, []string{"OS", "CPUS", "RAM", "STORAGE"})
	os := "NA"
	if result.AssetInfo != nil && result.AssetInfo.Os != nil && len(result.AssetInfo.Os.Name) > 0 {
		os = result.AssetInfo.Os.Name
	}
	cpus := "NA"
	ram := "NA"
	if result.AssetInfo != nil && result.AssetInfo.Hardware != nil {
		count := 0
		for _, cpu := range result.AssetInfo.Hardware.Cpus {
			count = count + int(cpu.NumCores)
		}
		cpus = fmt.Sprintf("%d", count)
		ram = fmt.Sprintf("%d", result.AssetInfo.Hardware.InstalledRam)
	}
	storage := "NA"
	if result.AssetInfo != nil && result.AssetInfo.Storage != nil && len(result.AssetInfo.Storage) > 0 {
		var total int64 = 0
		for _, storage := range result.AssetInfo.Storage {
			total = total + storage.TotalCapacity
		}
		storage = fmt.Sprintf("%d", total)
	}
	r = append(r, []string{os, cpus, ram, storage})

	return &ResultTable{r}
}

func FromDeviceList(result *grpc_public_api_go.DeviceList, labelLength int) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"ID", "DATE", "STATUS", "LABELS", "ENABLED"})

	for _, d := range result.Devices {
		r = append(r, []string{d.DeviceId, time.Unix(d.RegisterSince, 0).String(), d.DeviceStatusName, TransformLabels(d.Labels, labelLength), strconv.FormatBool(d.Enabled)})
	}

	return &ResultTable{r}
}

// ----
// Log
// ----

func FromLogResponse(result *grpc_unified_logging_go.LogResponse) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"TIMESTAMP", "MSG"})

	for _, e := range result.Entries {
		r = append(r, []string{time.Unix(e.Timestamp.Seconds, int64(e.Timestamp.Nanos)).String(), e.Msg})
	}

	return &ResultTable{r}
}

// ----
// Roles
// ----

func FromRole(result *grpc_public_api_go.Role) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"ID", "NAME", "PRIMITIVES"})
	r = append(r, []string{result.RoleId, result.Name, strings.Join(result.Primitives, ",")})
	return &ResultTable{r}
}

func FromRoleList(result *grpc_public_api_go.RoleList) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"ID", "NAME", "PRIMITIVES"})

	for _, role := range result.Roles {
		r = append(r, []string{role.RoleId, role.Name, strings.Join(role.Primitives, ",")})
	}

	return &ResultTable{r}
}

// ----
// EdgeController
// ----

func FromEICJoinToken(result *grpc_inventory_manager_go.EICJoinToken) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"TOKEN", "EXPIRES"})
	r = append(r, []string{result.Token, time.Unix(result.ExpiresOn, 0).String()})
	return &ResultTable{r}
}

func FromIEdgeController(result *grpc_inventory_go.EdgeController, labelLength int) *ResultTable {
	r := make([][]string, 0)
	name := "NA"
	if result.Name != "" {
		name = result.Name
	}
	r = append(r, []string{"NAME", "ID", "CREATED"})
	r = append(r, []string{name, result.EdgeControllerId, time.Unix(result.Created, 0).String()})
	r = append(r, []string{""})
	r = append(r, []string{"NAME", "GEOLOCATION", "LABELS"})

	location := "NA"
	if result.Location != nil {
		location = result.Location.Geolocation
	}
	labels := "NA"
	if result.Labels != nil {
		labels = TransformLabels(result.Labels, labelLength)
	}
	r = append(r, []string{name, location, labels})

	// Asset Info
	r = append(r, []string{""})
	r = append(r, []string{"OS", "CPUS", "RAM", "STORAGE"})
	os := "NA"
	if result.AssetInfo != nil && result.AssetInfo.Os != nil && len(result.AssetInfo.Os.Name) > 0 {
		os = result.AssetInfo.Os.Name
	}
	cpus := "NA"
	ram := "NA"
	if result.AssetInfo != nil && result.AssetInfo.Hardware != nil {
		count := 0
		for _, cpu := range result.AssetInfo.Hardware.Cpus {
			count = count + int(cpu.NumCores)
		}
		cpus = fmt.Sprintf("%d", count)
		ram = fmt.Sprintf("%d", result.AssetInfo.Hardware.InstalledRam)
	}
	storage := "NA"
	if result.AssetInfo != nil && result.AssetInfo.Storage != nil && len(result.AssetInfo.Storage) > 0 {
		var total int64 = 0
		for _, storage := range result.AssetInfo.Storage {
			total = total + storage.TotalCapacity
		}
		storage = fmt.Sprintf("%d", total)
	}
	r = append(r, []string{os, cpus, ram, storage})

	if result.LastOpResult != nil {
		r = append(r, []string{""})
		r = append(r, []string{"LAST OP"})
		r = append(r, []string{"OP_ID", "TIMESTAMP", "STATUS", "INFO"})
		r = append(r, []string{result.LastOpResult.OperationId, time.Unix(result.LastOpResult.Timestamp, 0).String(),
			result.LastOpResult.Status.String(), result.LastOpResult.Info})
	}

	return &ResultTable{r}
}

func FromEdgeControllerExtendedInfo(result *grpc_public_api_go.EdgeControllerExtendedInfo, labelLength int) *ResultTable {
	r := make([][]string, 0)
	geolocation := ""
	if result.Controller != nil {
		geolocation = result.Controller.Location.Geolocation
	}

	if result.Controller != nil {
		r = append(r, []string{"NAME", "LABELS", "LOCATION", "STATUS", "SEEN"})
		seen := "never"
		if result.Controller.LastAliveTimestamp != 0 {
			seen = time.Unix(result.Controller.LastAliveTimestamp, 0).String()
		}
		r = append(r, []string{result.Controller.Name, TransformLabels(result.Controller.Labels, labelLength), geolocation, result.Controller.StatusName, seen})
	}
	// Asset Info
	r = append(r, []string{""})
	r = append(r, []string{"OS", "CPUS", "RAM", "STORAGE"})
	os := "NA"
	if result.Controller.AssetInfo != nil && result.Controller.AssetInfo.Os != nil && len(result.Controller.AssetInfo.Os.Name) > 0 {
		os = result.Controller.AssetInfo.Os.Name
	}
	cpus := "NA"
	ram := "NA"
	if result.Controller.AssetInfo != nil && result.Controller.AssetInfo.Hardware != nil {
		count := 0
		for _, cpu := range result.Controller.AssetInfo.Hardware.Cpus {
			count = count + int(cpu.NumCores)
		}
		cpus = fmt.Sprintf("%d", count)
		ram = fmt.Sprintf("%d", result.Controller.AssetInfo.Hardware.InstalledRam)
	}
	storage := "NA"
	if result.Controller.AssetInfo != nil && result.Controller.AssetInfo.Storage != nil && len(result.Controller.AssetInfo.Storage) > 0 {
		var total int64 = 0
		for _, storage := range result.Controller.AssetInfo.Storage {
			total = total + storage.TotalCapacity
		}
		storage = fmt.Sprintf("%d", total)
	}
	r = append(r, []string{os, cpus, ram, storage})

	if result.Controller.LastOpResult != nil {
		r = append(r, []string{""})
		r = append(r, []string{"LAST OP"})
		r = append(r, []string{"OP_ID", "TIMESTAMP", "STATUS", "INFO"})
		r = append(r, []string{result.Controller.LastOpResult.OperationId, time.Unix(result.Controller.LastOpResult.Timestamp, 0).String(),
			result.Controller.LastOpResult.OpStatusName, result.Controller.LastOpResult.Info})
	}

	// Managed Assets
	if len(result.ManagedAssets) > 0 {
		r = append(r, []string{""})
		r = append(r, []string{"ASSET_ID", "IP", "STATUS", "SEEN"})
		for _, a := range result.ManagedAssets {
			seen := "never"
			if a.LastAliveTimestamp != 0 {
				seen = time.Unix(a.LastAliveTimestamp, 0).String()
			}
			r = append(r, []string{a.AssetId, a.EicNetIp, a.StatusName, seen})
		}
	}
	return &ResultTable{r}
}

// ----
// Agent
// -----

func FromECOpResponse(result *grpc_public_api_go.ECOpResponse) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"OPERATION"})
	r = append(r, []string{result.OperationId})
	return &ResultTable{r}
}

func FromAgentJoinToken(result *grpc_inventory_manager_go.AgentJoinToken) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"TOKEN", "EXPIRES"})
	r = append(r, []string{result.Token, time.Unix(result.ExpiresOn, 0).String()})
	return &ResultTable{r}
}

// ----
// Inventory
// ----

func FromInventoryList(result *grpc_public_api_go.InventoryList, labelLength int) *ResultTable {
	r := make([][]string, 0)

	r = append(r, []string{"TYPE", "ID", "LOCATION", "LABELS", "STATUS"})
	for _, device := range result.Devices {
		geolocation := "NA"
		if device.Location != nil {
			geolocation = device.Location.Geolocation
		}
		r = append(r, []string{"DEVICE", device.AssetDeviceId, geolocation, TransformLabels(device.Labels, labelLength), device.DeviceStatusName})
	}
	for _, ec := range result.Controllers {
		geolocation := "NA"
		if ec.Location != nil {
			geolocation = ec.Location.Geolocation
		}
		r = append(r, []string{"EC", ec.EdgeControllerId, geolocation, TransformLabels(ec.Labels, labelLength), ec.StatusName})
	}
	for _, asset := range result.Assets {
		geolocation := "NA"
		if asset.Location != nil {
			geolocation = asset.Location.Geolocation
		}
		r = append(r, []string{"ASSET", asset.AssetId, geolocation, TransformLabels(asset.Labels, labelLength), asset.StatusName})
	}
	return &ResultTable{r}
}

func FromInventorySummary(result *grpc_inventory_manager_go.InventorySummary) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"CPUs", "STORAGE (GB)", "RAM (GB)"})
	r = append(r, []string{strconv.FormatInt(result.TotalNumCpu, 10), strconv.FormatInt(result.TotalStorage, 10), strconv.FormatInt(result.TotalRam, 10)})

	return &ResultTable{r}
}

// ----
// Application Network
// ----
func FromAvailableInboundList(result *grpc_application_manager_go.AvailableInstanceInboundList) *ResultTable {
	r := make([][]string, len(result.InstanceInbounds))
	r = append(r, []string{"INSTANCE_ID", "INSTANCE_NAME", "INBOUND_NAME"})

	for _, inbound := range result.InstanceInbounds {
		r = append(r, []string{inbound.AppInstanceId, inbound.InstanceName, inbound.InboundName})
	}

	return &ResultTable{r}
}

func FromAvailableOutboundList(result *grpc_application_manager_go.AvailableInstanceOutboundList) *ResultTable {
	r := make([][]string, len(result.InstanceOutbounds))
	r = append(r, []string{"INSTANCE_ID", "INSTANCE_NAME", "OUTBOUND_NAME"})

	for _, outbound := range result.InstanceOutbounds {
		r = append(r, []string{outbound.AppInstanceId, outbound.InstanceName, outbound.OutboundName})
	}

	return &ResultTable{r}
}

func FromConnectionInstanceListResult(result *grpc_application_network_go.ConnectionInstanceList) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"SOURCE_INSTANCE_ID", "SOURCE_INSTANCE_NAME", "OUTBOUND", "TARGET_INSTANCE_ID", "TARGET_INSTANCE_NAME", "INBOUND"})
	for _, connection := range result.Connections {
		r = append(r, []string{connection.SourceInstanceId, connection.SourceInstanceName, connection.OutboundName, connection.TargetInstanceId,
			connection.TargetInstanceName, connection.InboundName})
	}

	return &ResultTable{r}
}

// ----
// Inventory monitoring
// ----

func FromQueryMetricsResult(result *grpc_monitoring_go.QueryMetricsResult) *ResultTable {
	r := [][]string{}
	r = append(r, []string{"TIMESTAMP", "METRIC", "ASSET", "AGGR", "VALUE"})

	for metric, assetMetric := range result.GetMetrics() {
		for _, metrics := range assetMetric.GetMetrics() {
			for _, value := range metrics.GetValues() {
				timestamp := time.Unix(value.GetTimestamp(), 0).Local().String()
				r = append(r, []string{timestamp, metric, metrics.GetAssetId(), metrics.GetAggregation().String(), strconv.FormatInt(value.GetValue(), 10)})
			}
		}
	}

	return &ResultTable{r}
}

func FromMetricsList(result *grpc_monitoring_go.MetricsList) *ResultTable {
	r := [][]string{}
	r = append(r, []string{"METRIC"})

	for _, metric := range result.GetMetrics() {
		r = append(r, []string{metric})
	}

	return &ResultTable{r}
}

// ----
// Assets
// ----

func FromAsset(result *grpc_public_api_go.Asset) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"ID", "CONTROLLER", "AGENT"})
	r = append(r, []string{result.AssetId, result.EdgeControllerId, result.AgentId})
	r = append(r, []string{""})
	r = append(r, []string{"IP", "SEEN", "STATUS"})
	r = append(r, []string{result.EicNetIp, time.Unix(result.LastAliveTimestamp, 0).String(), result.StatusName})
	r = append(r, []string{""})
	r = append(r, []string{"GEOLOCATION"})
	location := "NA"
	if result.Location != nil {
		location = result.Location.Geolocation
	}
	r = append(r, []string{location})
	r = append(r, []string{""})
	r = append(r, []string{"OS", "CPUS", "RAM", "STORAGE"})
	os := "NA"
	if result.Os != nil && len(result.Os.Name) > 0 {
		os = result.Os.Name
	}
	cpus := "NA"
	ram := "NA"
	if result.Hardware != nil {
		count := 0
		for _, cpu := range result.Hardware.Cpus {
			count = count + int(cpu.NumCores)
		}
		cpus = fmt.Sprintf("%d", count)
		ram = fmt.Sprintf("%d", result.Hardware.InstalledRam)
	}
	storage := "NA"
	if result.Storage != nil && len(result.Storage) > 0 {
		var total int64 = 0
		for _, storage := range result.Storage {
			total = total + storage.TotalCapacity
		}
		storage = fmt.Sprintf("%d", total)
	}
	r = append(r, []string{os, cpus, ram, storage})

	if result.LastOpSummary != nil {
		r = append(r, []string{""})
		r = append(r, []string{"LAST OP"})
		r = append(r, []string{"OP_ID", "TIMESTAMP", "STATUS", "INFO"})
		r = append(r, []string{result.LastOpSummary.OperationId, time.Unix(result.LastOpSummary.Timestamp, 0).String(), result.LastOpSummary.OpStatusName, result.LastOpSummary.Info})
	}

	return &ResultTable{r}
}

func FromIAsset(result *grpc_inventory_go.Asset, labelLength int) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"ID", "CONTROLLER", "AGENT"})
	r = append(r, []string{result.AssetId, result.EdgeControllerId, result.AgentId})
	r = append(r, []string{""})
	r = append(r, []string{"IP", "SEEN"})
	r = append(r, []string{result.EicNetIp, time.Unix(result.LastAliveTimestamp, 0).String()})
	r = append(r, []string{""})
	r = append(r, []string{"GEOLOCATION", "LABELS"})
	location := "NA"
	if result.Location != nil {
		location = result.Location.Geolocation
	}
	labels := "NA"
	if result.Labels != nil {
		labels = TransformLabels(result.Labels, labelLength)
	}
	r = append(r, []string{location, labels})
	r = append(r, []string{""})
	r = append(r, []string{"OS", "CPUS", "RAM", "STORAGE"})
	os := "NA"
	if result.Os != nil && len(result.Os.Name) > 0 {
		os = result.Os.Name
	}
	cpus := "NA"
	ram := "NA"
	if result.Hardware != nil {
		count := 0
		for _, cpu := range result.Hardware.Cpus {
			count = count + int(cpu.NumCores)
		}
		cpus = fmt.Sprintf("%d", count)
		ram = fmt.Sprintf("%d", result.Hardware.InstalledRam)
	}
	storage := "NA"
	if result.Storage != nil && len(result.Storage) > 0 {
		var total int64 = 0
		for _, storage := range result.Storage {
			total = total + storage.TotalCapacity
		}
		storage = fmt.Sprintf("%d", total)
	}
	r = append(r, []string{os, cpus, ram, storage})

	return &ResultTable{r}
}

func FromAgentOpResponse(result *grpc_public_api_go.AgentOpResponse) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"OPERATION_ID", "TIMESTAMP", "STATUS", "INFO"})
	r = append(r, []string{result.OperationId, time.Unix(result.Timestamp, 0).String(), result.Status, result.Info})
	return &ResultTable{r}
}

// ----
// Common
// ----

func FromSuccess(result *grpc_common_go.Success) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"RESULT"})
	r = append(r, []string{"OK"})
	return &ResultTable{r}
}

func FromOpResponse(result *grpc_public_api_go.OpResponse) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"REQUEST_ID", "TIMESTAMP", "STATUS", "INFO"})
	r = append(r, []string{result.RequestId, time.Unix(result.Timestamp, 0).String(), result.StatusName, result.Info})
	return &ResultTable{r}
}
