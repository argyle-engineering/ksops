package main

import (
	"fmt"
	"github.com/argyle-engineering/ksops/pkg/dummy"
	"github.com/argyle-engineering/ksops/pkg/schema"
	"go.mozilla.org/sops/v3/cmd/sops/formats"
	"go.mozilla.org/sops/v3/decrypt"
	"os"
	"path/filepath"
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
			if item.GetKind() == "ksops" && item.GetApiVersion() == "argyle.com/v1" {

				// Get the spec yaml & unmarshal it
				var spec schema.Spec
				err = yaml.Unmarshal([]byte(item.MustString()), &spec)
				if err != nil {
					return nil, fmt.Errorf("unable to parse ksops spec: %w\n", err)
				}

				// Generate secrets here
				for _, file := range spec.Files {

					var b, secret []byte

					b, err = os.ReadFile(file)

					// Sometimes we end up too deep in a directory
					// So we change to the parent and search there as well
					// #TODO figure out why this happens and fix it
					if err != nil {
						cwd, _ := os.Getwd()
						parent := filepath.Join(cwd, "..")
						b, err = os.ReadFile(filepath.Join(parent, file))
						if err != nil {
							return nil, fmt.Errorf("error opening file for decryption %s: \n\n%w\n\n\n", file, err)
						}
					}

					if ksopsGenerateDummySecrets {
						secret, err = dummy.GenerateDummySecret(b)
						if err != nil {
							return nil, fmt.Errorf("failed generating dummy file %s: %w\n", file, err)
						}
					} else {
						format := formats.FormatForPath(file)
						secret, err = decrypt.DataWithFormat(b, format)
						if err != nil && !spec.FailSilently {
							return nil, fmt.Errorf("failed decrypting file %s: \n\n%w\n\n", file, err)
						}
					}

					var node *yaml.RNode
					node, err = yaml.Parse(string(secret))
					if err != nil {
						return nil, fmt.Errorf("failed parse secret into yaml file %s: %w\n", file, err)
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
	api["ksops"] = make(map[string]kio.Filter)
	api["ksops"]["argyle.com/v1"] = kio.FilterFunc(fn)

	p := framework.VersionedAPIProcessor{FilterProvider: api}
	cmd := command.Build(&p, command.StandaloneDisabled, false)
	command.AddGenerateDockerfile(cmd)
	if err = cmd.Execute(); err != nil {
		fmt.Printf("\nerror: %s\n", err)
		os.Exit(1)
	}

}
