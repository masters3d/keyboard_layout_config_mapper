package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var compareRemoteCmd = &cobra.Command{
	Use:   "compare-remote [keyboard...]",
	Short: "Compare local configuration files with remote versions",
	Long: `Compare your local keyboard configuration files with the latest versions 
from their respective GitHub repositories. This helps you understand what 
changes you would lose before downloading/overriding your local files.

Supported keyboards:
  kinesis2  - Kinesis Advantage 2 configuration
  ergodox   - QMK ErgoDox keymap  
  glove80   - Glove80 ZMK keymap
  adv360    - Advantage360 ZMK keymap

Examples:
  klcm compare-remote                    # Compare all configurations
  klcm compare-remote adv360 glove80     # Compare only ZMK keyboards
  klcm compare-remote --show-unchanged   # Include files with no differences`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showUnchanged, _ := cmd.Flags().GetBool("show-unchanged")
		
		if len(args) == 0 {
			// Compare all if no specific keyboards specified
			return compareAllRemote(showUnchanged)
		}
		
		// Compare specific keyboards
		for _, keyboard := range args {
			_, err := compareKeyboardRemote(keyboard, showUnchanged)
			if err != nil {
				return fmt.Errorf("failed to compare %s: %w", keyboard, err)
			}
		}
		
		return nil
	},
}

func compareAllRemote(showUnchanged bool) error {
	fmt.Println("🔍 Comparing local configurations with remote versions...")
	
	hasChanges := false
	for _, kb := range keyboards {
		changed, err := compareKeyboardRemote(kb.name, showUnchanged)
		if err != nil {
			return err
		}
		if changed {
			hasChanges = true
		}
	}
	
	if !hasChanges {
		fmt.Println("\n✅ All local files are up to date with remote versions!")
	} else {
		fmt.Println("\n💡 Use 'klcm download --force' to update local files with remote changes")
	}
	
	return nil
}

func compareKeyboardRemote(name string, showUnchanged bool) (bool, error) {
	var kb *keyboardConfig
	for _, k := range keyboards {
		if k.name == name {
			kb = &k
			break
		}
	}
	
	if kb == nil {
		return false, fmt.Errorf("unknown keyboard: %s", name)
	}
	
	return compareFileRemote(kb, showUnchanged)
}

func compareFileRemote(kb *keyboardConfig, showUnchanged bool) (bool, error) {
	filePath := filepath.Join(kb.dir, kb.filename)
	
	// Check if local file exists
	localContent, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("\n📁 %s/%s\n", kb.dir, kb.filename)
			fmt.Printf("  🆕 Local file does not exist - would be created by download\n")
			return true, nil
		}
		return false, fmt.Errorf("failed to read local file %s: %w", filePath, err)
	}
	
	// Fetch remote content
	fmt.Printf("\n📁 %s/%s\n", kb.dir, kb.filename)
	fmt.Printf("  📡 Fetching remote version...\n")
	
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	resp, err := client.Get(kb.url)
	if err != nil {
		return false, fmt.Errorf("failed to fetch remote content from %s: %w", kb.url, err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to fetch remote content: HTTP %d", resp.StatusCode)
	}
	
	remoteContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read remote content: %w", err)
	}
	
	// Compare content
	localStr := strings.TrimSpace(string(localContent))
	remoteStr := strings.TrimSpace(string(remoteContent))
	
	if localStr == remoteStr {
		if showUnchanged {
			fmt.Printf("  ✅ No differences found\n")
		}
		return false, nil
	}
	
	// Show differences
	fmt.Printf("  ⚠️  Differences found:\n")
	
	localLines := strings.Split(localStr, "\n")
	remoteLines := strings.Split(remoteStr, "\n")
	
	// Simple diff algorithm - show a few key differences
	showDiff(localLines, remoteLines)
	
	return true, nil
}

func showDiff(localLines, remoteLines []string) {
	// Show summary stats
	fmt.Printf("    📊 Local: %d lines, Remote: %d lines\n", len(localLines), len(remoteLines))
	
	// Show first few differences
	maxLines := len(localLines)
	if len(remoteLines) > maxLines {
		maxLines = len(remoteLines)
	}
	
	diffCount := 0
	maxDiffs := 5 // Limit output
	
	for i := 0; i < maxLines && diffCount < maxDiffs; i++ {
		var localLine, remoteLine string
		
		if i < len(localLines) {
			localLine = strings.TrimSpace(localLines[i])
		}
		if i < len(remoteLines) {
			remoteLine = strings.TrimSpace(remoteLines[i])
		}
		
		if localLine != remoteLine {
			diffCount++
			fmt.Printf("    📝 Line %d:\n", i+1)
			
			if localLine != "" {
				fmt.Printf("      - %s\n", truncateLine(localLine))
			} else {
				fmt.Printf("      - (empty)\n")
			}
			
			if remoteLine != "" {
				fmt.Printf("      + %s\n", truncateLine(remoteLine))
			} else {
				fmt.Printf("      + (empty)\n")
			}
		}
	}
	
	if diffCount == maxDiffs {
		fmt.Printf("    ... (showing first %d differences)\n", maxDiffs)
	}
	
	// Helpful summary
	if len(localLines) != len(remoteLines) {
		diff := len(remoteLines) - len(localLines)
		if diff > 0 {
			fmt.Printf("    📈 Remote has %d more lines\n", diff)
		} else {
			fmt.Printf("    📉 Remote has %d fewer lines\n", -diff)
		}
	}
}

func truncateLine(line string) string {
	if len(line) <= 80 {
		return line
	}
	return line[:77] + "..."
}

func init() {
	rootCmd.AddCommand(compareRemoteCmd)
	compareRemoteCmd.Flags().BoolP("show-unchanged", "u", false, "Show files with no differences")
}