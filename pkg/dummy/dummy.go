package dummy

import (
	"fmt"
	"sigs.k8s.io/yaml"
)

func GenerateDummySecret(content []uint8) ([]byte, error) {
	var genericYaml map[string]interface{}
	var dummyString = map[string]string{
		"stringData": "SECRET",
		"data":       "U0VDUkVU",
	}
	err := yaml.Unmarshal(content, &genericYaml)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall in dummy: %w", err)
	}
	delete(genericYaml, "sops")
	for k, v := range genericYaml {
		if k == "stringData" || k == "data" {
			tempList := make(map[string]string)
			for a := range v.(map[string]interface{}) {
				tempList[a] = dummyString[k]
			}
			genericYaml[k] = tempList
		}
	}

	secretBytes, err := yaml.Marshal(&genericYaml)

	if err != nil {
		return nil, fmt.Errorf("failed to marshall in dummy: %w", err)
	}

	return secretBytes, nil
}
