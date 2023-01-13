package schema

type Spec struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	FailSilently bool     `yaml:"fail-silently"`
	Files        []string `yaml:"files"`
}
