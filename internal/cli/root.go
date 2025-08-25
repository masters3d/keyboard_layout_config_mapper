package cli

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "klcm",
	Short: "Keyboard Layout Configuration Mapper",
	Long: `KLCM (Keyboard Layout Configuration Mapper) helps manage keyboard configuration files
across different keyboard firmware systems including ZMK, QMK, and Kinesis.

Priority focus on ZMK keyboards (Advantage360 & Glove80) with support for:
- Pulling latest configurations from remote repository
- Comparing local vs remote configurations  
- GitHub PR automation for upstream changes
- Configuration validation and formatting`,
	Example: `  # Pull latest configurations for all keyboards
  klcm pull

  # Preview changes before pulling
  klcm pull --preview

  # Pull specific keyboards
  klcm pull adv360 glove80

  # Compare local vs remote files
  klcm compare-remote

  # Download specific configurations
  klcm download adv360`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.klcm.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}