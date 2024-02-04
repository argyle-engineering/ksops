#!/usr/bin/env bash

set -euo pipefail

# Display help message
display_help() {
	echo "Usage: $0 <version> <installation_directory>"
	echo "  version: Git tag of the release to download (e.g., v1.0.0)"
	echo "  installation_directory: Directory where the software will be installed"
}

# Ensure two parameters are provided
if [ "$#" -ne 2 ]; then
	echo "Error: Invalid number of arguments."
	echo ""
	display_help
	exit 1
fi

# Set the GitHub repository URL
repo_url="argyle-engineering/ksops"

# Set the version and installation directory
version="$1"
install_dir="$2"

# Determine the operating system and architecture
os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

# Download the latest release using gh cli based on the operating system and architecture
if [ "$os" == "darwin" ]; then
	if [ "$arch" == "x86_64" ]; then
		pattern='ksops_darwin_x86_64.tar.gz'
	elif [ "$arch" == "arm64" ]; then
		pattern='ksops_darwin_arm64.tar.gz'
	else
		echo "Unsupported architecture: $arch"
		exit 1
	fi
elif [ "$os" == "linux" ]; then
	if [ "$arch" == "x86_64" ]; then
		pattern='ksops_linux_x86_64.tar.gz'
	elif [ "$arch" == "i386" ]; then
		pattern='ksops_linux_i386.tar.gz'
	elif [ "$arch" == "arm64" ]; then
		pattern='ksops_linux_arm64.tar.gz'
	else
		echo "Unsupported architecture: $arch"
		exit 1
	fi
elif [ "$os" == "windows" ]; then
	echo "Unsupported operating system: $os"
	exit 1
fi

# Create a temporary directory
tmp_dir=$(mktemp -d)

# Cleanup function to remove the temporary directory on exit
cleanup() {
	rm -rf "$tmp_dir"
}
trap cleanup EXIT

# Determine the operating system and architecture
os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

# Download the specified release using gh cli based on the operating system and architecture
gh release download --clobber -R "$repo_url" --pattern "$pattern" -D "$tmp_dir" "$version"

# Extract the downloaded tarball in the temporary directory
tar -xzf "$tmp_dir/ksops_${os}_${arch}.tar.gz" -C "$tmp_dir"

# Move only the binary to the installation directory
mv "$tmp_dir/ksops" "$install_dir"

echo "Installation complete. The software is installed in $install_dir."