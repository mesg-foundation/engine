package acknowledgement

import (
	fmt "fmt"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	statusKey   = "status"
	statusReady = "ready"
)

// WaitForStreamToBeReady waits until the server publish a header with status ready on the stream.
func WaitForStreamToBeReady(stream grpc.ClientStream) error {
	// wait for header.
	header, err := stream.Header()
	if err != nil {
		return err
	}
	// header received. check status.
	statuses := header.Get(statusKey)
	statusesLen := len(statuses)
	if statusesLen == 0 {
		return nil // Ignore headers with no status key
	}
	lastStatus := statuses[statusesLen-1]
	if lastStatus != statusReady {
		return fmt.Errorf("stream header status is different than ready. Got %q", lastStatus)
	}
	return nil
}

// SetStreamReady send an header on the stream to notify the server is ready.
func SetStreamReady(stream grpc.ServerStream) error {
	return stream.SendHeader(metadata.Pairs(statusKey, statusReady))
}
