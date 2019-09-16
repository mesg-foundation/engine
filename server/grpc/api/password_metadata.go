package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GetAccountInfoFromContext retrives account name nad passwrod from request metadata.
func GetAccountInfoFromContext(ctx context.Context) (string, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	accountNameMd := md["account_name"]
	if len(accountNameMd) == 0 {
		return "", "", status.Errorf(codes.Unauthenticated, "missing account name")
	}

	accountPasswordMd := md["account_password"]
	if len(accountPasswordMd) == 0 {
		return "", "", status.Errorf(codes.Unauthenticated, "missing account password")
	}

	return accountNameMd[0], accountPasswordMd[0], nil
}
