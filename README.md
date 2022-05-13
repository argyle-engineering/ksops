# KSOPS 

*Inspired by: https://github.com/viaduct-ai/kustomize-sops*

A Flexible Kustomize Plugin for SOPS Encrypted Resources.

The main difference in this *fork* is the ability to fail silently to allow this plugin to be used in a CI or places 
where you don't want to allow for decryption.

##  Installation

Build and install binary to the following dir:

`$XDG_CONFIG_HOME/kustomize/plugin/argyle.com/v1/ksops/`

## Usage

Run kustomize with plugins enabled flag:

`kustomize build --enable-alpha-plugins <dir>`


## Manifest example

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

## Generate Dummy Secret

There is a case where our machine does not have an access to the decryptor key (i.e. our CICD), but we still want ksops to keep producing a secret with a placeholder 'secret' value in it.

With `fail-silently` set to `true`, ksops will outputing an `failed decrypting file` error message with 0 (zero) exit code. 

In order to generate a dummy secret without error message above, we need set `KSOPS_GENERATE_DUMMY_SECRETS` environment variable to `TRUE`. e.g `KSOPS_GENERATE_DUMMY_SECRETS=TRUE kustomize build --enable-alpha-plugins <dir>`