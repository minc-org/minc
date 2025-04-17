package minc

import (
	"github.com/minc-org/minc/pkg/cluster"
	"github.com/minc-org/minc/pkg/minc/types"
	"github.com/minc-org/minc/pkg/providers/register"
)

func Status(provider string) *types.StatusType {
	status := types.StatusType{
		Container: "stopped",
		APIServer: "stopped",
	}
	p, err := register.Register(provider)
	if err != nil {
		status.Error = err.Error()
		return &status
	}
	if _, err := p.List(); err != nil {
		status.Error = err.Error()
		return &status
	}
	status.Container = "running"
	config, err := p.GetKubeConfig()
	if err != nil {
		status.Error = err.Error()
		return &status
	}
	if err := cluster.GetPodStatus(config); err != nil {
		status.Error = err.Error()
		return &status
	}
	status.APIServer = "running"
	return &status
}
