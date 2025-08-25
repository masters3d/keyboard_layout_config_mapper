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