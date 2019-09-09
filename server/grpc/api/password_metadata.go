package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GetAccountFromContext(ctx context.Context) (string, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	accountNameMd := md["account_name"]
	if len(accountNameMd) == 0 {
		return "", "", status.Errorf(codes.Unauthenticated, "invalid account name")
	}

	accountPasswordMd := md["account_password"]
	if len(accountPasswordMd) == 0 {
		return "", "", status.Errorf(codes.Unauthenticated, "invalid account password")
	}

	return accountNameMd[0], accountPasswordMd[0], nil
}
