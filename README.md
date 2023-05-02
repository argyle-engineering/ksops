# KSOPS

KSOPS is a flexible Kustomize KRM-based plugin for SOPS encrypted resources. This repository provides a completely new KRM-based plugin with no affiliation with the existing Go-based ksops plugin.

## Features

- A flexible Kustomize KRM-based plugin for SOPS encrypted resources.
- Provides the ability to fail silently if the generator fails to decrypt files.
- Generates dummy secrets with the `KSOPS_GENERATE_DUMMY_SECRETS` environment variable.
- Example files and PGP key are provided with the repository to test KSOPS.

## Installation

To install KSOPS, download the binary and add it to your path. 

Additionally, if you are using non-KRM version, you also need to set the `XDG_CONFIG_HOME` environment variable in your shell. If the variable is not set, run the following command:

```shell
echo "export XDG_CONFIG_HOME=\$HOME/.config" >> $HOME/(.zshrc|.bashrc)
source $HOME/(.zshrc|.bashrc)
```

## Usage

To use KSOPS, follow these steps:

1. Import the GPG key: `gpg --import example/sops_functional_tests_key.asc`.
2. Build and decrypt the example files: `kustomize build --enable-alpha-plugins --enable-exec example/`.

To generate dummy secrets, set the `KSOPS_GENERATE_DUMMY_SECRETS` environment variable to `true`. For example: `KSOPS_GENERATE_DUMMY_SECRETS=TRUE kustomize build --enable-alpha-plugins <dir>`.

To allow KSOPS to fail silently, add the following to the generator:

```yaml
apiVersion: argyle.com/v1
kind: ksops
metadata:
  name: secret-generator
fail-silently: true
files:
- ./secret.enc.yaml
```

## Development

To release a new version, install `goreleaser` and set your GitHub token:

```shell
brew install goreleaser syft 
export GITHUB_TOKEN="YOUR_GH_TOKEN"
```

Then, create a tag and push it to GitHub:

```shell
git tag -a v0.1.0
git push origin v0.1.0
```

Finally, run the following command:

```shell
goreleaser release
```

or use docker 

```shell
docker buildx build --platform linux/arm64,linux/amd64 -t ksops:v1.0.3 --push .
```

## Build Status

The repository has the following badges to indicate the status of the build:

- Lint Status: [![Lint Status](https://github.com/argyle-engineering/ksops/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/argyle-engineering/ksops/actions/workflows/golangci-lint.yml)
- Fmt Status: [![fmt Status](https://github.com/argyle-engineering/ksops/actions/workflows/fmt.yaml/badge.svg)](https://github.com/argyle-engineering/ksops/actions/workflows/fmt.yaml)


