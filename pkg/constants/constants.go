package constants

import (
	"fmt"
	"runtime"
)

const (
	ContainerName = "microshift"
	HostName      = "127.0.0.1.nip.io"
	LabelKey      = "io.x-openshift.microshift.cluster"
	UShiftVersion = "4.19.0-okd-scos.17"
	Registry      = "quay.io"
	RegistryOrg   = "minc-org"
	ImageName     = "minc"
)

var (
	Version = "dev"
)

func GetImageRegistry() string {
	return fmt.Sprintf("%s/%s/%s", Registry, RegistryOrg, ImageName)
}

func GetUShiftImage(version string) string {
	return fmt.Sprintf("%s:%s-%s", GetImageRegistry(), version, runtime.GOARCH)
}
