package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync [source] [target]",
	Short: "Sync configuration changes between different keyboards",
	Long: `Sync allows you to copy configuration changes from one keyboard to another.
This is useful when you make changes to one keyboard layout and want to apply
similar changes to another keyboard.

Examples:
  # Compare adv360 and glove80 configurations
  klcm sync adv360 glove80 --preview

  # Apply changes from adv360 to glove80
  klcm sync adv360 glove80

  # Show available keyboards
  klcm sync --list`,
	Args: cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		listKeyboards, _ := cmd.Flags().GetBool("list")
		preview, _ := cmd.Flags().GetBool("preview")

		if listKeyboards {
			return showAvailableKeyboards()
		}

		if len(args) < 2 {
			return fmt.Errorf("please specify source and target keyboards (e.g., 'klcm sync adv360 glove80')")
		}

		source := args[0]
		target := args[1]

		return syncKeyboards(source, target, preview)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().Bool("list", false, "List available keyboards for syncing")
	syncCmd.Flags().Bool("preview", false, "Preview changes without applying them")
}

func showAvailableKeyboards() error {
	fmt.Println("📋 Available keyboards for syncing:")
	fmt.Println("  • adv360    - Kinesis Advantage360 (ZMK)")
	fmt.Println("  • glove80   - MoErgo Glove80 (ZMK)")
	fmt.Println("  • qmk_ergodox - ErgoDox (QMK)")
	fmt.Println("  • kinesis2  - Kinesis2 (Kinesis)")
	fmt.Println()
	fmt.Println("💡 Example: klcm sync adv360 glove80 --preview")
	return nil
}

func syncKeyboards(source, target string, preview bool) error {
	// Get file paths
	sourcePath := getKeyboardConfigPath(source)
	targetPath := getKeyboardConfigPath(target)

	if sourcePath == "" {
		return fmt.Errorf("unknown source keyboard: %s", source)
	}
	if targetPath == "" {
		return fmt.Errorf("unknown target keyboard: %s", target)
	}

	// Check if this is a cross-format sync (ZMK to Kinesis2 or vice versa)
	if isKinesis2(target) && !isKinesis2(source) {
		return syncZMKToKinesis2(source, target, preview, sourcePath, targetPath)
	}
	if isKinesis2(source) && !isKinesis2(target) {
		return syncKinesis2ToZMK(source, target, preview, sourcePath, targetPath)
	}

	// Read source file
	sourceContent, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file %s: %w", sourcePath, err)
	}

	// Read target file
	targetContent, err := os.ReadFile(targetPath)
	if err != nil {
		return fmt.Errorf("failed to read target file %s: %w", targetPath, err)
	}

	// Find the default layer in both files
	sourceDefault := findDefaultLayer(string(sourceContent))
	targetDefault := findDefaultLayer(string(targetContent))

	if sourceDefault == "" {
		return fmt.Errorf("could not find default layer in %s", source)
	}
	if targetDefault == "" {
		return fmt.Errorf("could not find default layer in %s", target)
	}

	// Compare and show differences
	fmt.Printf("🔄 Comparing %s → %s\n\n", source, target)
	
	differences := compareDefaultLayers(sourceDefault, targetDefault, source, target)
	
	if len(differences) == 0 {
		fmt.Println("✅ No differences found between default layers")
		return nil
	}

	fmt.Printf("📊 Found %d key differences:\n\n", len(differences))
	for i, diff := range differences {
		fmt.Printf("%d. Position %d:\n", i+1, diff.Position)
		fmt.Printf("   %s: %s\n", source, diff.SourceKey)
		fmt.Printf("   %s: %s\n", target, diff.TargetKey)
		fmt.Println()
	}

	if preview {
		fmt.Println("👀 Preview mode - no changes applied")
		return nil
	}

	// Ask for confirmation
	fmt.Printf("❓ Apply changes from %s to %s? (y/N): ", source, target)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response != "y" && response != "yes" {
		fmt.Println("❌ Sync cancelled")
		return nil
	}

	// Apply changes
	updatedContent := applyChangesToTarget(string(targetContent), targetDefault, sourceDefault)
	
	err = os.WriteFile(targetPath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated file: %w", err)
	}

	fmt.Printf("✅ Successfully synced changes from %s to %s\n", source, target)
	return nil
}

