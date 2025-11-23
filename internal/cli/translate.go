package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"masters3d.com/keyboard_layout_config_mapper/internal/mappers"
	"masters3d.com/keyboard_layout_config_mapper/internal/models"
	"masters3d.com/keyboard_layout_config_mapper/internal/parsers"
)

var (
	translateSourceType string
	translateTargetType string
	translateOutputFile string
	translateShowIR     bool
)

// translateCmd represents the translate command
var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate keyboard layouts between different types via intermediate representation",
	Long: `Translate command converts keyboard layouts from one type to another using the 
intermediate representation (IR) system. This allows automatic translation between 
different keyboard types that share similar layout structures.

The intermediate representation uses a universal 10x10 grid per hand that can map 
to most keyboard layouts, making it easy to port layouts between keyboards.

Examples:
  # Translate from Advantage360 to Glove80
  klcm translate --from adv360 --to glove80

  # Translate from Glove80 to Advanced Mod  
  klcm translate --from glove80 --to adv_mod

  # Save translation to specific file
  klcm translate --from adv360 --to glove80 --output glove80_from_adv360.keymap

  # Show the intermediate representation
  klcm translate --from adv360 --show-ir

Supported keyboard types:
  - adv360: Kinesis Advantage 360 (ZMK)
  - glove80: MoErgo Glove80 (ZMK)  
  - adv_mod: Kinesis Advantage with Pillz Mod (ZMK)`,
	RunE: runTranslate,
}

func runTranslate(cmd *cobra.Command, args []string) error {
	if verbose {
		fmt.Println("🔄 Starting layout translation...")
	}

	// Validate source type
	sourceType := models.KeyboardType(translateSourceType)
	if !isValidKeyboardType(sourceType) {
		return fmt.Errorf("unsupported source keyboard type: %s", translateSourceType)
	}

	// If only showing IR, we don't need a target type
	if !translateShowIR {
		if translateTargetType == "" {
			return fmt.Errorf("target keyboard type is required (use --to flag)")
		}

		targetType := models.KeyboardType(translateTargetType)
		if !isValidKeyboardType(targetType) {
			return fmt.Errorf("unsupported target keyboard type: %s", translateTargetType)
		}

		if sourceType == targetType {
			return fmt.Errorf("source and target types cannot be the same")
		}
	}

	// Load source layout
	fmt.Printf("📁 Loading %s layout...\n", sourceType)
	sourceLayout, err := loadLayout(sourceType)
	if err != nil {
		return fmt.Errorf("failed to load source layout: %w", err)
	}

	// Create layout translator and register mappers
	translator := mappers.NewLayoutTranslator()
	registerAllMappers(translator)

	// Get source mapper
	sourceMapper := getMapper(sourceType)
	if sourceMapper == nil {
		return fmt.Errorf("no mapper available for source type: %s", sourceType)
	}

	// Convert to IR
	fmt.Printf("🔄 Converting %s to intermediate representation...\n", sourceType)
	irLayout, err := sourceMapper.ToIR(sourceLayout)
	if err != nil {
		return fmt.Errorf("failed to convert to IR: %w", err)
	}

	// Show IR if requested
	if translateShowIR {
		fmt.Println("🎯 Intermediate Representation:")
		printIRSummary(irLayout)
		
		if translateOutputFile != "" {
			if err := saveIRToFile(irLayout, translateOutputFile); err != nil {
				return fmt.Errorf("failed to save IR: %w", err)
			}
			fmt.Printf("💾 IR saved to: %s\n", translateOutputFile)
		}
		return nil
	}

	// Convert from IR to target
	targetType := models.KeyboardType(translateTargetType)
	fmt.Printf("🔄 Converting IR to %s...\n", targetType)
	
	targetLayout, err := translator.Translate(sourceLayout, targetType)
	if err != nil {
		return fmt.Errorf("failed to translate layout: %w", err)
	}

	// Determine output file
	outputFile := translateOutputFile
	if outputFile == "" {
		// Generate default output filename
		targetConfigPath, _ := parsers.GetConfigPath(targetType)
		
		dir := filepath.Dir(targetConfigPath)
		filename := fmt.Sprintf("%s_from_%s.keymap", targetType, sourceType)
		outputFile = filepath.Join(dir, filename)
	}

	// Save translated layout
	fmt.Printf("💾 Saving translated layout to: %s\n", outputFile)
	if err := saveLayout(targetLayout, outputFile); err != nil {
		return fmt.Errorf("failed to save translated layout: %w", err)
	}

	// Print summary
	fmt.Printf("✅ Translation complete!\n")
	fmt.Printf("   Source: %s (%d layers, %d keys)\n", 
		sourceType, len(sourceLayout.Layers), countTotalKeys(sourceLayout))
	fmt.Printf("   Target: %s (%d layers, %d keys)\n", 
		targetType, len(targetLayout.Layers), countTotalKeys(targetLayout))
	fmt.Printf("   Output: %s\n", outputFile)

	return nil
}

