package authhelper

import (
	"context"
	"github.com/nalej/derrors"
	"google.golang.org/grpc/metadata"
)

type RequestMetadata struct {
	UserID         string
	OrganizationID string
}

func GetRequestMetadata(ctx context.Context) (*RequestMetadata, derrors.Error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, derrors.NewInvalidArgumentError("expecting JWT metadata")
	}
	//log.Debug().Interface("metadata", md).Msg("Metadata received")
	userID, found := md["user_id"]
	if !found {
		return nil, derrors.NewUnauthenticatedError("userID not found")
	}
	organizationID, found := md["organization_id"]
	if !found {
		return nil, derrors.NewUnauthenticatedError("organizationID not found")
	}
	return &RequestMetadata{userID[0], organizationID[0]}, nil
}
