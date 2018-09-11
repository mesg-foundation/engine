package service

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/core/utils/chunker"
	"github.com/mesg-foundation/core/x/xcolor"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/mesg-foundation/core/x/xstrings"
	"github.com/mesg-foundation/prefixer"
	"github.com/spf13/cobra"
)

// Logs of the core
var Logs = &cobra.Command{
	Use:   "logs",
	Short: "Show the logs of a service",
	Example: `mesg-core service logs SERVICE_ID
mesg-core service logs SERVICE_ID --dependency DEPENDENCY_NAME`,
	Run:               logsHandler,
	Args:              cobra.MinimumNArgs(1),
	DisableAutoGenTag: true,
}

var dependencies []string

func init() {
	Logs.Flags().StringArrayVarP(&dependencies, "dependency", "d", nil, "Name of the dependency to only show the logs from")
}

func logsHandler(cmd *cobra.Command, args []string) {
	closeReaders := showLogs(args[0], dependencies...)
	defer closeReaders()
	<-xsignal.WaitForInterrupt()
}

func showLogs(serviceID string, dependencies ...string) func() {
	prefixes, err := dependencyPrefixes(serviceID)
	utils.HandleError(err)

	ctx, cancel := context.WithCancel(context.Background())

	stream, err := cli().ServiceLogs(ctx, &core.ServiceLogsRequest{
		ServiceID:    serviceID,
		Dependencies: dependencies,
	})
	utils.HandleError(err)

	streams := make(map[chunkMeta]*chunker.Stream)

	go func() {
		for {
			data, err := stream.Recv()
			if err != nil {
				return
			}
			meta := chunkMeta{
				Dependency: data.Dependency,
				Type:       data.Type,
			}
			stream, ok := streams[meta]
			if !ok {
				stream = chunker.NewStream()
				go prefixedCopy(os.Stdout, stream, prefixes[meta.Dependency])
				streams[meta] = stream
			}
			stream.Provide(data.Data)
		}
	}()

	return func() {
		cancel()
		for _, stream := range streams {
			stream.Close()
		}
	}
}

// dependencyPrefixes returns dependency key, log prefix pair.
func dependencyPrefixes(serviceID string) (prefixes map[string]string, err error) {
	// maxCharLen is the char length of longest dependency key.
	var maxCharLen int

	// get list of services to calibrate spaces for short dependency keys.
	resp, err := cli().GetService(context.Background(), &core.GetServiceRequest{
		ServiceID: serviceID,
	})
	if err != nil {
		return nil, err
	}

	for key := range resp.Service.Dependencies {
		l := len(key)
		if l > maxCharLen {
			maxCharLen = l
		}
	}
	for key := range resp.Service.Dependencies {
		prefixes[key] = color.New(xcolor.NextColor()).Sprintf("%s |", fillSpace(key, maxCharLen))
	}

	return prefixes, nil
}

// prefixedReader wraps io.Reader by adding a prefix for each new line
// in the stream.
func prefixedReader(r io.Reader, prefix string) io.Reader {
	return prefixer.New(r, fmt.Sprintf("%s ", prefix))
}

// prefixedCopy copies src to dst by prefixing dependency key to each new line.
func prefixedCopy(dst io.Writer, src io.Reader, dep string) {
	io.Copy(dst, prefixedReader(src, dep))
}

// fillSpace fills the end of name with spaces until max chars limit hits.
func fillSpace(name string, max int) string {
	return xstrings.AppendSpaces(name, max-len(name))
}

// chunkMeta is a meta data for chunks.
type chunkMeta struct {
	Dependency string
	Type       core.LogData_Type
}
