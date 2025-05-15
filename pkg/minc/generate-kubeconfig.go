package minc

import (
	"github.com/minc-org/minc/pkg/kubeconfig"
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/providers/register"
)

func GenerateKubeConfig(provider string) error {
	p, err := register.Register(provider)
	if err != nil {
		return err
	}
	log.Debug("Provider Info", "Provider", p)
	if _, err := p.List(); err != nil {
		return err
	}
	config, err := p.GetKubeConfig()
	if err != nil {
		return err
	}
	if err := kubeconfig.UpdateKubeConfig(config); err != nil {
		return err
	}
	return nil
}
