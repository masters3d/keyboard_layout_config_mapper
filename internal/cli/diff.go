package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"masters3d.com/keyboard_layout_config_mapper/internal/mappers"
)

var (
	diffKeyboard1 string
	diffKeyboard2 string
	diffSemantic  bool
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff [keyboard1] [keyboard2]",
	Short: "Show differences between keyboard configurations",
	Long: `Show differences between keyboard configurations at various levels:

- Physical layout differences
- Key mapping differences  
- Layer structure differences
- Semantic differences (what functionality changed)

Useful for understanding compatibility and change impact.`,
	Example: `  # Compare two keyboards
  klcm diff adv360 glove80

  # Show semantic differences
  klcm diff --semantic adv360 glove80

  # Compare using flags
  klcm diff --keyboard1 adv360 --keyboard2 glove80`,
	Args: cobra.MaximumNArgs(2),
	RunE: runDiff,
}

func runDiff(cmd *cobra.Command, args []string) error {
	// Handle arguments vs flags
	var k1, k2 string
	
	if len(args) >= 2 {
		k1, k2 = args[0], args[1]
	} else {
		k1, k2 = diffKeyboard1, diffKeyboard2
	}

	if k1 == "" || k2 == "" {
		return fmt.Errorf("two keyboards must be specified")
	}

	if verbose {
		fmt.Printf("🔍 Comparing %s and %s configurations\n", k1, k2)
	}

	differ := mappers.NewDiffer()
	return differ.CompareKeyboards(k1, k2, diffSemantic)
}

func init() {
	rootCmd.AddCommand(diffCmd)

	diffCmd.Flags().StringVar(&diffKeyboard1, "keyboard1", "", "first keyboard to compare")
	diffCmd.Flags().StringVar(&diffKeyboard2, "keyboard2", "", "second keyboard to compare")
	diffCmd.Flags().BoolVar(&diffSemantic, "semantic", false, "show semantic differences (what functionality changed)")
}