package constants

import (
	"fmt"
	"runtime"
)

const (
	ContainerName = "microshift"
	HostName      = "127.0.0.1.nip.io"
	LabelKey      = "io.x-openshift.microshift.cluster"
	UShiftVersion = "4.18.0-okd-scos.1"
	Registry      = "quay.io"
	RegistryOrg   = "praveenkumar"
	ImageName     = "microshift-okd"
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
