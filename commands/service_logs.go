package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/mesg-foundation/core/x/xstrings"
	"github.com/mesg-foundation/prefixer"
	"github.com/spf13/cobra"
)

type serviceLogsCmd struct {
	baseCmd

	dependencies []string

	e ServiceExecutor
}

func newServiceLogsCmd(e ServiceExecutor) *serviceLogsCmd {
	c := &serviceLogsCmd{
		e: e,
	}

	c.cmd = newCommand(&cobra.Command{
		Use:   "logs",
		Short: "Show logs of a service",
		Example: `mesg-core service logs SERVICE
mesg-core service logs SERVICE --dependencies DEPENDENCY_NAME,DEPENDENCY_NAME,...`,
		Args: cobra.ExactArgs(1),
		RunE: c.runE,
	})

	c.cmd.Flags().StringArrayVarP(&c.dependencies, "dependencies", "d", c.dependencies, "Name of the dependency to show the logs from")
	return c
}

func (c *serviceLogsCmd) runE(cmd *cobra.Command, args []string) error {
	closer, err := showLogs(c.e, args[0], c.dependencies...)
	if err != nil {
		return err
	}
	defer closer()

	<-xsignal.WaitForInterrupt()
	return nil
}

func showLogs(e ServiceExecutor, serviceID string, dependencies ...string) (closer func(), err error) {
	logs, closer, err := e.ServiceLogs(serviceID, dependencies...)
	if err != nil {
		return nil, err
	}

	// if there was no dependiecies copy all returned
	// by service logs.
	if len(dependencies) == 0 {
		for _, log := range logs {
			dependencies = append(dependencies, log.Dependency)
		}
	}

	prefixes := dependencyPrefixes(dependencies)

	for _, log := range logs {
		go prefixedCopy(os.Stdout, log.Standard, prefixes[log.Dependency])
		go prefixedCopy(os.Stderr, log.Error, prefixes[log.Dependency])
	}

	return closer, nil
}

// dependencyPrefixes returns colored dependency name prefixes.
func dependencyPrefixes(dependencies []string) map[string]string {
	var (
		colors   = pretty.FgColors()
		prefixes = make(map[string]string, len(dependencies))
	)

	max := xstrings.FindLongest(dependencies)
	for i, dep := range dependencies {
		c := colors[i%len(colors)]
		prefixes[dep] = c.Sprintf("% *s | ", max, dep)
	}

	return prefixes
}

// prefixedCopy copies src to dst by prefixing dependency key to each new line.
func prefixedCopy(dst io.Writer, src io.Reader, dep string) {
	io.Copy(dst, prefixedReader(src, dep))
}

// prefixedReader wraps io.Reader by adding a prefix for each new line
// in the stream.
func prefixedReader(r io.Reader, prefix string) io.Reader {
	return prefixer.New(r, fmt.Sprintf("%s ", prefix))
}
