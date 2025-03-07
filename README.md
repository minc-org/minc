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

### Installation

Get the latest release from GitHub release page as per your platform.

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
