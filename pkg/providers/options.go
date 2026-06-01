package providers

import (
	"fmt"
	"github.com/minc-org/minc/pkg/constants"
)

type COptions struct {
	ContainerName       string
	ImageName           string
	UShiftConfig        string
	HttpPort            int
	HttpsPort           int
	DisableOverlayCache bool
	// HostContainerStorage is the host path to the container engine's graph root (e.g. Podman Store.GraphRoot).
	// When empty, the default rootful path /var/lib/containers/storage is used.
	HostContainerStorage string
	// AllowRootless enables rootless-specific container tweaks:
	//  - sysctl net.ipv6.conf.all.disable_ipv6=1: etcd binds 127.0.0.1 only,
	//    but glibc prefers ::1 for localhost (RFC 6724), causing kube-apiserver
	//    to fail connecting to etcd.
	//  - mount /dev/null as /dev/kmsg: kubelet needs /dev/kmsg for the OOM
	//    watcher, which is inaccessible in rootless user namespaces.
	AllowRootless bool
	// RootlessMicroShiftConfig is the host path to a MicroShift config.d YAML
	// that enables KubeletInUserNamespace and other rootless settings.
	RootlessMicroShiftConfig string
	// RootlessCRIOConfig is the host path to a CRI-O config drop-in that
	// switches the cgroup manager to cgroupfs for rootless operation.
	RootlessCRIOConfig string
	// RootlessCrunWrapper is the host path to a crun wrapper script that
	// forces --rootless mode so crun skips oom_score_adj writes.
	RootlessCrunWrapper string
}

func CreateOptions(r *COptions) []string {
	// in case http or https port is less than 1024 then macOS doesn't allow to bind with 127.0.0.1 so
	// need to bind with all the interfaces
	httpPortOption := "127.0.0.1:%d:80"
	httpsPortOption := "127.0.0.1:%d:443"
	if r.HttpPort < 1024 {
		httpPortOption = "%d:80"
	}
	if r.HttpsPort < 1024 {
		httpsPortOption = "%d:443"
	}
	createOptions := []string{
		"create",
		"--hostname", constants.HostName,
		"--label", fmt.Sprintf("%s=%s", constants.LabelKey, r.ContainerName),
		"-it", "--privileged",
		"-p", fmt.Sprintf(httpPortOption, r.HttpPort),
		"-p", fmt.Sprintf(httpsPortOption, r.HttpsPort),
		"-p", "127.0.0.1:6443:6443",
	}

	if r.AllowRootless {
		createOptions = append(createOptions,
			"--sysctl", "net.ipv6.conf.all.disable_ipv6=1",
			"-v", "/dev/null:/dev/kmsg",
		)
		if r.RootlessMicroShiftConfig != "" {
			createOptions = append(createOptions, "-v",
				fmt.Sprintf("%s:/etc/microshift/config.d/20-rootless.yaml:ro", r.RootlessMicroShiftConfig))
		}
		if r.RootlessCRIOConfig != "" {
			createOptions = append(createOptions, "-v",
				fmt.Sprintf("%s:/etc/crio/crio.conf.d/20-rootless.conf:ro", r.RootlessCRIOConfig))
		}
		if r.RootlessCrunWrapper != "" {
			createOptions = append(createOptions, "-v",
				fmt.Sprintf("%s:/usr/local/bin/crun-rootless:ro", r.RootlessCrunWrapper))
		}
	}

	// Handle overlay cache mount
	if !r.DisableOverlayCache {
		hostPath := "/var/lib/containers/storage"
		if r.HostContainerStorage != "" {
			hostPath = r.HostContainerStorage
		}
		createOptions = append(createOptions, "-v", fmt.Sprintf("%s:/host-container:ro,rshared", hostPath))
	} else {
		// Use named volume for better macOS/Docker compatibility
		// This allows CRI-O to function without accessing host storage
		// Note: Named volumes don't support bind options like 'rshared'
		createOptions = append(createOptions, "-v", "minc-container-storage:/host-container")
	}

	// Mount custom MicroShift config if provided
	if r.UShiftConfig != "" {
		createOptions = append(createOptions, "-v",
			fmt.Sprintf("%s:/etc/microshift/config.d/00-custom-config.yaml:ro,rshared", r.UShiftConfig))
	}

	return append(createOptions,
		"--name", r.ContainerName, r.ImageName)
}

func StartOptions(containerName string) []string {
	return []string{
		"start",
		containerName,
	}
}

func PullOptions(imageName string) []string {
	return []string{
		"pull",
		imageName,
	}
}

func ImageExistOptions(imageName string) []string {
	return []string{
		"image",
		"inspect",
		imageName,
	}
}

func ServiceWaitOption(service, containerName string) []string {
	return []string{
		"exec",
		containerName,
		"systemctl",
		"is-active",
		service,
	}
}

func KubeConfigOption(containerName, hostname string) []string {
	return []string{
		"exec",
		containerName,
		"cat",
		fmt.Sprintf("/var/lib/microshift/resources/kubeadmin/%s/kubeconfig", hostname),
	}
}

func DeleteOptions(containerName string) []string {
	return []string{
		"rm",
		"-f",
		containerName,
	}
}

func ListOptions(containerName string) []string {
	return []string{
		"ps",
		"-a",
		"-f", fmt.Sprintf("label=%s=%s", constants.LabelKey, containerName),
		"--format", "{{.Names}} {{.Ports}} {{.State}}",
	}
}
