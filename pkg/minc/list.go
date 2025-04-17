package minc

import (
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/providers/register"
)

func List(provider string) ([]byte, error) {
	p, err := register.Register(provider)
	if err != nil {
		return nil, err
	}
	log.Debug("Provider Info", "Provider", p)
	return p.List()
}
