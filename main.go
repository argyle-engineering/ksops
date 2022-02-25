package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"go.mozilla.org/sops/v3/cmd/sops/formats"
	"go.mozilla.org/sops/v3/decrypt"
	"sigs.k8s.io/yaml"
)

type resource struct {
	Files        []string `json:"files,omitempty" yaml:"files,omitempty"`
	FailSilently bool     `json:"fail-silently,omitempty" yaml:"fail-silently,omitempty"`
}

func main() {

	if len(os.Args) > 2 {
		_, _ = fmt.Fprintf(os.Stderr, "always invoke this via kustomize plugins, received too many args: %d\n", len(os.Args))
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "always invoke this via kustomize plugins, received too few args: %d\n", len(os.Args))
		os.Exit(1)
	}

	generateDummyEnv := os.Getenv("KSOPS_GENERATE_DUMMY_SECRETS")
	isGenerateDummy, envErr := strconv.ParseBool(generateDummyEnv)
	if envErr != nil && len(generateDummyEnv) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "error converting string to boolean, please use either false or true : %q\n", envErr)
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to read file: %s\n", os.Args[1])
		os.Exit(1)
	}

	var manifest resource
	err = yaml.Unmarshal(content, &manifest)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error unmarshalling manifest content: %q \n%s\n", err, content)
		os.Exit(1)
	}

	if manifest.Files == nil {
		_, _ = fmt.Fprintf(os.Stderr, "missing the required 'files' key in the ksops manifests: %s", content)
		os.Exit(1)
	}

	var output bytes.Buffer

	failed := false

	for _, file := range manifest.Files {

		var b, data []byte

		b, err = ioutil.ReadFile(file)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error reading %q: %q\n", file, err.Error())
			os.Exit(1)
		}

		format := formats.FormatForPath(file)
		data, err = decrypt.DataWithFormat(b, format)

		if err != nil {
			if manifest.FailSilently {
				if isGenerateDummy {
					dummySecret := generateDummySecret(b)
					output.Write(dummySecret)
					output.WriteString("\n---\n")
				} else {
					failed = true
				}
			} else {
				os.Exit(1)
			}
		}

		output.Write(data)
		output.WriteString("\n---\n")
	}

	// if we fail we never output anything
	if failed {
		os.Exit(0)
	}

	_, _ = fmt.Fprintf(os.Stdout, output.String())

}

func generateDummySecret(content []uint8) []byte {
	var yamlData map[string]interface{}
	var dummyString = map[string]string{
		"stringData": "SECRET",
		"data":       "U0VDUkVU",
	}
	yErr := yaml.Unmarshal(content, &yamlData)
	if yErr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error unmarshalling manifest content: %q \n%s\n", yErr, content)
		os.Exit(1)
	}
	delete(yamlData, "sops")
	for k, v := range yamlData {
		if k == "stringData" || k == "data" {
			tempList := make(map[string]string)
			for a, _ := range v.(map[string]interface{}) {
				tempList[a] = dummyString[k]
			}
			yamlData[k] = tempList
		}
	}

	dummyData, err := yaml.Marshal(&yamlData)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %q", err)
	}
	return dummyData
}
