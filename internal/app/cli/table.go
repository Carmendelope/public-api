package cli

import (
	"fmt"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-conductor-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-infrastructure-manager-go"
	"github.com/nalej/grpc-infrastructure-monitor-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const MinWidth = 5
const TabWidth = 2
const Padding = 3

const AppInstanceHeader = ""

var Headers = map[string]string{

}

type ResultTable struct {
	data [][]string
}

func AsTable(result interface{}) * ResultTable {
	switch result.(type){
	case *grpc_public_api_go.OrganizationInfo: return FromOrganizationInfo(result.(*grpc_public_api_go.OrganizationInfo))
	case *grpc_public_api_go.User: return FromUser(result.(*grpc_public_api_go.User))
	case *grpc_user_manager_go.User: return FromUserManagerUser(result.(*grpc_user_manager_go.User))
	case *grpc_public_api_go.UserList: return FromUserList(result.(*grpc_public_api_go.UserList))
	case *grpc_public_api_go.Cluster: return FromCluster(result.(*grpc_public_api_go.Cluster))
	case *grpc_infrastructure_monitor_go.ClusterSummary: return FromClusterSummary(result.(*grpc_infrastructure_monitor_go.ClusterSummary))
	case *grpc_public_api_go.ClusterList: return FromClusterList(result.(*grpc_public_api_go.ClusterList))
	case *grpc_infrastructure_manager_go.InstallResponse: return FromInstallResponse(result.(*grpc_infrastructure_manager_go.InstallResponse))
	case *grpc_public_api_go.AppInstanceList: return FromAppInstanceList(result.(*grpc_public_api_go.AppInstanceList))
	case *grpc_public_api_go.AppInstance: return FromAppInstance(result.(*grpc_public_api_go.AppInstance))
	case *grpc_conductor_go.DeploymentResponse: return FromDeploymentResponse(result.(*grpc_conductor_go.DeploymentResponse))
	case *grpc_application_go.AppDescriptorList: return FromAppDescriptorList(result.(*grpc_application_go.AppDescriptorList))
	case *grpc_application_go.AppDescriptor: return FromAppDescriptor(result.(*grpc_application_go.AppDescriptor))
	case *grpc_device_manager_go.DeviceGroup: return FromDeviceGroup(result.(*grpc_device_manager_go.DeviceGroup))
	case *grpc_device_manager_go.DeviceGroupList: return FromDeviceGroupList(result.(*grpc_device_manager_go.DeviceGroupList))
	case *grpc_public_api_go.Device: return FromDevice(result.(*grpc_public_api_go.Device))
	case *grpc_public_api_go.DeviceList: return FromDeviceList(result.(*grpc_public_api_go.DeviceList))
	case *grpc_unified_logging_go.LogResponse: return FromLogResponse(result.(*grpc_unified_logging_go.LogResponse))
	case *grpc_public_api_go.Node: return FromNode(result.(*grpc_public_api_go.Node))
	case *grpc_public_api_go.NodeList: return FromNodeList(result.(*grpc_public_api_go.NodeList))
	case *grpc_public_api_go.Role: return FromRole(result.(*grpc_public_api_go.Role))
	case *grpc_public_api_go.RoleList: return FromRoleList(result.(*grpc_public_api_go.RoleList))
	case *grpc_common_go.Success: return FromSuccess(result.(*grpc_common_go.Success))
	default: log.Fatal().Str("type", fmt.Sprintf("%T", result)).Msg("unsupported")
	}
	return nil
}

func (t * ResultTable) Print() {
	w := tabwriter.NewWriter(os.Stdout, MinWidth, TabWidth, Padding, ' ', 0)
	for _, d := range t.data{
		toPrint := strings.Join(d, "\t")
		fmt.Fprintln(w, toPrint)
	}
	w.Flush()
}

func PrintFromValues(header []string, values [][]string) {
	w := tabwriter.NewWriter(os.Stdout, MinWidth, TabWidth, Padding, ' ', 0)
	fmt.Fprintln(w,strings.Join(header, "\t"))
	for _, d := range values{
		toPrint := strings.Join(d, "\t")
		fmt.Fprintln(w, toPrint)
	}
	w.Flush()
}

func TransformLabels(labels map[string]string) string {
	r := make([]string, 0)
	for k, v := range labels{
		r = append(r, fmt.Sprintf("%s:%s", k, v))
	}
	return strings.Join(r, ",")
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
	for _, u := range user.Users{
		result = append(result, []string{u.Name, u.RoleName, u.Email})
	}
	return &ResultTable{result}
}

// ----
// Clusters
// ----

func FromCluster(result *grpc_public_api_go.Cluster) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"NAME", "ID", "NODES", "LABELS", "STATUS"})
	r = append(r, []string{result.Name, result.ClusterId, fmt.Sprintf("%d", result.TotalNodes), TransformLabels(result.Labels), result.StatusName})
	return &ResultTable{r}
}

