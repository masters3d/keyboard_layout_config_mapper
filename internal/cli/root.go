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
	Long: `KLCM (Keyboard Layout Configuration Mapper) helps synchronize keyboard layouts
across different keyboard firmware systems including ZMK, QMK, and Kinesis.

Priority focus on ZMK keyboards (Advantage360 & Glove80) with support for:
- Automatic sync between compatible keyboards
- Change detection and conflict resolution  
- GitHub PR automation for upstream changes
- Interactive mode for complex mappings`,
	Example: `  # Sync changes from Advantage360 to Glove80
  klcm sync --from adv360 --to glove80

  # Auto-detect and sync all changes
  klcm sync

  # Preview changes without applying
  klcm sync --preview

  # Create GitHub PRs for all changes
  klcm pr create --all`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.klcm.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}