# MINC: MicroShift in Container

MINC enables the deployment of [MicroShift](https://github.com/openshift/microshift), a lightweight OpenShift/Kubernetes distribution, within [Podman](https://podman.io/) as container.
This approach facilitates a streamlined and efficient environment for developing, testing, and running cloud-native applications.

## Features

- Containerized deployment of MicroShift using Podman.
- Simplified setup for lightweight Kubernetes clusters.
- Ideal for development, testing, and edge computing scenarios.

## Getting Started

### Prerequisites

- [Podman](https://podman.io/getting-started/installation) installed on your system.
- Basic understanding of Kubernetes and containerization concepts.

#### Windows

In windows on wsl environment make sure you have cgroupsv2 enabled which is not the case by default
- https://github.com/spurin/wsl-cgroupsv2
- Kernel command line parameter kernelCommandLine=systemd.unified_cgroup_hierarchy=1 results in creation 
of cgroup V1 and V2 hierarchy. It is prohibited now, microsoft/WSL#6662 (more details around cgroups-v1/v2)

### Installation

Get the latest release from GitHub release page as per your platform.

#### Linux

```bash
curl -L -o minc  https://github.com/minc-org/minc/releases/latest/download/minc_linux_amd64
chmod +x minc
```

#### Mac
```bash
curl -L -o minc  https://github.com/minc-org/minc/releases/latest/download/minc_darwin_arm64
chmod +x minc
```

#### Windows
```bash
curl -L -o minc.exe  https://github.com/minc-org/minc/releases/latest/download/minc.exe
```

## Usage

### Create the cluster 
```bash
minc create
```

### Delete the cluster
```bash
minc delete
```

### Get help and options
```bash
minc help
```
Once the container is running, you can interact with the MicroShift cluster using `kubectl` or `oc` tools.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.
