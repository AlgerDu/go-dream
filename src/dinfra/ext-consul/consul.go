package extconsul

import (
	"errors"

	"github.com/AlgerDu/go-dream/src/dinfra"
	capi "github.com/hashicorp/consul/api"
)

var (
	ErrNoFitService = errors.New("no fit service")
)

func consulServiceToService(cService *capi.AgentService) *dinfra.MicroService {

	infraService := &dinfra.MicroService{
		ID:       cService.ID,
		Name:     cService.Service,
		Address:  cService.Address,
		Port:     cService.Port,
		Env:      "",
		Tags:     []string{},
		Metadata: map[string]string{},
	}

	infraService.Tags = append(infraService.Tags, cService.Tags...)

	for k, v := range cService.Meta {

		if k == "Env" {
			infraService.Env = v
			continue
		}

		infraService.Metadata[k] = v
	}

	return infraService
}

func consulServiceSetByService(cService *capi.AgentService, infraService *dinfra.MicroService) {

	cService.ID = infraService.ID
	cService.Service = infraService.Name
	cService.Address = infraService.Address
	cService.Port = infraService.Port
	cService.Tags = append(cService.Tags, infraService.Tags...)

	for k, v := range infraService.Metadata {
		cService.Meta[k] = v
	}

	cService.Meta["Env"] = infraService.Env
}
