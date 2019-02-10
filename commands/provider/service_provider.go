package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/mesg-foundation/core/commands/provider/assets"
	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/chunker"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/mesg-foundation/core/utils/servicetemplate"
	"github.com/mesg-foundation/core/x/xerrors"
)

// ServiceProvider is a struct that provides all methods required by service command.
type ServiceProvider struct {
	client coreapi.CoreClient
}

// NewServiceProvider creates new ServiceProvider.
func NewServiceProvider(c coreapi.CoreClient) *ServiceProvider {
	return &ServiceProvider{client: c}
}

// ServiceByID finds service based on given id.
func (p *ServiceProvider) ServiceByID(id string) (*coreapi.Service, error) {
	serviceReply, err := p.client.GetService(context.Background(), &coreapi.GetServiceRequest{ServiceID: id})
	if err != nil {
		return nil, err
	}

	return serviceReply.Service, nil
}

// ServiceDeleteAll deletes all services.
func (p *ServiceProvider) ServiceDeleteAll(deleteData bool) error {
	rep, err := p.client.ListServices(context.Background(), &coreapi.ListServicesRequest{})
	if err != nil {
		return err
	}

	var (
		errs xerrors.SyncErrors
		wg   sync.WaitGroup
	)
	wg.Add(len(rep.Services))
	for _, s := range rep.Services {
		go func(id string) {
			_, err := p.client.DeleteService(context.Background(), &coreapi.DeleteServiceRequest{
				ServiceID:  id,
				DeleteData: deleteData,
			})
			if err != nil {
				errs.Append(err)
			}
			wg.Done()
		}(s.Sid)
	}
	wg.Wait()
	return errs.ErrorOrNil()
}

// ServiceDelete deletes service with given ids.
func (p *ServiceProvider) ServiceDelete(deleteData bool, ids ...string) error {
	var errs xerrors.Errors
	for _, id := range ids {
		if _, err := p.client.DeleteService(context.Background(), &coreapi.DeleteServiceRequest{
			ServiceID:  id,
			DeleteData: deleteData,
		}); err != nil {
			errs = append(errs, err)
		}
	}
	return errs.ErrorOrNil()
}

// ServiceListenEvents returns a channel with event data streaming..
func (p *ServiceProvider) ServiceListenEvents(id, eventFilter string) (chan *coreapi.EventData, chan error, error) {
	stream, err := p.client.ListenEvent(context.Background(), &coreapi.ListenEventRequest{
		ServiceID:   id,
		EventFilter: eventFilter,
	})
	if err != nil {
		return nil, nil, err
	}

	resultC := make(chan *coreapi.EventData)
	errC := make(chan error)

	go func() {
		<-stream.Context().Done()
		errC <- stream.Context().Err()
	}()
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

	if err := acknowledgement.WaitForStreamToBeReady(stream); err != nil {
		return nil, nil, err
	}

	return resultC, errC, nil
}

// ServiceListenResults returns a channel with event results streaming..
func (p *ServiceProvider) ServiceListenResults(id, taskFilter, outputFilter string, tagFilters []string) (chan *coreapi.ResultData, chan error, error) {
	stream, err := p.client.ListenResult(context.Background(), &coreapi.ListenResultRequest{
		ServiceID:    id,
		TaskFilter:   taskFilter,
		OutputFilter: outputFilter,
		TagFilters:   tagFilters,
	})
	if err != nil {
		return nil, nil, err
	}
	resultC := make(chan *coreapi.ResultData)
	errC := make(chan error)

	go func() {
		<-stream.Context().Done()
		errC <- stream.Context().Err()
	}()
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

	if err := acknowledgement.WaitForStreamToBeReady(stream); err != nil {
		return nil, nil, err
	}

	return resultC, errC, nil
}

// Log keeps dependency logs of service.
type Log struct {
	Dependency      string
	Standard, Error *chunker.Stream
}

// ServiceLogs returns logs reader for all service dependencies.
func (p *ServiceProvider) ServiceLogs(id string, dependencies ...string) (logs []*Log, close func(), errC chan error, err error) {
	if len(dependencies) == 0 {
		resp, err := p.client.GetService(context.Background(), &coreapi.GetServiceRequest{
			ServiceID: id,
		})
		if err != nil {
			return nil, nil, nil, err
		}

		for key := range resp.Service.Dependencies {
			dependencies = append(dependencies, key)
		}
		dependencies = append(dependencies, service.MainServiceKey)
	}

	ctx, cancel := context.WithCancel(context.Background())

	stream, err := p.client.ServiceLogs(ctx, &coreapi.ServiceLogsRequest{
		ServiceID:    id,
		Dependencies: dependencies,
	})
	if err != nil {
		cancel()
		return nil, nil, nil, err
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

	errC = make(chan error)
	go func() {
		<-stream.Context().Done()
		errC <- stream.Context().Err()
	}()
	go p.listenServiceLogs(stream, logs, errC)

	if err := acknowledgement.WaitForStreamToBeReady(stream); err != nil {
		closer()
		return nil, nil, nil, err
	}

	return logs, closer, errC, nil
}

// listenServiceLogs listen gRPC stream to get service logs.
func (p *ServiceProvider) listenServiceLogs(stream coreapi.Core_ServiceLogsClient, logs []*Log,
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
				case coreapi.LogData_Standard:
					out = log.Standard
				case coreapi.LogData_Error:
					out = log.Error
				}
				out.Provide(data.Data)
			}
		}
	}
}

// ServiceExecuteTask executes task on given service.
func (p *ServiceProvider) ServiceExecuteTask(id, taskKey, inputData string, tags []string) error {
	_, err := p.client.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID:     id,
		TaskKey:       taskKey,
		InputData:     inputData,
		ExecutionTags: tags,
	})
	return err
}

// ServiceStart starts a service.
func (p *ServiceProvider) ServiceStart(id string) error {
	_, err := p.client.StartService(context.Background(), &coreapi.StartServiceRequest{ServiceID: id})
	return err
}

// ServiceStop stops a service.
func (p *ServiceProvider) ServiceStop(id string) error {
	_, err := p.client.StopService(context.Background(), &coreapi.StopServiceRequest{ServiceID: id})
	return err
}

// ServiceValidate validates a service configuration and Dockerfile.
func (p *ServiceProvider) ServiceValidate(path string) (string, error) {
	if _, err := service.ReadDefinition(path); err != nil {
		return "", nil
	}
	return fmt.Sprintf(`%s mesg.yml is valid`, pretty.SuccessSign), nil
}

// ServiceGenerateDocs creates docs in given path.
func (p *ServiceProvider) ServiceGenerateDocs(path string) error {
	readmePath := filepath.Join(path, "README.md")
	service, err := service.ReadDefinition(path)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(readmePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	readmeTemplate, err := assets.Asset("commands/provider/assets/readme_template.md")
	if err != nil {
		return err
	}

	anchorEncode := func(a string) string {
		a = strings.Replace(a, " ", "-", -1)
		a = strings.Replace(a, "'", "", -1)
		a = strings.ToLower(a)
		return a
	}
	tpl, err := template.New("doc").Funcs(template.FuncMap{"anchorEncode": anchorEncode}).Parse(string(readmeTemplate))
	if err != nil {
		return err
	}
	return tpl.Execute(f, service)
}

// ServiceList lists all services.
func (p *ServiceProvider) ServiceList() ([]*coreapi.Service, error) {
	reply, err := p.client.ListServices(context.Background(), &coreapi.ListServicesRequest{})
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