func FromClusterList(result *grpc_public_api_go.ClusterList) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"NAME", "ID", "NODES", "LABELS", "STATUS"})
	for _, c := range result.Clusters{
		r = append(r, []string{c.Name, c.ClusterId, fmt.Sprintf("%d", c.TotalNodes), TransformLabels(c.Labels), c.StatusName})
	}
	return &ResultTable{r}
}

func FromClusterSummary(result *grpc_infrastructure_monitor_go.ClusterSummary) *ResultTable {
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

func FromNode(result *grpc_public_api_go.Node) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"ID", "IP", "STATE", "LABELS", "STATUS"})
	r = append(r, []string{result.NodeId, result.Ip, result.StateName, TransformLabels(result.Labels), result.StatusName})
	return &ResultTable{r}
}

func FromNodeList(result *grpc_public_api_go.NodeList) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"ID", "IP", "STATE", "LABELS", "STATUS"})
	for _, n := range result.Nodes{
		r = append(r, []string{n.NodeId, n.Ip, n.StateName, TransformLabels(n.Labels), n.StatusName})
	}
	return &ResultTable{r}
}

// ----
// Applications
// ----
func FromAppInstanceList(result *grpc_public_api_go.AppInstanceList) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"NAME", "ID", "LABELS", "STATUS"})
	for _, i := range result.Instances{
		r = append(r, []string{i.Name, i.AppInstanceId, TransformLabels(i.Labels), i.StatusName})
	}
	return &ResultTable{r}
}

func FromAppInstance(result *grpc_public_api_go.AppInstance) *ResultTable {
	r := make([][]string, 0)
	if result.StatusName == grpc_application_go.ApplicationStatus_DEPLOYMENT_ERROR.String() {
		r = append(r, []string{"STATUS", "INFO"})
		r = append(r, []string{result.StatusName, result.Info})
	}else{
		r = append(r, []string{"NAME", "REPLICAS", "STATUS", "ENDPOINTS"})
		for _, g := range result.Groups{
			groupReplicas := "NA"
			if g.Specs != nil{
				groupReplicas = strconv.Itoa(int(g.Specs.Replicas))
				if g.Specs.MultiClusterReplica{
					groupReplicas = "MULTI_CLUSTER"
				}
			}

			r= append(r, []string{fmt.Sprintf("[Group] %s",g.Name),groupReplicas, g.StatusName, strings.Join(g.GlobalFqdn, ", ")})
			for _, s:= range g.ServiceInstances{
				r = append(r, []string{s.Name, strconv.Itoa(int(s.Specs.Replicas)), s.StatusName, strings.Join(s.Endpoints, ", ")})
			}
		}
	}
	return &ResultTable{r}
}

func FromDeploymentResponse(result *grpc_conductor_go.DeploymentResponse) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"REQUEST", "ID", "STATUS"})
	r = append(r, []string{result.RequestId, result.AppInstanceId, result.Status.String()})
	return &ResultTable{r}
}


func FromAppDescriptorList(result *grpc_application_go.AppDescriptorList) *ResultTable {
	r := make([][]string, 0)
	r = append(r, []string{"NAME", "ID", "LABELS", "SERVICES"})
	for _, d := range result.Descriptors{
		numServices := 0
		for _, g := range d.Groups{
			numServices = numServices + len(g.Services)
		}
		r = append(r, []string{d.Name, d.AppDescriptorId, TransformLabels(d.Labels), strconv.Itoa(numServices)})
	}
	return &ResultTable{r}
}

