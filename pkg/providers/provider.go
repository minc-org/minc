package providers

import (
	"github.com/minc-org/minc/pkg/minc/types"
)

type Provider interface {
	Name() string
	Info() (*ProviderInfo, error)
	PullImage(image string) error
	Create(cType *types.CreateType) error
	WaitForMicroShiftService() error
	GetKubeConfig() ([]byte, error)
	Delete() error
	List() error
}

type ProviderInfo struct {
	Rootless bool
	CGroupV2 bool
}
