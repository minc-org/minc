package minc

import (
	"github.com/minc-org/minc/pkg/kubeconfig"
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/providers/register"
)

func Delete(provider string) error {
	p, err := register.Register(provider)
	if err != nil {
		return err
	}
	log.Debug("Provider Info", "Provider", p)
	if err := p.Delete(); err != nil {
		return err
	}
	log.Info("Removing entry from kubeconfig ...")
	if err := kubeconfig.RemoveClusterFromConfig(); err != nil {
		return err
	}
	return nil
}