func FromAppDescriptor(result *grpc_application_go.AppDescriptor) *ResultTable {
	r := make([][]string, 0)

	r = append(r, []string{"DESCRIPTOR", "ID", "LABELS"})
	r = append(r, []string{result.Name, result.AppDescriptorId, TransformLabels(result.Labels)})
	r = append(r, []string{"", "", ""})

	r = append(r, []string{"NAME", "IMAGE", "LABELS"})
		for _, g := range result.Groups{
			r= append(r, []string{fmt.Sprintf("[Group] %s",g.Name), "===", TransformLabels(g.Labels)})
			for _, s:= range g.Services{
				r = append(r, []string{s.Name, s.Image, TransformLabels(s.Labels)})
			}
		}

	return &ResultTable{r}
}

// ----
// Devices
// ----

func FromDeviceGroup(result *grpc_device_manager_go.DeviceGroup) *ResultTable{
	r := make([][]string, 0)
	r = append(r, []string{"ID", "NAME", "API_KEY", "ENABLED", "DEV_ENABLED"})
	r = append(r, []string{result.DeviceGroupId, result.Name, result.DeviceGroupApiKey, strconv.FormatBool(result.Enabled), strconv.FormatBool(result.DefaultDeviceConnectivity)})
	return &ResultTable{r}
}

func FromDeviceGroupList(result *grpc_device_manager_go.DeviceGroupList) *ResultTable{
	r := make([][]string, 0)
	r = append(r, []string{"ID", "NAME", "API_KEY", "ENABLED", "DEV_ENABLED"})

	for _, dg := range result.Groups{
		r = append(r, []string{dg.DeviceGroupId, dg.Name, dg.DeviceGroupApiKey, strconv.FormatBool(dg.Enabled), strconv.FormatBool(dg.DefaultDeviceConnectivity)})
	}

	return &ResultTable{r}
}

func FromDevice(result *grpc_public_api_go.Device) *ResultTable{
	r := make([][]string, 0)
	r = append(r, []string{"ID", "DATE", "STATUS", "LABELS", "ENABLED"})
	r = append(r, []string{result.DeviceId, time.Unix(result.RegisterSince, 0).String(), result.DeviceStatusName, TransformLabels(result.Labels), strconv.FormatBool(result.Enabled)})
	return &ResultTable{r}
}

func FromDeviceList(result *grpc_public_api_go.DeviceList) *ResultTable{
	r := make([][]string, 0)
	r = append(r, []string{"ID", "DATE", "STATUS", "LABELS", "ENABLED"})

	for _, d := range result.Devices{
		r = append(r, []string{d.DeviceId, time.Unix(d.RegisterSince, 0).String(), d.DeviceStatusName, TransformLabels(d.Labels), strconv.FormatBool(d.Enabled)})
	}

	return &ResultTable{r}
}

// ----
// Log
// ----

func FromLogResponse(result *grpc_unified_logging_go.LogResponse) *ResultTable{
	r := make([][]string, 0)
	r = append(r, []string{"TIMESTAMP", "MSG"})

	for _, e := range result.Entries{
		r = append(r, []string{time.Unix(e.Timestamp.Seconds, int64(e.Timestamp.Nanos)).String(), e.Msg})
	}

	return &ResultTable{r}
}


// ----
// Roles
// ----

func FromRole(result *grpc_public_api_go.Role) *ResultTable{
	r := make([][]string, 0)
	r = append(r, []string{"ID", "NAME", "PRIMITIVES"})
	r = append(r, []string{result.RoleId, result.Name, strings.Join(result.Primitives, ",")})
	return &ResultTable{r}
}

func FromRoleList(result *grpc_public_api_go.RoleList) *ResultTable{
	r := make([][]string, 0)
	r = append(r, []string{"ID", "NAME", "PRIMITIVES"})

	for _, role := range result.Roles{
		r = append(r, []string{role.RoleId, role.Name, strings.Join(role.Primitives, ",")})
	}

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