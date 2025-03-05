package minc

import (
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/providers/register"
)

func Create(provider string) error {
	p, err := register.Register(provider)
	if err != nil {
		return err
	}
	log.Info("Provider Info", "Provider", p)
	if err := p.Create(); err != nil {
		return err
	}
	return nil
}
