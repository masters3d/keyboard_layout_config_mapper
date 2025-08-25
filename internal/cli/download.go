package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
  klcm download --force ergodx     # Force re-download even if file exists`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		
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
		url:      "https://raw.githubusercontent.com/masters3d/glove80-zmk-config/cheyo/config/glove80.keymap",
	},
	{
		name:     "adv360",
		dir:      "configs/zmk_adv360",
		filename: "adv360.keymap",
		url:      "https://raw.githubusercontent.com/masters3d/Adv360-Pro-ZMK/cheyo/config/adv360.keymap",
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
	
	fmt.Println("\n🔗 Repository mapping:")
	fmt.Println("  - Kinesis Advantage 2: masters3d/supportfiles/master/1_qwerty.txt")
	fmt.Println("  - QMK ErgoDox: masters3d/qmk_firmware/masters3d/keyboards/ergodox_ez/keymaps/masters3d/keymap.c")
	fmt.Println("  - Glove80: masters3d/glove80-zmk-config/cheyo/config/glove80.keymap")
	fmt.Println("  - Advantage360: masters3d/Adv360-Pro-ZMK/cheyo/config/adv360.keymap")
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolP("force", "f", false, "Force re-download even if files exist")
}