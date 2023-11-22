package extconsul

import (
	"fmt"

	"github.com/AlgerDu/go-dream/src/dinfra"
	capi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

type (
	ConsulRegisterOptions struct {
		Address string
		Token   string

		Env string
	}

	ConsulRegister struct {
		logger  dinfra.Logger
		options *ConsulRegisterOptions

		client *capi.Client
	}
)

func New(
	logger dinfra.Logger,
	options *ConsulRegisterOptions,
) (*ConsulRegister, error) {

	client, err := capi.NewClient(&capi.Config{
		Address: options.Address,
		Token:   options.Token,
	})

	register := &ConsulRegister{
		logger:  logger,
		options: options,
		client:  client,
	}

	return register, err
}

func (register *ConsulRegister) Get(name string) (*dinfra.MicroService, error) {

	filter := fmt.Sprintf("Service == %s and Meta.Env == %s", name, register.options.Env)

	logger := register.logger.WithFields(logrus.Fields{
		"filter": filter,
	})

	services, err := register.client.Agent().ServicesWithFilter(filter)
	if err != nil {
		logger.WithError(err).Error("get service from consul err")
		return nil, err
	}

	if len(services) == 0 {
		logger.WithError(ErrNoFitService).Error("no fit service")
		return nil, ErrNoFitService
	}

	var cService *capi.AgentService
	for _, v := range services {
		cService = v
		break
	}

	logger.WithFields(logrus.Fields{
		"serviceID": cService.ID,
	}).Trace("get service success")

	return consulServiceToService(cService), nil
}

func (*ConsulRegister) Register(service *dinfra.MicroService) error {

	cService := &capi.AgentService{}
	consulServiceSetByService(cService, service)

	panic("unimplemented")
}
