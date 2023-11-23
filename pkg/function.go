package pkg

import (
	"fmt"
	"go.mozilla.org/sops/v3/cmd/sops/formats"
	"go.mozilla.org/sops/v3/decrypt"
	"os"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	"strconv"
	"strings"
)

var ksopsGenerateDummySecrets bool

func init() {
	var err error
	ke := os.Getenv("KSOPS_GENERATE_DUMMY_SECRETS")
	if len(ke) == 0 { // env not set
		ke = "false"
	}

	ksopsGenerateDummySecrets, err = strconv.ParseBool(ke)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error converting string to boolean, please use either false or true : %q\n", err)
		os.Exit(1)
	}
}

func Ksops(items []*yaml.RNode) ([]*yaml.RNode, error) {
	var filteredItems []*yaml.RNode
	for i := range items {
		item := items[i]

		// All other resources get passed along unmodified
		if strings.ToLower(item.GetKind()) != "ksops" || strings.ToLower(item.GetApiVersion()) != "argyle.com/v1" {
			filteredItems = append(filteredItems, item)
			continue
		}

		// Get the spec yaml & unmarshal it
		var spec Spec
		err := yaml.Unmarshal([]byte(item.MustString()), &spec)
		if err != nil {
			return nil, fmt.Errorf("unable to parse ksops spec: %w\n", err)
		}

		// Generate secrets here
		for _, file := range spec.Files {

			var b, secret []byte

			b, err = os.ReadFile(file)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w\n", file, err)
			}

			if ksopsGenerateDummySecrets {
				secret, err = GenerateDummySecret(b)
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
			_ = node.SetAnnotations(map[string]string{"kustomize.config.k8s.io/needs-hash": "true"})
			if err != nil {
				return nil, fmt.Errorf("failed parse secret into yaml file %s: %w\n", file, err)
			}

			filteredItems = append(filteredItems, node)
		}

	}
	return filteredItems, nil
}
