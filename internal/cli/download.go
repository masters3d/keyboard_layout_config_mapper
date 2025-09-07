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

var downloadCmd = &cobra.Command{
	Use:   "download [keyboard...]",
	Short: "Download keyboard configurations from GitHub repositories",
	Long: `Download the latest keyboard configurations from their respective GitHub repositories.

Supported keyboards:
  kinesis2  - Kinesis Advantage 2 configuration
  ergodox   - QMK ErgoDox keymap  
  glove80   - Glove80 ZMK keymap
  adv360    - Advantage360 ZMK keymap

Examples:
  klcm download                    # Download all configurations
  klcm download adv360 glove80     # Download only ZMK keyboards
  klcm download --force ergodx     # Force re-download even if file exists
  klcm download --preview          # Preview changes before downloading`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		preview, _ := cmd.Flags().GetBool("preview")
		
		// Handle preview mode
		if preview {
			fmt.Println("🔍 Previewing changes before download...")
			
			keyboardsToPreview := args
			if len(args) == 0 {
				// Preview all keyboards
				keyboardsToPreview = make([]string, 0, len(keyboards))
				for _, kb := range keyboards {
					keyboardsToPreview = append(keyboardsToPreview, kb.name)
				}
			}
			
			hasChanges := false
			for _, keyboard := range keyboardsToPreview {
				changed, err := previewKeyboardChanges(keyboard)
				if err != nil {
					return fmt.Errorf("failed to preview %s: %w", keyboard, err)
				}
				if changed {
					hasChanges = true
				}
			}
			
			if !hasChanges {
				fmt.Println("\n✅ All local files are up to date. No download needed.")
				return nil
			}
			
			fmt.Print("\n❓ Proceed with download? (y/N): ")
			var response string
			fmt.Scanln(&response)
			
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				fmt.Println("📦 Download cancelled.")
				return nil
			}
			
			fmt.Println("🚀 Starting download...")
			
			// Force download after preview confirmation
			force = true
		}
		
		if len(args) == 0 {
			// Download all if no specific keyboards specified
			return downloadAll(force)
		}
		
		// Download specific keyboards
		for _, keyboard := range args {
			if err := downloadKeyboard(keyboard, force); err != nil {
				return fmt.Errorf("failed to download %s: %w", keyboard, err)
			}
		}
		
		return nil
	},
}

type keyboardConfig struct {
	name     string
	dir      string
	filename string
	url      string
}

var keyboards = []keyboardConfig{
	{
		name:     "kinesis2",
		dir:      "configs/kinesis2",
		filename: "1_qwerty.txt",
		url:      "https://raw.githubusercontent.com/masters3d/supportfiles/master/1_qwerty.txt",
	},
	{
		name:     "ergodox", 
		dir:      "configs/qmk_ergodox",
		filename: "keymap.c",
		url:      "https://raw.githubusercontent.com/masters3d/qmk_firmware/masters3d/keyboards/ergodox_ez/keymaps/masters3d/keymap.c",
	},
	{
		name:     "glove80",
		dir:      "configs/zmk_glove80",
		filename: "glove80.keymap",
		url:      "https://raw.githubusercontent.com/masters3d/keyboard_layout_config_mapper/v7_target/configs/zmk_glove80/glove80.keymap",
	},
	{
		name:     "adv360",
		dir:      "configs/zmk_adv360",
		filename: "adv360.keymap",
		url:      "https://raw.githubusercontent.com/masters3d/keyboard_layout_config_mapper/v7_target/configs/zmk_adv360/adv360.keymap",
	},
	{
		name:     "adv_mod",
		dir:      "configs/zmk_adv_mod",
		filename: "adv_mod.keymap",
		url:      "https://raw.githubusercontent.com/masters3d/zmk-config-pillzmod-nicenano/main/config/adv_mod.keymap",
	},
}

func downloadAll(force bool) error {
	fmt.Println("🚀 Starting keyboard configuration download...")
	
	for _, kb := range keyboards {
		if err := downloadKeyboard(kb.name, force); err != nil {
			return err
		}
	}
	
	fmt.Println("\n🎉 All configurations downloaded successfully!")
	printSummary()
	return nil
}

func downloadKeyboard(name string, force bool) error {
	var kb *keyboardConfig
	for _, k := range keyboards {
		if k.name == name {
			kb = &k
			break
		}
	}
	
	if kb == nil {
		return fmt.Errorf("unknown keyboard: %s", name)
	}
	
	fmt.Printf("\n📁 Processing %s...\n", kb.dir)
	return downloadFile(kb.dir, kb.filename, kb.url, force)
}

