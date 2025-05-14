package providers

import (
	"fmt"
	"github.com/minc-org/minc/pkg/constants"
)

type COptions struct {
	ContainerName string
	ImageName     string
	UShiftConfig  string
	HttpPort      int
	HttpsPort     int
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
		"-v", "/var/lib/containers/storage:/host-container:ro,rshared",
		"-p", fmt.Sprintf(httpPortOption, r.HttpPort),
		"-p", fmt.Sprintf(httpsPortOption, r.HttpsPort),
		"-p", "127.0.0.1:6443:6443",
	}
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
		"--format", "{{.Names}} {{.Ports}}",
	}
}
