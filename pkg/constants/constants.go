package constants

import (
	"fmt"
	"runtime"
)

const (
	ContainerName = "microshift"
	HostName      = "127.0.0.1.nip.io"
	LabelKey      = "io.x-openshift.microshift.cluster"
)

var (
	ImageName = fmt.Sprintf("quay.io/praveenkumar/microshift-okd:4.18.0-okd-scos.1-%s", runtime.GOARCH)
)
