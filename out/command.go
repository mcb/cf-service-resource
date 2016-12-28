package out

import (
	"time"

	"github.com/mcb/cf-service-resource"
)

type Command struct {
	paas PAAS
}

func NewCommand(paas PAAS) *Command {
	return &Command{
		paas: paas,
	}
}

func (command *Command) Run(request Request) (Response, error) {
	err := command.paas.Login(
		request.Source.API,
		request.Source.Username,
		request.Source.Password,
		request.Source.SkipCertCheck,
	)
	if err != nil {
		return Response{}, err
	}

	err = command.paas.Target(
		request.Source.Organization,
		request.Source.Space,
	)
	if err != nil {
		return Response{}, err
	}

	err = command.paas.CreateService(
		request.Params.Service,
		request.Params.Plan,
		request.Params.InstanceName,
	)
	if err != nil {
		return Response{}, err
	}

  if !request.Params.SkipBinding {
		err = command.paas.BindService(
			request.Params.CurrentAppName,
			request.Params.InstanceName,
		)
		if err != nil {
			return Response{}, err
		}
	}

	if !request.Params.SkipRestage {
		err = command.paas.RestageApp(
			request.Params.CurrentAppName,
		)
		if err != nil {
			return Response{}, err
		}
	}

	return Response{
		Version: resource.Version{
			Timestamp: time.Now(),
		},
		Metadata: []resource.MetadataPair{
			{
				Name:  "organization",
				Value: request.Source.Organization,
			},
			{
				Name:  "space",
				Value: request.Source.Space,
			},
		},
	}, nil
}
