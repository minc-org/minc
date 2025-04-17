package providers

import (
	"fmt"
	"github.com/minc-org/minc/pkg/constants"
)

func RunOptions(containerName, imageName, uShiftConfig string) []string {
	runOptions := []string{
		"run",
		"--hostname", constants.HostName,
		"--label", fmt.Sprintf("%s=%s", constants.LabelKey, containerName),
		"--detach",
		"-it", "--privileged",
		"-v", "/var/lib/containers/storage:/host-container:ro,rshared",
		"-p", "127.0.0.1:9080:80",
		"-p", "127.0.0.1:9443:443",
		"-p", "127.0.0.1:6443:6443",
	}
	if uShiftConfig != "" {
		runOptions = append(runOptions, "-v",
			fmt.Sprintf("%s:/etc/microshift/config.d/00-custom-config.yaml:ro,rshared", uShiftConfig))
	}
	return append(runOptions,
		"--name", containerName, imageName)
}

func PullOptions(imageName string) []string {
	return []string{
		"pull",
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
		"-s",
		"-f", fmt.Sprintf("label=%s=%s", constants.LabelKey, containerName),
		"--format", "{{.Names}} {{.Ports}}",
	}
}
