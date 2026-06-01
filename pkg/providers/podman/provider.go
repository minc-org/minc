package podman

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/minc-org/minc/pkg/minc/types"

	"github.com/minc-org/minc/pkg/constants"
	"github.com/minc-org/minc/pkg/exec"
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/providers"
	"github.com/minc-org/minc/pkg/retry"
)

// Provider implements provider.Provider
type provider struct {
	info          *providers.ProviderInfo
	useSudo       bool
	allowRootless bool
}

// New builds a Podman provider. When allowRootless is true, minc stays on the user's rootless
// Podman (no sudo) and skips the rootful-only guard; MicroShift may still fail at runtime.
func New(allowRootless bool) (providers.Provider, error) {
	rootlessRuntime, err := rootlessFromUserPodman()
	if err != nil {
		return nil, err
	}
	useSudo := rootlessRuntime && !allowRootless

	p := &provider{
		allowRootless: allowRootless,
		useSudo:       useSudo,
	}
	info, err := p.fetchProviderInfo()
	if err != nil {
		return nil, err
	}
	p.info = info
	return p, nil
}

func rootlessFromUserPodman() (bool, error) {
	cmd := exec.Command("podman", "info", "--format", "{{.Host.Security.Rootless}}")
	out, err := exec.Output(cmd)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(strings.TrimSpace(string(out)))
}

func (p *provider) Name() string {
	return "podman"
}

func (p *provider) Info() (*providers.ProviderInfo, error) {
	return p.fetchProviderInfo()
}

func (p *provider) ImageExists(image string) bool {
	cmd := p.podmanCmd(providers.ImageExistOptions(image))
	_, err := exec.Output(cmd)
	if err != nil {
		return false
	}
	return true
}

func (p *provider) PullImage(image string) error {
	if err := p.checkCGroupsAndRootFulMode(); err != nil {
		return err
	}
	if p.ImageExists(image) {
		return nil
	}
	cmd := p.podmanCmd(providers.PullOptions(image))
	out, err := exec.Output(cmd)
	if err != nil {
		return err
	}
	log.Debug(string(out))
	return nil
}

