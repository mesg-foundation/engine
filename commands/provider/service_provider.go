package provider

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/mesg-foundation/core/commands/provider/assets"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/utils/chunker"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/mesg-foundation/core/utils/servicetemplate"
	"github.com/mesg-foundation/core/x/xerrors"
)

// ServiceProvider is a struct that provides all methods required by service command.
type ServiceProvider struct {
	client core.CoreClient
}

// NewServiceProvider creates new ServiceProvider.
func NewServiceProvider(c core.CoreClient) *ServiceProvider {
	return &ServiceProvider{client: c}
}

// ServiceByID finds service based on given id.
func (p *ServiceProvider) ServiceByID(id string) (*core.Service, error) {
	serviceReply, err := p.client.GetService(context.Background(), &core.GetServiceRequest{ServiceID: id})
	if err != nil {
		return nil, err
	}

	return serviceReply.Service, nil
}

// ServiceDeleteAll deletes all services.
func (p *ServiceProvider) ServiceDeleteAll() error {
	rep, err := p.client.ListServices(context.Background(), &core.ListServicesRequest{})
	if err != nil {
		return err
	}

	var errs xerrors.Errors
	for _, s := range rep.Services {
		_, err := p.client.DeleteService(context.Background(), &core.DeleteServiceRequest{ServiceID: s.ID})
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs.ErrorOrNil()
}

// ServiceDelete deletes service with given ids.
func (p *ServiceProvider) ServiceDelete(ids ...string) error {
	var errs xerrors.Errors
	for _, id := range ids {
		_, err := p.client.DeleteService(context.Background(), &core.DeleteServiceRequest{ServiceID: id})
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs.ErrorOrNil()
}

// ServiceListenEvents returns a channel with event data streaming..
func (p *ServiceProvider) ServiceListenEvents(id, eventFilter string) (chan *core.EventData, chan error, error) {
	stream, err := p.client.ListenEvent(context.Background(), &core.ListenEventRequest{
		ServiceID:   id,
		EventFilter: eventFilter,
	})
	if err != nil {
		return nil, nil, err
	}

	resultC := make(chan *core.EventData)
	errC := make(chan error)

	go func() {
		for {
			if res, err := stream.Recv(); err != nil {
				errC <- err
				break
			} else {
				resultC <- res
			}
		}
	}()
	return resultC, errC, nil

}

// ServiceListenResults returns a channel with event results streaming..
func (p *ServiceProvider) ServiceListenResults(id, taskFilter, outputFilter string, tagFilters []string) (chan *core.ResultData, chan error, error) {
	stream, err := p.client.ListenResult(context.Background(), &core.ListenResultRequest{
		ServiceID:    id,
		TaskFilter:   taskFilter,
		OutputFilter: outputFilter,
		TagFilters:   tagFilters,
	})
	if err != nil {
		return nil, nil, err
	}
	resultC := make(chan *core.ResultData)
	errC := make(chan error)

	go func() {
		for {
			if res, err := stream.Recv(); err != nil {
				errC <- err
				break
			} else {
				resultC <- res
			}
		}
	}()
	return resultC, errC, nil
}

// Log keeps dependency logs of service.
type Log struct {
	Dependency      string
	Standard, Error *chunker.Stream
}

// ServiceLogs returns logs reader for all service dependencies.
func (p *ServiceProvider) ServiceLogs(id string, dependencies ...string) (logs []*Log, close func(), err error) {
	if len(dependencies) == 0 {
		resp, err := p.client.GetService(context.Background(), &core.GetServiceRequest{
			ServiceID: id,
		})
		if err != nil {
			return nil, nil, err
		}
		for _, dep := range resp.Service.Dependencies {
			dependencies = append(dependencies, dep.Key)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	stream, err := p.client.ServiceLogs(ctx, &core.ServiceLogsRequest{
		ServiceID:    id,
		Dependencies: dependencies,
	})
	if err != nil {
		cancel()
		return nil, nil, err
	}

	for _, key := range dependencies {
		log := &Log{
			Dependency: key,
			Standard:   chunker.NewStream(),
			Error:      chunker.NewStream(),
		}
		logs = append(logs, log)
	}

	closer := func() {
		cancel()
		for _, log := range logs {
			log.Standard.Close()
			log.Error.Close()
		}
	}

	errC := make(chan error, len(logs))
	go p.listenServiceLogs(stream, logs, errC)
	go func() {
		<-errC
		closer()
	}()

	return logs, closer, nil
}

// listenServiceLogs listen gRPC stream to get service logs.
func (p *ServiceProvider) listenServiceLogs(stream core.Core_ServiceLogsClient, logs []*Log,
	errC chan error) {
	for {
		data, err := stream.Recv()
		if err != nil {
			errC <- err
			return
		}

		for _, log := range logs {
			if log.Dependency == data.Dependency {
				var out *chunker.Stream
				switch data.Type {
				case core.LogData_Standard:
					out = log.Standard
				case core.LogData_Error:
					out = log.Error
				}
				out.Provide(data.Data)
			}
		}
	}
}

// ServiceExecuteTask executes task on given service.
func (p *ServiceProvider) ServiceExecuteTask(id, taskKey, inputData string, tags []string) error {
	_, err := p.client.ExecuteTask(context.Background(), &core.ExecuteTaskRequest{
		ServiceID:     id,
		TaskKey:       taskKey,
		InputData:     inputData,
		ExecutionTags: tags,
	})
	return err
}

// ServiceStart starts a service.
func (p *ServiceProvider) ServiceStart(id string) error {
	_, err := p.client.StartService(context.Background(), &core.StartServiceRequest{ServiceID: id})
	return err
}

// ServiceStop stops a service.
func (p *ServiceProvider) ServiceStop(id string) error {
	_, err := p.client.StopService(context.Background(), &core.StopServiceRequest{ServiceID: id})
	return err
}

// ServiceValidate validates a service configuration and Dockerfile.
func (p *ServiceProvider) ServiceValidate(path string) (string, error) {
	validation, err := importer.Validate(path)
	if err != nil {
		return "", err
	}

	if !validation.ServiceFileExist {
		return fmt.Sprintf("%s File 'mesg.yml' does not exist", pretty.FailSign), nil
	}

	if len(validation.ServiceFileWarnings) > 0 {
		var msg = fmt.Sprintf("%s File 'mesg.yml' is not valid. See documentation: https://docs.mesg.com/guide/service/service-file.html\n", pretty.FailSign)
		for _, warning := range validation.ServiceFileWarnings {
			msg += fmt.Sprintf("\t* %s\n", warning)
		}
		return msg, nil
	}

	if !validation.DockerfileExist {
		return fmt.Sprintf("%s Dockerfile does not exist", pretty.FailSign), nil
	}
	if !validation.IsValid() {
		return fmt.Sprintf("%s Service is not valid", pretty.FailSign), nil
	}

	return fmt.Sprintf(`%s Dockerfile exists\n
%s mesg.yml is valid
%s Service is valid`, pretty.SuccessSign, pretty.SuccessSign, pretty.SuccessSign), nil
}

// ServiceGenerateDocs creates docs in given path.
func (p *ServiceProvider) ServiceGenerateDocs(path string) error {
	readmePath := filepath.Join(path, "README.md")
	service, err := importer.From(path)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(readmePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer f.Close()

	readmeTemplate := assets.MustAsset("readme_template.md")

	tmpl := template.Must(template.New("doc").Parse(string(readmeTemplate)))
	return tmpl.Execute(f, service)
}

// ServiceList lists all services.
func (p *ServiceProvider) ServiceList() ([]*core.Service, error) {
	reply, err := p.client.ListServices(context.Background(), &core.ListServicesRequest{})
	if err != nil {
		return nil, err
	}
	return reply.Services, nil
}

// ServiceInitTemplateList downloads services templates list from awesome github repo.
func (p *ServiceProvider) ServiceInitTemplateList() ([]*servicetemplate.Template, error) {
	return servicetemplate.List()
}

// ServiceInitDownloadTemplate download given service template.
func (p *ServiceProvider) ServiceInitDownloadTemplate(t *servicetemplate.Template, dst string) error {
	return servicetemplate.Download(t, dst)
}
