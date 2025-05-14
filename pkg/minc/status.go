package minc

import (
	"github.com/minc-org/minc/pkg/cluster"
	"github.com/minc-org/minc/pkg/minc/types"
	"github.com/minc-org/minc/pkg/providers/register"
	"strings"
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
	out, err := p.List()
	if err != nil {
		status.Error = err.Error()
		return &status
	}
	if strings.Contains(string(out), "running") {
		status.Container = "running"
	}
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
