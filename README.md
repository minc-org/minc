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

```bash
curl -LsSf -o minc https://github.com/minc-org/minc/releases/latest/download/minc_linux_amd64
chmod +x minc
```

By default, minc requires sudo to run Podman commands. To run without
sudo, see [Rootless Mode (Linux)](#rootless-mode-linux) below.

#### Mac
```bash
curl -LsSf -o minc https://github.com/minc-org/minc/releases/latest/download/minc_darwin_arm64
chmod +x minc
```

#### Windows
```bash
curl.exe -LsSf -o minc.exe https://github.com/minc-org/minc/releases/latest/download/minc.exe
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

### Regenerate kubeconfig file for cluster
```bash
minc generate-kubeconfig
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
| Parameter            | Description                                                                                                                                           |
|----------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| `microshift-config`  | Custom MicroShift config file to change MicroShift defaults. [More info](https://github.com/openshift/microshift/blob/main/docs/user/howto_config.md) |
| `microshift-version` | MicroShift version, check available tags at `quay.io/minc-org/minc`                                                                                   |
| `log-level`          | Log level (default: `info`)                                                                                                                           |
| `provider`           | Container runtime provider, e.g., `docker`, `podman` (default: `podman`)                                                                              |
| `https-port`         | Different port to use for exposing https service (default:`9443`)                                                                                     |
| `http-port`          | Different port to use for exposing http service (default:`9080`)                                                                                      |
| `allow-rootless`     | Use rootless Podman without sudo (default: `false`). See [Rootless Mode](#rootless-mode-linux)                                                        |
| `disable-overlay-cache` | Disable container overlay storage cache mount (default: `false`)                                                                                  |


Once the container is running, you can interact with the MicroShift cluster using `kubectl` or `oc` tools.

## Rootless Mode (Linux)

Minc can run without sudo using rootless Podman. This is experimental and
requires a one-time host setup.

### Host Prerequisites

Steps 1-3 require root/sudo. All settings persist across reboots.

**1. Kernel parameters (sysctl):**
```bash
sudo tee /etc/sysctl.d/99-minc-rootless.conf <<'EOF'
net.ipv4.ip_forward = 1
net.ipv4.ip_unprivileged_port_start = 80
fs.inotify.max_user_instances = 1024
EOF
sudo sysctl --system
```

- `ip_forward` - required for container networking (packet forwarding)
- `ip_unprivileged_port_start = 80` - allows rootless Podman to bind HTTP/HTTPS ports
- `max_user_instances = 1024` - increases inotify watch limit for MicroShift services

**2. Load `ip_tables` kernel module:**
```bash
sudo tee /etc/modules-load.d/minc-rootless.conf <<'EOF'
ip_tables
EOF
sudo modprobe ip_tables
```

**3. Delegate cgroup controllers to user sessions:**
```bash
sudo mkdir -p /etc/systemd/system/user@.service.d/
sudo tee /etc/systemd/system/user@.service.d/delegate.conf <<'EOF'
[Service]
Delegate=cpuset cpu io memory pids
EOF
sudo systemctl daemon-reload
```

**4. Expand subordinate UID/GID ranges** (required for OpenShift SCC UIDs):
```bash
sudo usermod --add-subuids 165536-1265535999 --add-subgids 165536-1265535999 $(whoami)
```

Then, as your regular user (not root), migrate Podman to pick up the new ranges:
```bash
podman system migrate
```

### Rootless Usage

```bash
# Create a rootless cluster:
minc create --allow-rootless

# Or set it persistently via config:
minc config set allow-rootless true
minc create

# Delete (auto-detects rootless mode):
minc delete
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.
