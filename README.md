# KSOPS 

A Flexible Kustomize KRM based Plugin for SOPS Encrypted Resources.

This is a completely new KRM based plugin with no affiliation with the [existing Go-based ksops plugin](https://github.com/viaduct-ai/kustomize-sops).

##  Installation
Download the binary and add it to your path.

## Fail silently (in case you want the generator to just skip files that it fails to decrypt)
To allow it to fail silently just add the following to your generator:

```yaml
apiVersion: argyle.com/v1
kind: ksops
metadata:
  name: secret-generator
fail-silently: true
files:
- ./secret.enc.yaml
```

## Dummy Secrets 

In order to generate a dummy secrets, we need set `KSOPS_GENERATE_DUMMY_SECRETS` environment variable to `true`.
e.g `KSOPS_GENERATE_DUMMY_SECRETS=TRUE kustomize build --enable-alpha-plugins <dir>`_


## Example usage:
If you want to test ksops without having to do a bunch of setup, you can use the example files and pgp key provided with the repository:

Install gpg and sops and kustomize using brew (or figure it out if you're on Linux)
```shell
brew install sops gnugpg kustomize
```

then:

```shell
gpg --import example/sops_functional_tests_key.asc
kustomize build --enable-alpha-plugins --enable-exec example/
```

This last step will decrypt example.yaml using the test private key.