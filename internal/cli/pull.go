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
	"masters3d.com/keyboard_layout_config_mapper/internal/models"
	"masters3d.com/keyboard_layout_config_mapper/internal/parsers"
)

var (
	pullPreview bool
	pullAll     bool
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull [keyboards...]",
	Short: "Pull latest configuration files from remote repository",
	Long: `Pull command downloads the latest configuration files from the remote repository.
	
Similar to 'git pull', this updates your local configurations with the latest remote versions.
Supports preview mode to see changes before applying them.`,
	Example: `  # Pull updates for all keyboards
  klcm pull

  # Pull updates for specific keyboards
  klcm pull adv360 glove80

  # Preview changes before pulling
  klcm pull --preview

  # Preview changes for specific keyboards  
  klcm pull --preview adv360`,
	RunE: runPull,
}

func runPull(cmd *cobra.Command, args []string) error {
	if verbose {
		fmt.Println("🔄 Starting configuration pull...")
	}

	// Determine which keyboards to pull
	keyboards := args
	if len(keyboards) == 0 || pullAll {
		// Default to all supported keyboards
		keyboards = []string{"adv360", "glove80", "adv_mod"}
	}

	if pullPreview {
		fmt.Println("🔍 Previewing changes...")
		fmt.Println("📋 Preview mode - no changes will be applied")
	}

	for _, keyboard := range keyboards {
		keyboardType := models.KeyboardType(keyboard)
		
		// Get local file path
		configPath, err := parsers.GetConfigPath(keyboardType)
		if err != nil {
			fmt.Printf("❌ Failed to get config path for %s: %v\n", keyboard, err)
			continue
		}

		if err := pullKeyboard(keyboardType, configPath, pullPreview); err != nil {
			fmt.Printf("❌ Failed to pull %s: %v\n", keyboard, err)
			continue
		}
	}

	return nil
}

func pullKeyboard(keyboardType models.KeyboardType, configPath string, preview bool) error {
	// Get remote URL for this keyboard
	remoteURL, err := getRemoteURL(keyboardType)
	if err != nil {
		return fmt.Errorf("failed to get remote URL: %v", err)
	}

	fmt.Printf("📁 %s\n", filepath.Base(configPath))
	fmt.Println("  📡 Fetching remote version...")

	// Fetch remote content
	remoteContent, err := fetchRemoteContent(remoteURL)
	if err != nil {
		return fmt.Errorf("failed to fetch remote content: %v", err)
	}

	// Read local content if it exists
	var localContent string
	if _, err := os.Stat(configPath); err == nil {
		localBytes, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to read local file: %v", err)
		}
		localContent = string(localBytes)
	} else {
		localContent = "" // File doesn't exist locally
	}

	// Compare content
	if localContent == remoteContent {
		fmt.Println("  ✅ Already up to date")
		return nil
	}

	// Show differences
	fmt.Println("  ⚠️  Changes detected:")
	localLines := strings.Split(localContent, "\n")
	remoteLines := strings.Split(remoteContent, "\n")
	
	fmt.Printf("    📊 Local: %d lines, Remote: %d lines\n", len(localLines), len(remoteLines))
	
	// Generate and show git-style diff
	diff := UnifiedDiff("local", "remote", localContent, remoteContent, DefaultDiffOptions())
	if diff != "" {
		fmt.Print(diff)
	}

	if preview {
		return nil // Don't actually apply changes in preview mode
	}

	// Ask for confirmation
	fmt.Print("\n❓ Apply changes? (y/N): ")
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))
	
	if response != "y" && response != "yes" {
		fmt.Println("❌ Pull cancelled")
		return nil
	}

	// Apply changes
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(configPath, []byte(remoteContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Println("  ✅ Changes applied successfully")
	return nil
}

func getRemoteURL(keyboardType models.KeyboardType) (string, error) {
	baseURL := "https://raw.githubusercontent.com/masters3d/keyboard_layout_config_mapper/main/configs"
	
	switch keyboardType {
	case models.KeyboardZMKAdv360:
		return baseURL + "/zmk_adv360/adv360.keymap", nil
	case models.KeyboardZMKGlove80:
		return baseURL + "/zmk_glove80/glove80.keymap", nil
	case models.KeyboardZMKAdvMod:
		return "https://raw.githubusercontent.com/masters3d/zmk-config-pillzmod-nicenano/cheyo/config/pillzmod_pro.keymap", nil
	default:
		return "", fmt.Errorf("unsupported keyboard type: %s", keyboardType)
	}
}

func fetchRemoteContent(url string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}
	
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	return string(content), nil
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().BoolVarP(&pullPreview, "preview", "p", false, "preview changes without applying")
	pullCmd.Flags().BoolVar(&pullAll, "all", false, "pull updates for all keyboards")
}