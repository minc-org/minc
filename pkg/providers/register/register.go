package register

import (
	"github.com/minc-org/minc/pkg/providers"
	"github.com/minc-org/minc/pkg/providers/moby"
	"github.com/minc-org/minc/pkg/providers/podman"
	"github.com/minc-org/minc/pkg/rootlessmarker"
	"github.com/spf13/viper"
)

func Register(provider string) (providers.Provider, error) {
	allowRootless := viper.GetBool("allow-rootless") || rootlessmarker.Present()
	switch provider {
	case "podman":
		return podman.New(allowRootless)
	case "docker":
		return moby.New()
	default:
		return podman.New(allowRootless)
	}
}
