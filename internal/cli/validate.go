package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"masters3d.com/keyboard_layout_config_mapper/internal/parsers"
)

var (
	validateAll     bool
	validateCompile bool
	validateKeyboard string
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate keyboard configuration files",
	Long: `Validate keyboard configuration files for syntax errors and compatibility.

Supports validation for:
- ZMK keymap files (.keymap)
- QMK keymap files (.c)  
- Kinesis2 configuration files (.txt)

Can also attempt compilation validation where possible.`,
	Example: `  # Validate all keyboards
  klcm validate --all

  # Validate specific keyboard
  klcm validate --keyboard adv360

  # Validate with compilation check
  klcm validate --compile`,
	RunE: runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	if verbose {
		fmt.Println("✅ Starting keyboard configuration validation...")
	}

	validator := parsers.NewValidator()

	if validateAll || validateKeyboard == "" {
		// Validate all keyboards
		if verbose {
			fmt.Println("🔍 Validating all keyboard configurations...")
		}
		return validator.ValidateAll(validateCompile)
	}

	// Validate specific keyboard
	if verbose {
		fmt.Printf("🎯 Validating %s configuration\n", validateKeyboard)
	}

	return validator.ValidateKeyboard(validateKeyboard, validateCompile)
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().BoolVar(&validateAll, "all", false, "validate all keyboard configurations")
	validateCmd.Flags().BoolVar(&validateCompile, "compile", false, "attempt compilation validation")
	validateCmd.Flags().StringVar(&validateKeyboard, "keyboard", "", "specific keyboard to validate (adv360, glove80, qmk_ergodx, kinesis2)")
}