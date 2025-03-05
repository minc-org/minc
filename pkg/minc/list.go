package minc

import (
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/providers/register"
)

func List(provider string) error {
	p, err := register.Register(provider)
	if err != nil {
		return err
	}
	log.Info("Provider Info", "Provider", p)
	if err := p.List(); err != nil {
		return err
	}
	return nil
}