func downloadFile(dir, filename, url string, force bool) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	
	filePath := filepath.Join(dir, filename)
	
	// Check if file exists and skip if not forcing
	if !force {
		if _, err := os.Stat(filePath); err == nil {
			fmt.Printf("  ⏭️  %s already exists (use --force to re-download)\n", filename)
			return nil
		}
	}
	
	fmt.Printf("  📥 Downloading %s...\n", filename)
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Download file
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download from %s: %w", url, err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s: HTTP %d", filename, resp.StatusCode)
	}
	
	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()
	
	// Copy content
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}
	
	fmt.Printf("    ✅ Successfully downloaded %s\n", filename)
	return nil
}

func printSummary() {
	fmt.Println("\n📋 Summary:")
	fmt.Println("  - configs/kinesis2/: Kinesis Advantage 2 reference configuration (1_qwerty.txt)")
	fmt.Println("  - configs/qmk_ergodox/: QMK ErgoDox keymap (keymap.c)")
	fmt.Println("  - configs/zmk_glove80/: Glove80 ZMK keymap (glove80.keymap)")
	fmt.Println("  - configs/zmk_adv360/: Advantage360 ZMK keymap (adv360.keymap)")
	fmt.Println("  - configs/zmk_adv_mod/: Kinesis Advantage (Pillz Mod) ZMK keymap (adv_mod.keymap)")
	
	fmt.Println("\n🔗 Repository mapping:")
	fmt.Println("  - Kinesis Advantage 2: masters3d/supportfiles/master/1_qwerty.txt")
	fmt.Println("  - QMK ErgoDox: masters3d/qmk_firmware/masters3d/keyboards/ergodox_ez/keymaps/masters3d/keymap.c")
	fmt.Println("  - Glove80: masters3d/glove80-zmk-config/cheyo/config/glove80.keymap")
	fmt.Println("  - Advantage360: masters3d/Adv360-Pro-ZMK/cheyo/config/adv360.keymap")
}

func previewKeyboardChanges(name string) (bool, error) {
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
	
	filePath := filepath.Join(kb.dir, kb.filename)
	
	// Check if local file exists
	localContent, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("\n📁 %s/%s\n", kb.dir, kb.filename)
			fmt.Printf("  🆕 Local file does not exist - would be created\n")
			return true, nil
		}
		return false, fmt.Errorf("failed to read local file %s: %w", filePath, err)
	}
	
	// Fetch remote content
	fmt.Printf("\n📁 %s/%s\n", kb.dir, kb.filename)
	fmt.Printf("  📡 Checking remote version...\n")
	
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
		fmt.Printf("  ✅ Up to date\n")
		return false, nil
	}
	
	// Show differences using git-style diff
	fmt.Printf("  ⚠️  Changes detected:\n")
	
	localLines := strings.Split(localStr, "\n")
	remoteLines := strings.Split(remoteStr, "\n")
	
	// Generate git-style diff
	opts := DefaultDiffOptions()
	opts.ShowHeader = false // We'll show our own header
	opts.MaxWidth = 100     // Slightly smaller for indented output
	
	diff := UnifiedDiff("local", "remote", localStr, remoteStr, opts)
	
	if diff != "" {
		// Show summary first
		fmt.Printf("    📊 Local: %d lines, Remote: %d lines\n", len(localLines), len(remoteLines))
		
		// Show truncated diff (first few chunks only)
		diffLines := strings.Split(strings.TrimSpace(diff), "\n")
		chunkCount := 0
		maxChunks := 3 // Limit preview chunks
		
		for _, line := range diffLines {
			if strings.HasPrefix(line, "@@") {
				chunkCount++
				if chunkCount > maxChunks {
					fmt.Printf("    ... (%d more changes not shown in preview)\n", 
						strings.Count(diff, "@@")-maxChunks)
					break
				}
			}
			fmt.Printf("    %s\n", line)
		}
		
		// Show summary
		summary := SimpleDiffSummary(localLines, remoteLines)
		fmt.Printf("    📈 Summary: %s\n", summary)
	}
	
	return true, nil
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolP("force", "f", false, "Force re-download even if files exist")
	downloadCmd.Flags().BoolP("preview", "p", false, "Preview changes before downloading")
}