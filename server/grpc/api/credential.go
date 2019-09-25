package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GetCredentialFromContext retrives credential name and password from request metadata.
func GetCredentialFromContext(ctx context.Context) (string, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	credentialUsernameMd := md["credential_username"]
	if len(credentialUsernameMd) == 0 {
		return "", "", status.Errorf(codes.Unauthenticated, "missing credential_username from metadata")
	}

	credentialPassphraseMd := md["credential_passphrase"]
	if len(credentialPassphraseMd) == 0 {
		return "", "", status.Errorf(codes.Unauthenticated, "missing credential_passphrase from metadata")
	}

	return credentialUsernameMd[0], credentialPassphraseMd[0], nil
}
