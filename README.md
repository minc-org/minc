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
- `kubectl` or `oc` cli tool installed

#### Windows

In windows on wsl environment make sure you have cgroupsv2 enabled which is not the case by default
- https://github.com/spurin/wsl-cgroupsv2
- Kernel command line parameter kernelCommandLine=systemd.unified_cgroup_hierarchy=1 results in creation 
of cgroup V1 and V2 hierarchy. It is prohibited now, microsoft/WSL#6662 (more details around cgroups-v1/v2)

### Installation

Get the latest release from GitHub release page as per your platform.

#### Linux

In Linux, minc require sudo permission to run podman command because it is not working with rootless mode as of now.
Check: https://github.com/minc-org/minc/issues/22

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
curl.exe -L -o minc.exe  https://github.com/minc-org/minc/releases/latest/download/minc.exe
```

## Usage

### Create the cluster 
```bash
minc create
```

### Status of the cluster

This command provide output in `json` format
```bash
minc status
{
  "container": "running",
  "apiserver": "running"
}
```
In case of error output would be look like below
```bash
minc status
{
  "container": "stopped",
  "apiserver": "stopped",
  "error": "no microshift containers found, use 'create' command to create it"
}
```

### Delete the cluster
```bash
minc delete
```

### Get help and options
```bash
minc help
```

### Config options
```bash
minc config -h
```

### Available Config Settings
| Parameter           | Description                                                                                                     |
|---------------------|-----------------------------------------------------------------------------------------------------------------|
| `microshift-config` | Custom MicroShift config file to change MicroShift defaults. [More info](https://github.com/openshift/microshift/blob/main/docs/user/howto_config.md) |
| `microshift-version`| MicroShift version, check available tags at `quay.io/minc-org/minc`                                             |
| `log-level`         | Log level (default: `info`)                                                                                     |
| `provider`          | Container runtime provider, e.g., `docker`, `podman` (default: `podman`)                                        |


Once the container is running, you can interact with the MicroShift cluster using `kubectl` or `oc` tools.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.
