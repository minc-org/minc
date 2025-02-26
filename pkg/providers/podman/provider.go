package podman

import (
	"github.com/minc-org/minc/pkg/providers"
	"log"
)

// Provider implements provider.Provider
// see NewProvider
type provider struct {
	logger log.Logger
	info   *providers.ProviderInfo
}

func New(info *providers.ProviderInfo) providers.Provider {
	return &provider{
		logger: 
	}
}
