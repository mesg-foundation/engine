package service

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/core/x/xsignal"
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

// dependencyLogs keeps dependency info and corresponding std & err log streams.
type dependencyLogs struct {
	Dependency      string
	Standard, Error *logReader
}

func showLogs(serviceID string, dependencies ...string) func() {
	stream, err := cli().ServiceLogs(context.Background(), &core.ServiceLogsRequest{
		ServiceID:    serviceID,
		Dependencies: dependencies,
	})
	utils.HandleError(err)

	// first received dependency list.
	data, err := stream.Recv()
	utils.HandleError(err)

	var (
		rstds, rerrs []*logReader
	)

	for _, dep := range data.Depedencies {
		var (
			rstd = newLogReader(dep, core.LogData_Data_Standard)
			rerr = newLogReader(dep, core.LogData_Data_Error)
		)

		rstds = append(rstds, rstd)
		rerrs = append(rerrs, rerr)
	}

	for _, r := range rstds {
		go prefixedCopy(os.Stdout, r, r.dependency)
	}
	for _, r := range rerrs {
		go prefixedCopy(os.Stderr, r, r.dependency)
	}

	for {
		data, err := stream.Recv()
		if err != nil {
			break
		}
		for _, l := range append(rstds, rerrs...) {
			l.process(data)
		}
	}

	return func() {
		for _, c := range append(rstds, rerrs...) {
			c.Close()
		}
	}
}

// logReader implements io.Reader to combine log data chunks being received
// from gRPC stream.
type logReader struct {
	dependency string
	typ        core.LogData_Data_Type

	recv chan []byte
	done chan struct{}

	data []byte
	i    int64
}

// newLogReader returns a new log reader.
func newLogReader(dependency string, typ core.LogData_Data_Type) *logReader {
	return &logReader{
		dependency: dependency,
		typ:        typ,
		recv:       make(chan []byte, 0),
		done:       make(chan struct{}, 0),
	}
}

// process processes log data received from gRPC stream and checks if it belongs
// to this log stream.
func (r *logReader) process(data *core.LogData) {
	if r.dependency == data.Data.Dependency &&
		r.typ == data.Data.Type {
		r.recv <- data.Data.Data
	}
}

// Read implements io.Reader.
func (r *logReader) Read(p []byte) (n int, err error) {
	if r.i >= int64(len(r.data)) {
		for {
			select {
			case <-r.done:
				return 0, io.EOF

			case data := <-r.recv:
				if err != nil {
					return 0, err
				}
				r.data = data
				r.i = 0
				return r.Read(p)
			}
		}
	}
	n = copy(p, r.data[r.i:])
	r.i += int64(n)
	return n, nil
}

// Close closes log reader.
func (r *logReader) Close() error {
	close(r.done)
	return nil
}

// prefixedReader wraps io.Reader by adding a prefix for each new line
// in the stream.
func prefixedReader(r io.Reader, prefix string) io.Reader {
	return prefixer.New(r, fmt.Sprintf("%s ", prefix))
}

// prefixedCopy copies src to dst by prefixing dependency key to each new line.
func prefixedCopy(dst io.Writer, src io.Reader, dep string) {
	prefix := color.New(randColor()).Sprintf("%s |", dep)
	io.Copy(dst, prefixedReader(src, prefix))
}

// randColor returns a random color.
func randColor() color.Attribute {
	attrs := []color.Attribute{
		color.FgRed,
		color.FgGreen,
		color.FgYellow,
		color.FgBlue,
		color.FgMagenta,
		color.FgCyan,
	}
	rand.Seed(time.Now().UnixNano())
	return attrs[rand.Intn(len(attrs))]
}
