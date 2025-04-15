package moby

import (
	"encoding/json"
	"fmt"
	"github.com/minc-org/minc/pkg/minc/types"
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
	return "docker"
}

func (p *provider) Info() (*providers.ProviderInfo, error) {
	return getProviderInfo()
}

func (p *provider) PullImage(image string) error {
	if err := checkCGroupsAndRootFulMode(p.info); err != nil {
		return err
	}
	cmd := exec.Command("docker",
		providers.PullOptions(image)...,
	)
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
	cmd := exec.Command("docker",
		providers.RunOptions(constants.ContainerName, constants.GetUShiftImage(cType.UShiftVersion), cType.UShiftConfig)...,
	)
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
		cmd := exec.Command("docker",
			providers.ServiceWaitOption("microshift", constants.ContainerName)...,
		)
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
	cmd := exec.Command("docker",
		providers.KubeConfigOption(constants.ContainerName, constants.HostName)...,
	)
	return exec.Output(cmd)
}

func (p *provider) Delete() error {
	if err := checkCGroupsAndRootFulMode(p.info); err != nil {
		return err
	}
	cmd := exec.Command("docker",
		providers.DeleteOptions(constants.ContainerName)...,
	)
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
	cmd := exec.Command("docker",
		providers.ListOptions(constants.ContainerName)...,
	)
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
	cmd := exec.Command("docker", "info", "--format", "json")
	out, err := exec.Output(cmd)
	if err != nil {
		return nil, err
	}

	type Response struct {
		CgroupsVersion  string   `json:"CgroupVersion"`
		SecurityOptions []string `json:"SecurityOptions"`
	}

	var res Response
	var cGroupV2, rootless bool
	if err := json.Unmarshal(out, &res); err != nil {
		return nil, err
	}

	if res.CgroupsVersion == "2" {
		cGroupV2 = true
	}

	for _, group := range res.SecurityOptions {
		// sudo docker info -f "{{println .SecurityOptions}}"
		// [name=seccomp,profile=builtin name=selinux name=cgroupns]
		if group == "name=rootless" {
			rootless = true
		}
	}

	return &providers.ProviderInfo{
		Rootless: rootless,
		CGroupV2: cGroupV2,
	}, nil

}

func checkCGroupsAndRootFulMode(pInfo *providers.ProviderInfo) error {
	if !pInfo.CGroupV2 {
		return fmt.Errorf("docker provider requires cgroup v2")
	}
	if pInfo.Rootless {
		return fmt.Errorf("docker provider requires rootful mode")
	}
	return nil
}

// String implements fmt.Stringer
// NOTE: the value of this should not currently be relied upon for anything!
// This is only used for setting the Node's providerID
func (p *provider) String() string {
	return "docker"
}
