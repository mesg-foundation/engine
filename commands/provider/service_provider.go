package provider

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/mesg-foundation/core/commands/provider/assets"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/mesg-foundation/core/utils/servicetemplate"
	"github.com/mesg-foundation/core/x/xerrors"
)

type ServiceProvider struct {
	client core.CoreClient
}

func NewServiceProvider(c core.CoreClient) *ServiceProvider {
	return &ServiceProvider{client: c}
}

func (p *ServiceProvider) ServiceByID(id string) (*core.Service, error) {
	serviceReply, err := p.client.GetService(context.Background(), &core.GetServiceRequest{ServiceID: id})
	if err != nil {
		return nil, err
	}

	return serviceReply.Service, nil
}

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

func (p *ServiceProvider) ServiceListenEvents(id, eventFilter string) (chan *core.EventData, chan error, error) {
	stream, err := p.client.ListenEvent(context.Background(), &core.ListenEventRequest{
		ServiceID:   id,
		EventFilter: eventFilter,
	})
	if err != nil {
		return nil, nil, err
	}

	reslutC := make(chan *core.EventData)
	errC := make(chan error)

	go func() {
		for {
			if res, err := stream.Recv(); err != nil {
				errC <- err
			} else {
				reslutC <- res
			}
		}
	}()
	return reslutC, errC, nil

}

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
	reslutC := make(chan *core.ResultData)
	errC := make(chan error)

	go func() {
		for {
			if res, err := stream.Recv(); err != nil {
				errC <- err
			} else {
				reslutC <- res
			}
		}
	}()
	return reslutC, errC, nil
}

func (p *ServiceProvider) ServiceLogs(id string) (io.ReadCloser, error) {
	rs, err := p.ServiceDependencyLogs(id, "*")
	if err != nil {
		return nil, err
	}

	if len(rs) != 1 {
		return nil, errors.New("no valid readers")
	}
	return rs[0], nil
}

func (p *ServiceProvider) ServiceDependencyLogs(id string, dependency string) ([]io.ReadCloser, error) {
	// TODO: wait for feature fix-cmd-logs to be merged
	return nil, errors.New("logs unimplementd")
}

func (p *ServiceProvider) ServiceExecuteTask(id, taskKey, inputData string, tags []string) error {
	_, err := p.client.ExecuteTask(context.Background(), &core.ExecuteTaskRequest{
		ServiceID:     id,
		TaskKey:       taskKey,
		InputData:     inputData,
		ExecutionTags: tags,
	})
	return err
}

func (p *ServiceProvider) ServiceStart(id string) error {
	_, err := p.client.StartService(context.Background(), &core.StartServiceRequest{ServiceID: id})
	return err
}

func (p *ServiceProvider) ServiceStop(id string) error {
	_, err := p.client.StopService(context.Background(), &core.StopServiceRequest{ServiceID: id})
	return err
}

func (p *ServiceProvider) ServiceValidate(path string) (string, error) {
	validation, err := importer.Validate(path)
	if err != nil {
		return "", err
	}

	if !validation.ServiceFileExist {
		return fmt.Sprintf("%s File 'mesg.yml' does not exist\n", pretty.FailSign), nil
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

	return fmt.Sprintf(`
		%s Dockerfile exists
		%s mesg.yml is valid
		%s Service is valid
	`, pretty.SuccessSign, pretty.SuccessSign, pretty.SuccessSign), nil
}

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

func (p *ServiceProvider) ServiceList() ([]*core.Service, error) {
	reply, err := p.client.ListServices(context.Background(), &core.ListServicesRequest{})
	if err != nil {
		return nil, err
	}
	return reply.Services, nil
}

func (p *ServiceProvider) ServiceInitTemplateList() ([]*servicetemplate.Template, error) {
	return servicetemplate.List()
}

func (p *ServiceProvider) ServiceInitDownloadTemplate(t *servicetemplate.Template, dst string) error {
	return servicetemplate.Download(t, dst)
}
