package utils

import "google.golang.org/grpc/metadata"

// StatusReady is a metadata to notify a stream is ready.
var StatusReady = metadata.Pairs("mesg-status", "ready")
