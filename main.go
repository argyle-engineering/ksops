package main

import (
	"github.com/argyle-engineering/ksops/v2/pkg"
	"github.com/spf13/cobra"
	"os"
	"sigs.k8s.io/kustomize/kyaml/errors"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/kio"
)

func main() {
	rootCmd.SetVersionTemplate("{{.Version}}\n")
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ksops",
	Short: "KSOPS is a flexible Kustomize KRM-based plugin for SOPS encrypted resources",
	Long: `KSOPS is a flexible Kustomize KRM-based plugin for SOPS encrypted resources.
- Provides the ability to fail silently if the generator fails to decrypt files.
- Generates dummy secrets with the 'KSOPS_GENERATE_DUMMY_SECRETS' environment variable.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// No config is required
		p := framework.SimpleProcessor{Config: nil, Filter: kio.FilterFunc(pkg.Ksops)}

		// STDIN and STDOUT will be used if no reader or writer respectively is provided.
		err := framework.Execute(p, nil)

		return errors.Wrap(err)
	},
	Version: version,
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ksops.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
