package main

import (
	"fmt"
	"github.com/argyle-engineering/ksops/pkg/dummy"
	"github.com/argyle-engineering/ksops/pkg/schema"
	"go.mozilla.org/sops/v3/cmd/sops/formats"
	"go.mozilla.org/sops/v3/decrypt"
	"os"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	"strconv"
)

func main() {

	ke := os.Getenv("KSOPS_GENERATE_DUMMY_SECRETS")
	if len(ke) == 0 { // env not set
		ke = "false"
	}

	ksopsGenerateDummySecrets, err := strconv.ParseBool(ke)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error converting string to boolean, please use either false or true : %q\n", err)
		os.Exit(1)
	}

	fn := func(items []*yaml.RNode) ([]*yaml.RNode, error) {
		var filteredItems []*yaml.RNode
		for i := range items {
			item := items[i]
			if item.GetKind() == "KSOPSGenerator" && item.GetApiVersion() == "argyle.com/v1" {

				// Get the spec RNode
				rawSpec := item.Field("spec")
				if rawSpec == nil {
					return nil, fmt.Errorf("no spec found in KSOPSGenerator")
				}

				// Get the spec yaml & unmarshal it
				var spec schema.Spec
				err = yaml.Unmarshal([]byte(rawSpec.Value.MustString()), &spec)
				if err != nil {
					return nil, fmt.Errorf("unable to parse KSOPSGenerator spec: %w", err)
				}

				// Generate secrets here
				for _, file := range spec.Files {

					var b, secret []byte
					//var b []byte

					b, err = os.ReadFile(file)

					if err != nil {
						return nil, fmt.Errorf("error reading %s: %w", file, err)
					}

					if ksopsGenerateDummySecrets {
						secret, err = dummy.GenerateDummySecret(b)
						if err != nil {
							return nil, fmt.Errorf("failed generating dummy file %s: %w", file, err)
						}
					} else {
						format := formats.FormatForPath(file)
						secret, err = decrypt.DataWithFormat(b, format)
						if err != nil && !spec.FailSilently {
							return nil, fmt.Errorf("failed decrypting file %s: %w -- %s", file, err, string(secret))
						}
					}

					var node *yaml.RNode
					node, err = yaml.Parse(string(secret))
					if err != nil {
						return nil, fmt.Errorf("failed parse secret into yaml file %s: %w", file, err)
					}

					filteredItems = append(filteredItems, node)
				}
			} else {
				// All other resources get passed along unmodified
				filteredItems = append(filteredItems, item)
			}
		}
		return filteredItems, nil
	}

	api := make(framework.GVKFilterMap)
	api["KSOPSGenerator"] = make(map[string]kio.Filter)
	api["KSOPSGenerator"]["argyle.com/v1"] = kio.FilterFunc(fn)

	p := framework.VersionedAPIProcessor{FilterProvider: api}
	cmd := command.Build(&p, command.StandaloneDisabled, false)
	command.AddGenerateDockerfile(cmd)
	if err = cmd.Execute(); err != nil {
		fmt.Printf("\nerror: %s\n", err)
		os.Exit(1)
	}

}
