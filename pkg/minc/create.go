package minc

import (
	"github.com/minc-org/minc/pkg/cluster"
	"github.com/minc-org/minc/pkg/kubeconfig"
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

	log.Info("Waiting for API server to start...")
	if err := p.WaitForAPI(); err != nil {
		return err
	}

	log.Info("Waiting for KubeConfig ...")
	config, err := p.GetKubeConfig()
	if err != nil {
		return err
	}
	if err := kubeconfig.UpdateKubeConfig(config); err != nil {
		return err
	}
	log.Info("Waiting for pods to be ready...")
	if err := cluster.GetPodStatus(config); err != nil {
		return err
	}
	return nil
}