func (p *provider) storeGraphRoot() (string, error) {
	cmd := p.podmanCmd([]string{"info", "--format", "{{.Store.GraphRoot}}"})
	out, err := exec.Output(cmd)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (p *provider) Create(cType *types.CreateType) error {
	if err := p.checkCGroupsAndRootFulMode(); err != nil {
		return err
	}
	if out, _ := p.List(); len(out) == 0 {
		graphRoot, err := p.storeGraphRoot()
		if err != nil {
			return fmt.Errorf("podman store graph root: %w", err)
		}
		cOptions := &providers.COptions{
			ContainerName:        constants.ContainerName,
			ImageName:            constants.GetUShiftImage(cType.UShiftImage, cType.UShiftVersion),
			UShiftConfig:         cType.UShiftConfig,
			HttpPort:             cType.HTTPPort,
			HttpsPort:            cType.HTTPSPort,
			DisableOverlayCache:  cType.DisableOverlayCache,
			HostContainerStorage: graphRoot,
			AllowRootless:        p.allowRootless,
		}
		if p.allowRootless {
			rlConf, err := writeRootlessConfigs()
			if err != nil {
				return fmt.Errorf("writing rootless configs: %w", err)
			}
			cOptions.RootlessMicroShiftConfig = rlConf.microshiftConf
			cOptions.RootlessCRIOConfig = rlConf.crioConf
			cOptions.RootlessCrunWrapper = rlConf.crunWrapper
		}
		cmd := p.podmanCmd(providers.CreateOptions(cOptions))
		out, err := exec.Output(cmd)
		if err != nil {
			return err
		}
		log.Debug(string(out))
	}
	cmd := p.podmanCmd(providers.StartOptions(constants.ContainerName))
	out, err := exec.Output(cmd)
	if err != nil {
		return err
	}
	log.Debug(string(out))
	return nil
}

func (p *provider) WaitForMicroShiftService() error {
	if err := p.checkCGroupsAndRootFulMode(); err != nil {
		return err
	}
	cmdFunc := func() error {
		cmd := p.podmanCmd(providers.ServiceWaitOption("microshift", constants.ContainerName))
		out, err := exec.Output(cmd)
		if err != nil {
			return err
		}

		log.Debug(string(out))
		return nil
	}
	return retry.Retry(cmdFunc, providers.MicroShiftServiceMaxRetries, providers.MicroShiftServiceInitialRetryDelay)
}

func (p *provider) GetKubeConfig() ([]byte, error) {
	if err := p.checkCGroupsAndRootFulMode(); err != nil {
		return nil, err
	}
	cmd := p.podmanCmd(providers.KubeConfigOption(constants.ContainerName, constants.HostName))
	return exec.Output(cmd)
}

func (p *provider) Delete() error {
	if err := p.checkCGroupsAndRootFulMode(); err != nil {
		return err
	}
	cmd := p.podmanCmd(providers.DeleteOptions(constants.ContainerName))
	out, err := exec.Output(cmd)
	if err != nil {
		return err
	}
	log.Debug(string(out))
	return nil
}

func (p *provider) List() ([]byte, error) {
	var out []byte
	if err := p.checkCGroupsAndRootFulMode(); err != nil {
		return out, err
	}
	cmd := p.podmanCmd(providers.ListOptions(constants.ContainerName))
	out, err := exec.Output(cmd)
	if err != nil {
		return out, err
	}
	log.Debug(string(out))
	if len(out) == 0 {
		return out, fmt.Errorf("no %s containers found, use 'create' command to create it", constants.ContainerName)
	}
	if !strings.Contains(string(out), "running") {
		return out, fmt.Errorf("%s container is not running, use 'create' command to run it", constants.ContainerName)
	}
	return out, nil
}

func (p *provider) fetchProviderInfo() (*providers.ProviderInfo, error) {
	cmd := p.podmanCmd([]string{"info", "--format", "json"})
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

func (p *provider) checkCGroupsAndRootFulMode() error {
	if p.info == nil {
		info, err := p.fetchProviderInfo()
		if err != nil {
			return err
		}
		p.info = info
	}
	if !p.info.CGroupV2 {
		return fmt.Errorf("podman provider requires cgroup v2")
	}
	if p.info.Rootless && !p.allowRootless {
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

func (p *provider) podmanCmd(args []string) exec.Cmd {
	if p.useSudo && runtime.GOOS == "linux" {
		log.Debug("Running with sudo:", "podman", strings.Join(args, " "))
		return exec.Command("sudo", append([]string{"podman"}, args...)...)
	}
	return exec.Command("podman", args...)
}

// rootlessConfigs holds the host paths for all rootless configuration files.
type rootlessConfigs struct {
	microshiftConf string
	crioConf       string
	crunWrapper    string
}

// writeRootlessConfigs creates MicroShift, CRI-O config files and a crun
// wrapper needed for rootless operation. The files are written to the user's
// minc config directory (~/.config/minc/).
func writeRootlessConfigs() (*rootlessConfigs, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(configDir, "minc")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// MicroShift config: enable KubeletInUserNamespace so the kubelet skips
	// host-level sysctl writes and /dev/kmsg checks that fail in user namespaces.
	msPath := filepath.Join(dir, "rootless-microshift.yaml")
	msContent := `kubelet:
  featureGates:
    KubeletInUserNamespace: true
`
	if err := os.WriteFile(msPath, []byte(msContent), 0644); err != nil {
		return nil, fmt.Errorf("writing microshift rootless config: %w", err)
	}

	// CRI-O config: switch cgroup manager from systemd to cgroupfs and
	// register a custom runtime that forces crun --rootless mode.
	// The drop-in is numbered 20- so it overrides 10-microshift.conf.
	crioPath := filepath.Join(dir, "rootless-crio.conf")
	crioContent := `[crio.runtime]
cgroup_manager = "cgroupfs"
conmon_cgroup = "pod"
default_runtime = "crun-rootless"
# Clear default_sysctls set by 10-microshift.conf; pinns cannot write
# net.ipv4.ping_group_range inside rootless user namespaces.
default_sysctls = []

[crio.runtime.runtimes.crun-rootless]
runtime_path = "/usr/local/bin/crun-rootless"
runtime_type = "oci"
`
	if err := os.WriteFile(crioPath, []byte(crioContent), 0644); err != nil {
		return nil, fmt.Errorf("writing crio rootless config: %w", err)
	}

	// crun wrapper: patch the OCI config.json to remove oomScoreAdj and set
	// linux.rootless=true before calling crun. In user namespaces, the kernel
	// only allows INCREASING oom_score_adj; pods requesting a lower value
	// (e.g. -997 for Guaranteed QoS) cause crun to fail with EPERM.
	wrapperPath := filepath.Join(dir, "crun-rootless")
	wrapperContent := `#!/bin/sh
# Patch OCI spec for rootless: strip oomScoreAdj and set linux.rootless.
# Only applies to "create" / "run" subcommands that carry --bundle.
bundle=""
prev=""
for arg in "$@"; do
  if [ "$prev" = "--bundle" ] || [ "$prev" = "-b" ]; then
    bundle="$arg"
    break
  fi
  prev="$arg"
done
if [ -n "$bundle" ] && [ -f "$bundle/config.json" ]; then
  jq 'if .process then .process |= del(.oomScoreAdj) else . end
      | if .linux then .linux.rootless = true else . end' \
    "$bundle/config.json" > "$bundle/config.json.tmp" \
    && mv "$bundle/config.json.tmp" "$bundle/config.json"
fi
exec /usr/bin/crun "$@"
`
	if err := os.WriteFile(wrapperPath, []byte(wrapperContent), 0755); err != nil {
		return nil, fmt.Errorf("writing crun rootless wrapper: %w", err)
	}

	return &rootlessConfigs{
		microshiftConf: msPath,
		crioConf:       crioPath,
		crunWrapper:    wrapperPath,
	}, nil
}
