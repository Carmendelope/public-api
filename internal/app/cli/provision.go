/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package cli

import (
    "fmt"
    "github.com/golang/protobuf/jsonpb"
    "github.com/google/uuid"
    "github.com/nalej/derrors"
    "github.com/nalej/grpc-common-go"
    "github.com/nalej/grpc-infrastructure-go"
    "github.com/nalej/grpc-installer-go"
    "github.com/nalej/grpc-provisioner-go"
    "github.com/nalej/grpc-public-api-go"
    "github.com/rs/zerolog/log"
    "os"
    "time"
)

const CliProvisionCheckSleepTime=time.Second * 20

type Provision struct {
    Connection
    Credentials
}

func NewProvision(address string, port int, insecure bool, useTLS bool, caCertPath string, output string,
    labelLength int) *Provision {
    return &Provision{
        Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
        Credentials: *NewEmptyCredentials(DefaultPath),
    }
}

func (p *Provision) Cluster(organizationId string, clusterName string, azureCredentialsPath string,
    azureDnsZoneName string, azureResourceGroup string, clusterType grpc_infrastructure_go.ClusterType, isManagementCluster bool,
    isProduction bool, kubernetesVersion string, nodeType string, numNodes int64, targetPlatform grpc_installer_go.Platform,
    zone string) {
    err := p.LoadCredentials()
    if err != nil {
        log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
    }

    azureCredentials, err := p.loadAzureCredentials(azureCredentialsPath)
    if err != nil {
        p.PrintResultOrError("", err, "error loading azure credentials")
        return
    }

    c, err := p.GetConnection()
    if err != nil {
        log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
    }
    defer c.Close()

    provClient := grpc_public_api_go.NewProvisionClient(c)



    request := grpc_provisioner_go.ProvisionClusterRequest{
        RequestId: uuid.New().String(),
        OrganizationId: organizationId,
        ClusterId: uuid.New().String(),
        ClusterName: clusterName,
        AzureCredentials: azureCredentials,
        AzureOptions: &grpc_provisioner_go.AzureProvisioningOptions{
            DnsZoneName: azureDnsZoneName,
            ResourceGroup: azureResourceGroup,
        },
        ClusterType: clusterType,
        IsManagementCluster: isManagementCluster,
        IsProduction: isProduction,
        KubernetesVersion: kubernetesVersion,
        NodeType: nodeType,
        NumNodes: numNodes,
        TargetPlatform: targetPlatform,
        Zone: zone,
    }

    ctx,cancel := p.GetContext()
    defer cancel()

    resp, errReq := provClient.ProvisionCluster(ctx, &request)

    p.PrintResultOrError(resp, errReq, "cannot provision cluster")
}

func (p *Provision) CheckProgress(requestId string) {
    err := p.LoadCredentials()
    if err != nil {
        log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
    }

    c, err := p.GetConnection()
    if err != nil {
        log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
    }
    defer c.Close()


    provClient := grpc_public_api_go.NewProvisionClient(c)
    // force initial check
    stop := p.checkCall(provClient, requestId)
    if stop {
        return
    }

    // Check periodically
    ticker := time.NewTicker(CliProvisionCheckSleepTime)

    for {
        select {
        case _ = <- ticker.C:
            stop := p.checkCall(provClient, requestId)
            if stop {
                return
            }
        }
    }
}

func (p *Provision) checkCall(client grpc_public_api_go.ProvisionClient, requestId string) bool{
    ctx,cancel := p.GetContext()
    resultCheck, errCheck := client.CheckProgress(ctx, &grpc_common_go.RequestId{RequestId: requestId})
    cancel()
    if errCheck != nil {
        p.PrintResultOrError(nil, errCheck, "error checking cluster provision")
        log.Info().Msgf("%+v\n",resultCheck)
        return true
    }
    p.printProgress(resultCheck)
    if resultCheck.State == grpc_provisioner_go.ProvisionProgress_FINISHED ||
        resultCheck.State == grpc_provisioner_go.ProvisionProgress_ERROR {
        // print the final message
        log.Info().Msgf("%+v\n",resultCheck)
        return true
    }
    // do not stop
    return false
}

func (p *Provision) printProgress(progress *grpc_provisioner_go.ProvisionClusterResponse) {
    date := time.Unix(progress.ElapsedTime,0)
    strDate := date.Format("20060102-150405")
    fmt.Printf("%s\t%s\t%s\t%s\n",strDate,progress.RequestId, progress.State, progress.Error)
}

func (p *Provision) RemoveProvision(requestId string) {
    err := p.LoadCredentials()
    if err != nil {
        log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
    }

    c, err := p.GetConnection()
    if err != nil {
        log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
    }
    defer c.Close()
    ctx,cancel := p.GetContext()
    defer cancel()

    provClient := grpc_public_api_go.NewProvisionClient(c)

    resp, errReq := provClient.RemoveProvision(ctx, &grpc_common_go.RequestId{RequestId: requestId})

    p.PrintResultOrError(resp, errReq, "cannot provision cluster")
}


// LoadAzureCredentials loads the content of a file into the grpc structure.
func (p *Provision) loadAzureCredentials(credentialsPath string) (*grpc_provisioner_go.AzureCredentials, derrors.Error) {
    credentials := &grpc_provisioner_go.AzureCredentials{}
    file, err := os.Open(credentialsPath)
    if err != nil {
        return nil, derrors.AsError(err, "cannot open credentials path")
    }
    // The unmarshalling using jsonpb is available due to the fact that the naming of the JSON fields produced
    // by Azure matches the ones defined in the protobuf json mapping.
    err = jsonpb.Unmarshal(file, credentials)
    if err != nil {
        return nil, derrors.AsError(err, "cannot unmarshal content")
    }
    log.Debug().Interface("tenantId", credentials.TenantId).Msg("azure credentials have been loaded")
    return credentials, nil
}