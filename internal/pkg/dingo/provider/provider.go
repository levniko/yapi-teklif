package provider

import (
	"github.com/sarulabs/dingo/v4"

	services "github.com/yapi-teklif/internal/pkg/dingo/servicedefs"
)

type Provider struct {
	dingo.BaseProvider
}

func (p *Provider) Load() error {
	if err := p.AddDefSlice(services.HandlersDefs); err != nil {
		return err
	}
	if err := p.AddDefSlice(services.ExternalConnectionDefs); err != nil {
		return err
	}
	if err := p.AddDefSlice(services.ManagersDefs); err != nil {
		return err
	}
	if err := p.AddDefSlice(services.ServicesDefs); err != nil {
		return err
	}
	if err := p.AddDefSlice(services.RepositoriesDefs); err != nil {
		return err
	}

	return nil
}
