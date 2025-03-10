package podman

import (
	"encoding/json"
	"fmt"
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

func (p *provider) PullImage() error {
	if err := checkCGroupsAndRootFulMode(p.info); err != nil {
		return err
	}
	cmd := exec.Command("podman",
		"pull",
		"--retry=5",
		constants.ImageName,
	)
	out, err := exec.Output(cmd)
	if err != nil {
		return err
	}
	log.Debug(string(out))
	return nil
}

func (p *provider) Create() error {
	if err := checkCGroupsAndRootFulMode(p.info); err != nil {
		return err
	}
	cmd := exec.Command("podman",
		"run",
		"--hostname", constants.HostName,
		"--detach",
		"--rm", "-it", "--privileged",
		"-v", "/var/lib/containers/storage:/host-container:ro,rshared",
		"-p", "9080:80",
		"-p", "9443:443",
		"-p", "6443:6443",
		"--name", constants.ContainerName, constants.ImageName,
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
		cmd := exec.Command("podman",
			"exec",
			constants.ContainerName,
			"systemctl",
			"is-active",
			"microshift",
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
	cmd := exec.Command("podman",
		"exec",
		constants.ContainerName,
		"cat",
		fmt.Sprintf("/var/lib/microshift/resources/kubeadmin/%s/kubeconfig", constants.HostName),
	)
	return exec.Output(cmd)
}

func (p *provider) Delete() error {
	if err := checkCGroupsAndRootFulMode(p.info); err != nil {
		return err
	}
	cmd := exec.Command("podman",
		"rm", "-f",
		constants.ContainerName)
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
	cmd := exec.Command("podman",
		"ps",
		"-s",
		"-f", fmt.Sprintf("name=%s", constants.ContainerName),
		"--format", "{{.Names}} {{.Ports}}",
	)
	out, err := exec.Output(cmd)
	if err != nil {
		return err
	}
	fmt.Printf("%s", out)
	return nil
}

func getProviderInfo() (*providers.ProviderInfo, error) {
	cmd := exec.Command("podman", "info", "--format", "json")
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