func getKeyboardConfigPath(keyboard string) string {
	paths := map[string]string{
		"adv360":     "configs/zmk_adv360/adv360.keymap",
		"glove80":    "configs/zmk_glove80/glove80.keymap",
		"qmk_ergodx": "configs/qmk_ergodox/keymap.c",
		"kinesis2":   "configs/kinesis2/1_qwerty.txt",
	}
	return paths[keyboard]
}

func findDefaultLayer(content string) string {
	lines := strings.Split(content, "\n")
	inDefaultLayer := false
	var layerLines []string

	for _, line := range lines {
		// Look for default layer start
		if strings.Contains(line, "layer0_default") || strings.Contains(line, "default_layer") {
			inDefaultLayer = true
			layerLines = append(layerLines, line)
			continue
		}

		if inDefaultLayer {
			layerLines = append(layerLines, line)
			
			// Look for end of layer (next layer or closing brace)
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "layer") && !strings.Contains(line, "layer0_default") {
				// Remove the last line as it's the start of next layer
				layerLines = layerLines[:len(layerLines)-1]
				break
			}
			if trimmed == "};" || trimmed == "}" {
				break
			}
		}
	}

	return strings.Join(layerLines, "\n")
}

type KeyDifference struct {
	Position  int
	SourceKey string
	TargetKey string
}

func compareDefaultLayers(sourceLayer, targetLayer, sourceName, targetName string) []KeyDifference {
	sourceKeys := extractKeybindings(sourceLayer)
	targetKeys := extractKeybindings(targetLayer)

	var differences []KeyDifference
	maxLen := len(sourceKeys)
	if len(targetKeys) > maxLen {
		maxLen = len(targetKeys)
	}

	for i := 0; i < maxLen; i++ {
		var sourceKey, targetKey string
		if i < len(sourceKeys) {
			sourceKey = sourceKeys[i]
		}
		if i < len(targetKeys) {
			targetKey = targetKeys[i]
		}

		if sourceKey != targetKey && sourceKey != "" && targetKey != "" {
			differences = append(differences, KeyDifference{
				Position:  i + 1,
				SourceKey: sourceKey,
				TargetKey: targetKey,
			})
		}
	}

	return differences
}

func extractKeybindings(layerContent string) []string {
	var keys []string
	lines := strings.Split(layerContent, "\n")
	
	for _, line := range lines {
		// Look for lines with key bindings
		if strings.Contains(line, "&kp ") || strings.Contains(line, "&mo ") || strings.Contains(line, "&to ") {
			// Extract individual key bindings from the line
			parts := strings.Fields(line)
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "&") {
					// Clean up the key binding
					key := strings.TrimSuffix(part, ",")
					key = strings.TrimSuffix(key, ";")
					if key != "" {
						keys = append(keys, key)
					}
				}
			}
		}
	}
	
	return keys
}

func applyChangesToTarget(targetContent, targetDefault, sourceDefault string) string {
	sourceKeys := extractKeybindings(sourceDefault)
	targetKeys := extractKeybindings(targetDefault)

	// For now, let's do a simple replacement of the specific key that changed
	// This is a simplified version - in a full implementation, you'd want more sophisticated mapping
	
	// Find what changed in source vs target
	for i := 0; i < len(sourceKeys) && i < len(targetKeys); i++ {
		if sourceKeys[i] != targetKeys[i] {
			// Replace the target key with source key
			targetContent = strings.Replace(targetContent, targetKeys[i], sourceKeys[i], 1)
		}
	}

	return targetContent
}

// Helper functions for kinesis2 support
func isKinesis2(keyboard string) bool {
	return keyboard == "kinesis2"
}

