package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"masters3d.com/keyboard_layout_config_mapper/internal/mappers"
)

var (
	syncFrom     string
	syncTo       string
	syncPreview  bool
	syncAll      bool
	interactive  bool
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize keyboard layouts between different keyboards",
	Long: `Sync command helps synchronize keyboard layouts between different keyboards.
	
Primary workflow is ZMK-to-ZMK sync between Advantage360 and Glove80.
Also supports syncing to QMK ErgoDox and Kinesis2 formats.

Auto-detection will find the most recent changes and sync them appropriately.`,
	Example: `  # Auto-detect and sync recent changes
  klcm sync

  # Sync specific direction  
  klcm sync --from adv360 --to glove80

  # Preview changes without applying
  klcm sync --preview

  # Sync to all keyboards
  klcm sync --from adv360 --to all

  # Interactive mode for complex changes
  klcm sync --interactive`,
	RunE: runSync,
}

func runSync(cmd *cobra.Command, args []string) error {
	if verbose {
		fmt.Println("🔄 Starting keyboard layout sync...")
	}

	// Initialize the sync manager
	syncManager := mappers.NewSyncManager()

	// Handle different sync modes
	if interactive {
		return syncManager.SyncInteractive()
	}

	if syncAll || (syncFrom == "" && syncTo == "") {
		// Auto-detect mode
		if verbose {
			fmt.Println("🔍 Auto-detecting changes...")
		}
		return syncManager.SyncAutoDetect(syncPreview)
	}

	// Specific sync direction
	if syncFrom == "" || syncTo == "" {
		return fmt.Errorf("both --from and --to must be specified, or use auto-detect mode")
	}

	if verbose {
		fmt.Printf("🎯 Syncing from %s to %s\n", syncFrom, syncTo)
	}

	return syncManager.SyncDirection(syncFrom, syncTo, syncPreview)
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().StringVar(&syncFrom, "from", "", "source keyboard (adv360, glove80, qmk_ergodx, kinesis2)")
	syncCmd.Flags().StringVar(&syncTo, "to", "", "target keyboard (adv360, glove80, qmk_ergodx, kinesis2, all)")
	syncCmd.Flags().BoolVarP(&syncPreview, "preview", "p", false, "preview changes without applying")
	syncCmd.Flags().BoolVar(&syncAll, "all", false, "sync to all compatible keyboards")
	syncCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "interactive mode for complex mappings")
}