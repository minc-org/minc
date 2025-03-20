package minc

import (
	"fmt"
	"github.com/minc-org/minc/pkg/cluster"
	"github.com/minc-org/minc/pkg/constants"
	"github.com/minc-org/minc/pkg/kubeconfig"
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/minc/types"
	"github.com/minc-org/minc/pkg/providers/register"
	"github.com/minc-org/minc/pkg/spinner"
	"time"
)

func Create(cType *types.CreateType) error {
	p, err := register.Register(cType.Provider)
	if err != nil {
		return err
	}
	log.Debug("Provider Info", "Provider", p)
	img := constants.GetUShiftImage(cType.UShiftVersion)
	log.Info(fmt.Sprintf("Ensuring cluster image (%s) ...", img))
	s := spinner.New(time.Second)
	s.Start()
	if err := p.PullImage(img); err != nil {
		return err
	}
	s.Stop()

	s.Start()
	if err := p.Create(cType); err != nil {
		return err
	}
	s.Stop()

	log.Info("Waiting for MicroShift service to start...")
	s.Start()
	if err := p.WaitForMicroShiftService(); err != nil {
		return err
	}
	s.Stop()

	log.Info("Waiting for KubeConfig ...")
	s.Start()
	config, err := p.GetKubeConfig()
	if err != nil {
		return err
	}
	s.Stop()
	if err := kubeconfig.UpdateKubeConfig(config); err != nil {
		return err
	}
	log.Info("Waiting for pods to be ready...")
	if err := cluster.GetPodStatus(config); err != nil {
		return err
	}
	return nil
}