func syncZMKToKinesis2(source, target string, preview bool, sourcePath, targetPath string) error {
	fmt.Printf("🔄 Cross-format sync: %s (ZMK) → %s (Kinesis2)\n\n", source, target)

	// Read both files
	sourceContent, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	targetContent, err := os.ReadFile(targetPath)
	if err != nil {
		return fmt.Errorf("failed to read target file: %w", err)
	}

	// Analyze ZMK changes and find kinesis2 equivalents
	changes := analyzeZMKChanges(string(sourceContent), source)
	kinesisChanges := mapZMKToKinesis2(changes)

	if len(kinesisChanges) == 0 {
		fmt.Println("✅ No relevant changes found to sync to kinesis2")
		return nil
	}

	fmt.Printf("📊 Found %d mappable changes:\n\n", len(kinesisChanges))
	for i, change := range kinesisChanges {
		fmt.Printf("%d. %s\n", i+1, change.Description)
		fmt.Printf("   Current: %s\n", change.CurrentMapping)
		fmt.Printf("   Proposed: %s\n", change.ProposedMapping)
		fmt.Println()
	}

	if preview {
		fmt.Println("👀 Preview mode - no changes applied")
		return nil
	}

	// Ask for confirmation
	fmt.Printf("❓ Apply changes to %s? (y/N): ", target)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response != "y" && response != "yes" {
		fmt.Println("❌ Sync cancelled")
		return nil
	}

	// Apply changes to kinesis2 file
	updatedContent := applyKinesis2Changes(string(targetContent), kinesisChanges)
	
	err = os.WriteFile(targetPath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated file: %w", err)
	}

	fmt.Printf("✅ Successfully synced %s changes to %s\n", source, target)
	return nil
}

func syncKinesis2ToZMK(source, target string, preview bool, sourcePath, targetPath string) error {
	return fmt.Errorf("kinesis2 to ZMK sync not yet implemented")
}

type ZMKChange struct {
	Position    string
	OldKey      string
	NewKey      string
	Description string
}

type Kinesis2Change struct {
	LineNumber      int
	CurrentMapping  string
	ProposedMapping string
	Description     string
}

func analyzeZMKChanges(content, keyboard string) []ZMKChange {
	var changes []ZMKChange

	// Enhanced ZMK change detection for adv360
	if keyboard == "adv360" {
		// Check for ESCAPE→BSPC change
		if strings.Contains(content, "&kp BSPC        &kp Q") {
			changes = append(changes, ZMKChange{
				Position:    "top-left",
				OldKey:      "ESCAPE",
				NewKey:      "BSPC",
				Description: "Top-left key changed from ESCAPE to BSPC",
			})
		}

		// Check for TAB→BSPC change in thumb cluster
		if strings.Contains(content, "[kp-tab]>[bspace]") || strings.Contains(content, "[tab]>[bspace]") {
			changes = append(changes, ZMKChange{
				Position:    "left-thumb-tab",
				OldKey:      "ESCAPE",
				NewKey:      "BSPC",
				Description: "Tab key remapped to Backspace",
			})
		}

		// Check for shift key F18/F19 mappings
		if strings.Contains(content, "[lshift]>[F18]") {
			changes = append(changes, ZMKChange{
				Position:    "left-shift",
				OldKey:      "LSHIFT",
				NewKey:      "F18",
				Description: "Left Shift mapped to F18",
			})
		}

		if strings.Contains(content, "[rshift]>[F19]") {
			changes = append(changes, ZMKChange{
				Position:    "right-shift",
				OldKey:      "RSHIFT", 
				NewKey:      "F19",
				Description: "Right Shift mapped to F19",
			})
		}

		// Check for morphed keys (behaviors)
		if strings.Contains(content, "morph_dot") {
			changes = append(changes, ZMKChange{
				Position:    "right-home-row",
				OldKey:      "PERIOD",
				NewKey:      "PERIOD/COLON",
				Description: "Period key has mod-morph to Colon when shifted",
			})
		}

		if strings.Contains(content, "morph_comma") {
			changes = append(changes, ZMKChange{
				Position:    "right-bottom-row",
				OldKey:      "COMMA",
				NewKey:      "COMMA/SEMICOLON",
				Description: "Comma key has mod-morph to Semicolon when shifted",
			})
		}

		// Check for quote morphs
		if strings.Contains(content, "morph_quote_single") {
			changes = append(changes, ZMKChange{
				Position:    "top-row-quote",
				OldKey:      "SINGLE_QUOTE",
				NewKey:      "SINGLE_QUOTE/GRAVE",
				Description: "Single quote has mod-morph to Grave when shifted",
			})
		}
	}

	return changes
}

