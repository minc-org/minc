package register

import (
	"github.com/minc-org/minc/pkg/providers"
	"github.com/minc-org/minc/pkg/providers/moby"
	"github.com/minc-org/minc/pkg/providers/podman"
)

func Register(provider string) (providers.Provider, error) {
	switch provider {
	case "podman":
		return podman.New()
	case "docker":
		return moby.New()
	default:
		return podman.New()
	}
}
