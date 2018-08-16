package service

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/docker/docker/pkg/archive"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/spinner"
	"github.com/spf13/cobra"
)

// Deploy a service to the marketplace.
var Deploy = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"publish"},
	Short:   "Deploy a service",
	Long: `Deploy a service.

To get more information, see the [deploy page from the documentation](https://docs.mesg.com/guide/service/deploy-a-service.html)`,
	Example:           `mesg-core service deploy PATH_TO_SERVICE`,
	Run:               deployHandler,
	DisableAutoGenTag: true,
}

func deployHandler(cmd *cobra.Command, args []string) {
	path := defaultPath(args)
	serviceID, err := deployService(path)
	utils.HandleError(err)

	fmt.Println("Service deployed with ID:", aurora.Green(serviceID))
	fmt.Printf("To start it, run the command: mesg-core service start %s\n", serviceID)

}

func deployService(path string) (serviceID string, err error) {
	stream, err := cli().DeployService(context.Background())
	if err != nil {
		return "", err
	}

	deployment := make(chan deploymentResult)
	go readDeployReply(stream, deployment)

	if govalidator.IsURL(path) {
		if err := stream.Send(&core.DeployServiceRequest{
			Value: &core.DeployServiceRequest_Url{Url: path},
		}); err != nil {
			return "", err
		}
	} else {
		if err := deployServiceSendServiceContext(path, stream); err != nil {
			return "", err
		}
	}

	if err := stream.CloseSend(); err != nil {
		return "", err
	}

	result := <-deployment
	return result.serviceID, result.err
}

func deployServiceSendServiceContext(path string, stream core.Core_DeployServiceClient) error {
	archive, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	for {
		n, err := archive.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&core.DeployServiceRequest{
			Value: &core.DeployServiceRequest_Chunk{Chunk: buf[:n]},
		}); err != nil {
			return err
		}
	}

	return nil
}

type deploymentResult struct {
	serviceID string
	err       error
}

func readDeployReply(stream core.Core_DeployServiceClient, deployment chan deploymentResult) {
	sp := spinner.New(utils.SpinnerCharset, utils.SpinnerDuration)

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			deployment <- deploymentResult{err: err}
			return
		}

		var (
			status          = message.GetStatus()
			serviceID       = message.GetServiceID()
			validationError = message.GetValidationError()
		)

		switch {
		case status != "":
			sp.Start()
			sp.Suffix = " " + status
		case serviceID != "":
			sp.Stop()
			deployment <- deploymentResult{serviceID: serviceID}
			return
		case validationError != "":
			sp.Stop()
			fmt.Println(aurora.Red(validationError))
			fmt.Println("Run the command 'service validate' for more details")
			os.Exit(0)
		}

	}
}