// loadLayout loads a keyboard layout from its config file
func loadLayout(keyboardType models.KeyboardType) (*models.KeyboardLayout, error) {
	configPath, err := parsers.GetConfigPath(keyboardType)
	if err != nil {
		return nil, err
	}

	parser, err := parsers.NewParser(keyboardType)
	if err != nil {
		return nil, err
	}

	return parser.Parse(configPath)
}

// saveLayout saves a keyboard layout to a file
func saveLayout(layout *models.KeyboardLayout, outputFile string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return err
	}

	// For now, save as JSON. In the future, we could generate the actual keymap format
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(layout)
}

// saveIRToFile saves intermediate representation to a JSON file
func saveIRToFile(irLayout *models.IRLayout, outputFile string) error {
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return err
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(irLayout)
}

// printIRSummary prints a summary of the intermediate representation
func printIRSummary(irLayout *models.IRLayout) {
	fmt.Printf("  Name: %s\n", irLayout.Name)
	fmt.Printf("  Source: %s\n", irLayout.Source)
	fmt.Printf("  Layers: %d\n", len(irLayout.Layers))
	
	// Count keys by hand and zone
	leftKeys := make(map[string]int)
	rightKeys := make(map[string]int)
	
	for _, layer := range irLayout.Layers {
		fmt.Printf("    Layer %d (%s): %d bindings\n", layer.Index, layer.Name, len(layer.Bindings))
		
		for _, binding := range layer.Bindings {
			zone := string(binding.Position.Zone)
			if binding.Position.Hand == "left" {
				leftKeys[zone]++
			} else {
				rightKeys[zone]++
			}
		}
	}
	
	fmt.Printf("  Left hand zones:\n")
	for zone, count := range leftKeys {
		fmt.Printf("    %s: %d keys\n", zone, count)
	}
	
	fmt.Printf("  Right hand zones:\n")
	for zone, count := range rightKeys {
		fmt.Printf("    %s: %d keys\n", zone, count)
	}
	
	if len(irLayout.Combos) > 0 {
		fmt.Printf("  Combos: %d\n", len(irLayout.Combos))
	}
	
	if len(irLayout.Behaviors) > 0 {
		fmt.Printf("  Behaviors: %d\n", len(irLayout.Behaviors))
	}
}

// countTotalKeys counts the total number of key bindings in a layout
func countTotalKeys(layout *models.KeyboardLayout) int {
	total := 0
	for _, layer := range layout.Layers {
		total += len(layer.Bindings)
	}
	return total
}

// isValidKeyboardType checks if a keyboard type is supported for translation
func isValidKeyboardType(keyboardType models.KeyboardType) bool {
	switch keyboardType {
	case models.KeyboardZMKAdv360, models.KeyboardZMKGlove80, models.KeyboardZMKAdvMod:
		return true
	default:
		return false
	}
}

// getMapper returns the appropriate mapper for a keyboard type
func getMapper(keyboardType models.KeyboardType) mappers.KeyboardMapper {
	switch keyboardType {
	case models.KeyboardZMKAdv360:
		return mappers.NewZMKAdv360Mapper()
	case models.KeyboardZMKGlove80:
		return mappers.NewZMKGlove80Mapper()
	case models.KeyboardZMKAdvMod:
		return mappers.NewZMKAdvModMapper()
	default:
		return nil
	}
}

// registerAllMappers registers all available mappers with the translator
func registerAllMappers(translator *mappers.LayoutTranslator) {
	translator.RegisterMapper(mappers.NewZMKAdv360Mapper())
	translator.RegisterMapper(mappers.NewZMKGlove80Mapper())
	translator.RegisterMapper(mappers.NewZMKAdvModMapper())
}

func init() {
	rootCmd.AddCommand(translateCmd)

	translateCmd.Flags().StringVar(&translateSourceType, "from", "", "source keyboard type (required)")
	translateCmd.Flags().StringVar(&translateTargetType, "to", "", "target keyboard type (required unless --show-ir)")
	translateCmd.Flags().StringVarP(&translateOutputFile, "output", "o", "", "output file path (default: auto-generated)")
	translateCmd.Flags().BoolVar(&translateShowIR, "show-ir", false, "show intermediate representation instead of translating")

	translateCmd.MarkFlagRequired("from")
}