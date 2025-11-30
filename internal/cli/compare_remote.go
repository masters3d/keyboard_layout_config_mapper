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
  glove80   - Glove80 ZMK keymap
  adv360    - Advantage360 ZMK keymap
  adv_mod   - Kinesis Advantage (Pillz Mod) ZMK keymap

Examples:
  klcm compare-remote                    # Compare all configurations
  klcm compare-remote adv360 glove80     # Compare specific keyboards
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
	
	// Generate git-style diff
	localContent := strings.Join(localLines, "\n")
	remoteContent := strings.Join(remoteLines, "\n")
	
	opts := DefaultDiffOptions()
	opts.ShowHeader = false // We'll show our own header
	opts.MaxWidth = 100     // Slightly smaller for indented output
	
	diff := UnifiedDiff("local", "remote", localContent, remoteContent, opts)
	
	if diff != "" {
		// Indent the diff output
		diffLines := strings.Split(strings.TrimSpace(diff), "\n")
		for _, line := range diffLines {
			if strings.HasPrefix(line, "@@") {
				fmt.Printf("    %s\n", line)
			} else {
				fmt.Printf("    %s\n", line)
			}
		}
	}
	
	// Show summary
	summary := SimpleDiffSummary(localLines, remoteLines)
	fmt.Printf("    📈 Summary: %s\n", summary)
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