func mapZMKToKinesis2(zmkChanges []ZMKChange) []Kinesis2Change {
	var kinesisChanges []Kinesis2Change

	for _, change := range zmkChanges {
		switch {
		case change.Position == "top-left" && change.NewKey == "BSPC":
			// Top-left escape key changed to backspace
			kinesisChanges = append(kinesisChanges, Kinesis2Change{
				LineNumber:      67, // [kp-tab]>[bspace] line  
				CurrentMapping:  "[kp-tab]>[escape]",
				ProposedMapping: "[kp-tab]>[bspace]",
				Description:     "Map Keypad Tab to Backspace (matches ZMK top-left change)",
			})
			kinesisChanges = append(kinesisChanges, Kinesis2Change{
				LineNumber:      68, // [tab]>[bspace] line
				CurrentMapping:  "[tab]>[escape]",
				ProposedMapping: "[tab]>[bspace]",
				Description:     "Map Tab key to Backspace (matches ZMK top-left change)",
			})

		case change.Position == "left-shift" && change.NewKey == "F18":
			// Left shift mapped to F18
			kinesisChanges = append(kinesisChanges, Kinesis2Change{
				LineNumber:      73, // [lshift]>[F18] line
				CurrentMapping:  "[lshift]>[null]",
				ProposedMapping: "[lshift]>[F18]",
				Description:     "Map Left Shift to F18 (matches ZMK)",
			})
			kinesisChanges = append(kinesisChanges, Kinesis2Change{
				LineNumber:      74, // [kp-lshift]>[F18] line
				CurrentMapping:  "[kp-lshift]>[null]",
				ProposedMapping: "[kp-lshift]>[F18]",
				Description:     "Map Keypad Left Shift to F18 (matches ZMK)",
			})

		case change.Position == "right-shift" && change.NewKey == "F19":
			// Right shift mapped to F19 - already exists in the config
			// No changes needed as it's already [rshift]>[F19]

		case change.Position == "right-home-row" && change.NewKey == "PERIOD/COLON":
			// Period with colon mod-morph - add macro for shift+period
			kinesisChanges = append(kinesisChanges, Kinesis2Change{
				LineNumber:      168, // After the semicolon mappings
				CurrentMapping:  "",
				ProposedMapping: "{lshift}{.}>{speed9}{-lshift}{;}{+lshift}",
				Description:     "Add macro for Shift+Period → Colon (matches ZMK mod-morph)",
			})

		case change.Position == "right-bottom-row" && change.NewKey == "COMMA/SEMICOLON":
			// Comma with semicolon mod-morph - add macro for shift+comma  
			kinesisChanges = append(kinesisChanges, Kinesis2Change{
				LineNumber:      183, // After the slash mappings
				CurrentMapping:  "",
				ProposedMapping: "{lshift}{,}>{speed9}{-lshift}{;}{+lshift}",
				Description:     "Add macro for Shift+Comma → Semicolon (matches ZMK mod-morph)",
			})

		case change.Position == "top-row-quote" && change.NewKey == "SINGLE_QUOTE/GRAVE":
			// Single quote with grave mod-morph
			kinesisChanges = append(kinesisChanges, Kinesis2Change{
				LineNumber:      83, // After [1]>['] line
				CurrentMapping:  "",
				ProposedMapping: "{lshift}{'}>{speed9}{`}",
				Description:     "Add macro for Shift+SingleQuote → Grave (matches ZMK mod-morph)",
			})
		}
	}

	return kinesisChanges
}

func applyKinesis2Changes(content string, changes []Kinesis2Change) string {
	lines := strings.Split(content, "\n")

	for _, change := range changes {
		if change.CurrentMapping == "" {
			// This is a new line to add
			if change.LineNumber-1 <= len(lines) {
				// Insert new line at the specified position
				lines = append(lines[:change.LineNumber-1], 
					append([]string{change.ProposedMapping}, lines[change.LineNumber-1:]...)...)
			} else {
				// Append at the end
				lines = append(lines, change.ProposedMapping)
			}
		} else {
			// This is a replacement of existing line
			if change.LineNumber-1 < len(lines) {
				// Check if the current line matches what we expect
				currentLine := strings.TrimSpace(lines[change.LineNumber-1])
				expectedLine := strings.TrimSpace(change.CurrentMapping)
				if currentLine == expectedLine || strings.Contains(currentLine, expectedLine) {
					lines[change.LineNumber-1] = change.ProposedMapping
				} else {
					// Find the line with the expected content
					for i, line := range lines {
						if strings.TrimSpace(line) == expectedLine || strings.Contains(line, expectedLine) {
							lines[i] = change.ProposedMapping
							break
						}
					}
				}
			}
		}
	}

	return strings.Join(lines, "\n")
}