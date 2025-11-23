#!/usr/bin/env bash
set -e

VM_NAME="directio-vm"
PROJECT_DIR="$(pwd)"
GO_VERSION="1.22.0"

echo "=== Checking Lima installation ==="
if ! command -v limactl >/dev/null 2>&1; then
    echo "Installing Lima..."
    brew install lima
fi

echo "=== Ensuring VM exists ==="
if ! limactl list | grep -q "${VM_NAME}"; then
    echo "Creating new VM: ${VM_NAME}"
    cat <<EOF | limactl start --name="${VM_NAME}" -
images:
  - location: "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-arm64.img"
    arch: "aarch64"

# IMPORTANT: NO writable mounts from macOS!
mounts:
  - location: "~"
    writable: false

cpus: 4
memory: "4GiB"
disk: "20GiB"
EOF
else
    limactl start "${VM_NAME}"
fi

echo "=== Detecting REAL Linux home directory inside VM ==="
VM_HOME=$(limactl shell "${VM_NAME}" -- bash -lc "echo \$HOME")
echo "VM HOME = ${VM_HOME}"

echo "=== Creating project directory ==="
limactl shell "${VM_NAME}" -- bash -lc "mkdir -p ${VM_HOME}/project"

echo "=== Copying ONLY current directory into VM ==="
limactl copy . "${VM_NAME}:${VM_HOME}/project"

echo "=== Installing Go inside VM (if needed) ==="
limactl shell "${VM_NAME}" -- bash -lc "
if ! command -v go >/dev/null 2>&1; then
    wget -q https://go.dev/dl/go${GO_VERSION}.linux-arm64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-arm64.tar.gz
    echo 'export PATH=\$PATH:/usr/local/go/bin' >> ~/.bashrc
fi
"

echo "=== Running benchmarks inside VM ==="
limactl shell "${VM_NAME}" -- bash -lc "
export PATH=\$PATH:/usr/local/go/bin
cd ${VM_HOME}/project
go test -bench . -benchmem -run NONE
"

echo "=== Stopping VM ==="
limactl stop "${VM_NAME}"

echo "=== DONE ==="
