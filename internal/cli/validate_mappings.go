package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"masters3d.com/keyboard_layout_config_mapper/internal/mappers"
)

// validateMappingsCmd represents the validate-mappings command
var validateMappingsCmd = &cobra.Command{
	Use:   "validate-mappings",
	Short: "Validate intermediate representation mappings for consistency",
	Long: `Validate-mappings command checks that all keyboard mappers have consistent
bidirectional mappings between keyboard-specific coordinates and the intermediate
representation (IR) coordinate system.

This ensures that layout translations will work correctly and that no key positions
are lost during translation between keyboards.

Examples:
  # Validate all keyboard mappings
  klcm validate-mappings

  # Verbose output showing detailed mapping information
  klcm validate-mappings --verbose`,
	RunE: runValidateMappings,
}

func runValidateMappings(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Validating intermediate representation mappings...")

	// Create all mappers
	adv360Mapper := mappers.NewZMKAdv360Mapper()
	glove80Mapper := mappers.NewZMKGlove80Mapper()
	advModMapper := mappers.NewZMKAdvModMapper()

	keyboards := []struct {
		name   string
		mapper mappers.KeyboardMapper
	}{
		{"Advantage360", adv360Mapper},
		{"Glove80", glove80Mapper},
		{"Advanced Mod", advModMapper},
	}

	allValid := true

	for _, kb := range keyboards {
		fmt.Printf("\n📁 Validating %s mapping...\n", kb.name)

		toIRMapping := kb.mapper.GetPositionMapping()
		fromIRMapping := kb.mapper.GetReverseMapping()

		if verbose {
			fmt.Printf("  📊 Keyboard positions: %d\n", len(toIRMapping))
			fmt.Printf("  📊 IR positions: %d\n", len(fromIRMapping))
		}

		// Validate bidirectional consistency
		if err := mappers.ValidateMapping(toIRMapping, fromIRMapping); err != nil {
			fmt.Printf("  ❌ Validation failed: %v\n", err)
			allValid = false
			continue
		}

		// Check for missing mappings
		if len(toIRMapping) != len(fromIRMapping) {
			fmt.Printf("  ⚠️  Mapping count mismatch: %d keyboard positions, %d IR positions\n",
				len(toIRMapping), len(fromIRMapping))
		}

		if verbose {
			// Show some example mappings
			fmt.Printf("  📝 Sample mappings:\n")
			count := 0
			for kbPos, irPos := range toIRMapping {
				if count >= 5 {
					break
				}
				fmt.Printf("    %s → %s (zone: %s)\n", kbPos, irPos.IRCoordinate(), irPos.Zone)
				count++
			}
		}

		fmt.Printf("  ✅ %s mapping is valid\n", kb.name)
	}

	// Cross-keyboard compatibility check
	fmt.Printf("\n🔄 Checking cross-keyboard compatibility...\n")

	// Get all IR positions used by each keyboard
	adv360IR := make(map[string]bool)
	glove80IR := make(map[string]bool)
	advModIR := make(map[string]bool)

	for _, irPos := range adv360Mapper.GetPositionMapping() {
		adv360IR[irPos.IRCoordinate()] = true
	}

	for _, irPos := range glove80Mapper.GetPositionMapping() {
		glove80IR[irPos.IRCoordinate()] = true
	}

	for _, irPos := range advModMapper.GetPositionMapping() {
		advModIR[irPos.IRCoordinate()] = true
	}

	// Find common positions
	commonPositions := 0
	for irCoord := range adv360IR {
		if glove80IR[irCoord] && advModIR[irCoord] {
			commonPositions++
		}
	}

	fmt.Printf("  📊 Common IR positions across all keyboards: %d\n", commonPositions)
	fmt.Printf("  📊 Adv360-specific positions: %d\n", len(adv360IR)-commonPositions)
	fmt.Printf("  📊 Glove80-specific positions: %d\n", len(glove80IR)-commonPositions)
	fmt.Printf("  📊 AdvMod-specific positions: %d\n", len(advModIR)-commonPositions)

	if commonPositions > 0 {
		fmt.Printf("  ✅ Cross-keyboard translation possible (%d shared positions)\n", commonPositions)
	} else {
		fmt.Printf("  ⚠️  No common positions found - translations may be limited\n")
		allValid = false
	}

	// Final result
	fmt.Printf("\n🎯 Validation Summary:\n")
	if allValid {
		fmt.Printf("✅ All mappings are valid and consistent!\n")
		fmt.Printf("🔄 Translation between keyboards should work correctly.\n")
	} else {
		fmt.Printf("❌ Some mappings have issues that need to be addressed.\n")
		return fmt.Errorf("mapping validation failed")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(validateMappingsCmd)
}