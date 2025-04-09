package podman

import (
	"encoding/json"
	"fmt"
	"github.com/minc-org/minc/pkg/minc/types"
	"runtime"
	"time"

	"github.com/minc-org/minc/pkg/constants"
	"github.com/minc-org/minc/pkg/exec"
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/providers"
	"github.com/minc-org/minc/pkg/retry"
)

// Provider implements provider.Provider
// see NewProvider
type provider struct {
	info *providers.ProviderInfo
}

func New() (providers.Provider, error) {
	pInfo, err := getProviderInfo()
	if err != nil {
		return nil, err
	}
	return &provider{pInfo}, nil
}

func (p *provider) Name() string {
	return "podman"
}

func (p *provider) Info() (*providers.ProviderInfo, error) {
	return getProviderInfo()
}

func (p *provider) PullImage(image string) error {
	if err := checkCGroupsAndRootFulMode(p.info); err != nil {
		return err
	}
	cmd := podmanCmd(providers.PullOptions(image))
	out, err := exec.Output(cmd)
	if err != nil {
		return err
	}
	log.Debug(string(out))
	return nil
}

func (p *provider) Create(cType *types.CreateType) error {
	if err := checkCGroupsAndRootFulMode(p.info); err != nil {
		return err
	}
	cmd := podmanCmd(providers.RunOptions(constants.ContainerName, constants.GetUShiftImage(cType.UShiftVersion)))
	out, err := exec.Output(cmd)
	if err != nil {
		return err
	}
	log.Debug(string(out))
	return nil
}

func (p *provider) WaitForMicroShiftService() error {
	if err := checkCGroupsAndRootFulMode(p.info); err != nil {
		return err
	}
	cmdFunc := func() error {
		cmd := podmanCmd(providers.ServiceWaitOption("microshift", constants.ContainerName))
		out, err := exec.Output(cmd)
		if err != nil {
			return err
		}

		log.Debug(string(out))
		return nil
	}
	return retry.Retry(cmdFunc, 15, 2*time.Second)
}

func (p *provider) GetKubeConfig() ([]byte, error) {
	if err := checkCGroupsAndRootFulMode(p.info); err != nil {
		return nil, err
	}
	cmd := podmanCmd(providers.KubeConfigOption(constants.ContainerName, constants.HostName))
	return exec.Output(cmd)
}

func (p *provider) Delete() error {
	if err := checkCGroupsAndRootFulMode(p.info); err != nil {
		return err
	}
	cmd := podmanCmd(providers.DeleteOptions(constants.ContainerName))
	out, err := exec.Output(cmd)
	if err != nil {
		return err
	}
	log.Debug(string(out))
	return nil
}

func (p *provider) List() error {
	if err := checkCGroupsAndRootFulMode(p.info); err != nil {
		return err
	}
	cmd := podmanCmd(providers.ListOptions(constants.ContainerName))
	out, err := exec.Output(cmd)
	if err != nil {
		return err
	}
	if string(out) == "" {
		return fmt.Errorf("no %s containers found, use 'create' command to create it", constants.ContainerName)
	}
	fmt.Printf("%s", out)
	return nil
}

func getProviderInfo() (*providers.ProviderInfo, error) {
	cmd := podmanCmd([]string{"info", "--format", "json"})
	out, err := exec.Output(cmd)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Host struct {
			CgroupsVersion string `json:"cgroupVersion"`
			Security       struct {
				Rootless bool `json:"rootless"`
			} `json:"security"`
		} `json:"host"`
	}

	var res Response
	var cGroupV2 bool
	if err := json.Unmarshal(out, &res); err != nil {
		return nil, err
	}

	if res.Host.CgroupsVersion == "v2" {
		cGroupV2 = true
	}

	return &providers.ProviderInfo{
		Rootless: res.Host.Security.Rootless,
		CGroupV2: cGroupV2,
	}, nil

}

func checkCGroupsAndRootFulMode(pInfo *providers.ProviderInfo) error {
	if !pInfo.CGroupV2 {
		return fmt.Errorf("podman provider requires cgroup v2")
	}
	if pInfo.Rootless {
		return fmt.Errorf("podman provider requires rootful mode")
	}
	return nil
}

// String implements fmt.Stringer
// NOTE: the value of this should not currently be relied upon for anything!
// This is only used for setting the Node's providerID
func (p *provider) String() string {
	return "podman"
}

func podmanCmd(args []string) exec.Cmd {
	if runtime.GOOS == "linux" {
		return exec.Command("sudo", append([]string{"podman"}, args...)...)
	} else {
		return exec.Command("podman", args...)
	}
}
