package minc

import (
	"fmt"
	"github.com/minc-org/minc/pkg/cluster"
	"github.com/minc-org/minc/pkg/constants"
	"github.com/minc-org/minc/pkg/kubeconfig"
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/providers/register"
)

func Create(provider string) error {
	p, err := register.Register(provider)
	if err != nil {
		return err
	}

	log.Debug("Provider Info", "Provider", p)
	log.Info(fmt.Sprintf("Ensuring cluster image (%s) ...", constants.ImageName))
	if err := p.PullImage(); err != nil {
		return err
	}

	if err := p.Create(); err != nil {
		return err
	}

	log.Info("Waiting for MicroShift service to start...")
	if err := p.WaitForMicroShiftService(); err != nil {
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
