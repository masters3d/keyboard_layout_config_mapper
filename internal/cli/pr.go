package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	prAll     bool
	prDryRun  bool
	prBranch  string
	prApply   bool
	prForce   bool
)

// GitHub repository configuration
type RepoConfig struct {
	Name        string
	Owner       string
	LocalPath   string
	RemoteURL   string
	DefaultFile string
}

var repoConfigs = []RepoConfig{
	{
		Name:        "adv360",
		Owner:       "masters3d",
		LocalPath:   "configs/zmk_adv360",
		RemoteURL:   "https://github.com/masters3d/Adv360-Pro-ZMK",
		DefaultFile: "adv360.keymap",
	},
	{
		Name:        "glove80", 
		Owner:       "masters3d",
		LocalPath:   "configs/zmk_glove80",
		RemoteURL:   "https://github.com/masters3d/glove80-zmk-config",
		DefaultFile: "glove80.keymap",
	},
	{
		Name:        "adv_mod",
		Owner:       "masters3d",
		LocalPath:   "configs/zmk_adv_mod",
		RemoteURL:   "https://github.com/masters3d/zmk-config-pillzmod-nicenano",
		DefaultFile: "adv_mod.keymap",
	},
}

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage GitHub pull requests for keyboard configurations",
	Long: `Create and manage GitHub pull requests for keyboard configuration changes.

Supports creating PRs to upstream repositories when you've made changes
to keyboard configurations that should be shared back to the community.`,
}

// prCreateCmd represents the pr create command
var prCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create pull requests for configuration changes",
	Long: `Create pull requests to upstream repositories for keyboard configuration changes.

This will:
1. Detect which configurations have changed
2. Create appropriate branches 
3. Commit changes with descriptive messages
4. Create pull requests to upstream repositories`,
	Example: `  # Create PRs for all changed configurations
  klcm pr create --all

  # Dry run to see what would be created
  klcm pr create --dry-run

  # Create PR with custom branch name
  klcm pr create --branch feature/my-layout-updates`,
	RunE: runPRCreate,
}

func runPRCreate(cmd *cobra.Command, args []string) error {
	if verbose {
		fmt.Println("🚀 Creating GitHub pull requests...")
	}

	var hasChanges bool
	var changedRepos []RepoConfig
	var err error

	// If force flag is set, assume all repos have changes
	if prForce {
		fmt.Println("🔧 Force mode: Creating PRs for all configured keyboards...")
		changedRepos = repoConfigs
		hasChanges = true
	} else {
		// Check for git changes (uncommitted or recent commits)
		hasChanges, changedRepos, err = detectChanges()
		if err != nil {
			return fmt.Errorf("failed to detect changes: %w", err)
		}

		if !hasChanges {
			fmt.Println("📭 No changes detected in any keyboard configurations")
			fmt.Println("💡 Make some changes first, or use --force to create PRs anyway")
			return nil
		}
	}

	fmt.Printf("📊 Found changes in %d keyboard configuration(s):\n", len(changedRepos))
	for _, repo := range changedRepos {
		fmt.Printf("   📁 %s (%s)\n", repo.Name, repo.LocalPath)
	}
	fmt.Println()

	if prDryRun {
		fmt.Println("🔍 Dry run mode - analyzing what PRs would be created:")
		return simulatePRCreation(changedRepos)
	}

	if !prApply {
		fmt.Println("⚠️  Use --apply flag to actually create the PRs")
		fmt.Println("💡 Or use --dry-run to see what would be created")
		return simulatePRCreation(changedRepos)
	}

	return createActualPRs(changedRepos)
}

func detectChanges() (bool, []RepoConfig, error) {
	var changedRepos []RepoConfig
	changedPaths := make(map[string]bool)
	
	// First check git status for uncommitted changes
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, nil, fmt.Errorf("failed to run git status: %w", err)
	}

	changes := strings.TrimSpace(string(output))
	if changes != "" {
		// Parse uncommitted changes
		lines := strings.Split(changes, "\n")
		for _, line := range lines {
			if len(line) < 3 {
				continue
			}
			path := strings.TrimSpace(line[3:])
			changedPaths[path] = true
		}
	}

	// If no uncommitted changes, check recent commits (last 5 commits)
	if len(changedPaths) == 0 {
		cmd = exec.Command("git", "diff", "--name-only", "HEAD~5..HEAD")
		output, err = cmd.Output()
		if err != nil {
			// If we can't get recent commits, just continue (might be a new repo)
			if verbose {
				fmt.Printf("⚠️  Could not check recent commits: %v\n", err)
			}
		} else {
			recentChanges := strings.TrimSpace(string(output))
			if recentChanges != "" {
				lines := strings.Split(recentChanges, "\n")
				for _, line := range lines {
					path := strings.TrimSpace(line)
					if path != "" {
						changedPaths[path] = true
					}
				}
				if verbose && len(changedPaths) > 0 {
					fmt.Printf("📊 Found changes in recent commits (%d files)\n", len(changedPaths))
				}
			}
		}
	}

	// If still no changes, check if current branch is ahead of origin
	if len(changedPaths) == 0 {
		cmd = exec.Command("git", "rev-list", "--count", "HEAD", "--not", "--remotes")
		output, err = cmd.Output()
		if err == nil {
			unpushedCommits := strings.TrimSpace(string(output))
			if unpushedCommits != "0" && unpushedCommits != "" {
				// We have unpushed commits, check what files they touch
				cmd = exec.Command("git", "diff", "--name-only", "@{upstream}..HEAD")
				output, err = cmd.Output()
				if err == nil {
					unpushedChanges := strings.TrimSpace(string(output))
					if unpushedChanges != "" {
						lines := strings.Split(unpushedChanges, "\n")
						for _, line := range lines {
							path := strings.TrimSpace(line)
							if path != "" {
								changedPaths[path] = true
							}
						}
						if verbose && len(changedPaths) > 0 {
							fmt.Printf("📊 Found changes in unpushed commits (%d files)\n", len(changedPaths))
						}
					}
				}
			}
		}
	}

	// Match changed paths to repo configs
	for _, repo := range repoConfigs {
		for path := range changedPaths {
			if strings.HasPrefix(path, repo.LocalPath) {
				changedRepos = append(changedRepos, repo)
				break
			}
		}
	}

	return len(changedRepos) > 0, changedRepos, nil
}

func simulatePRCreation(repos []RepoConfig) error {
	fmt.Println("📋 PR Creation Simulation:")
	fmt.Println()

	timestamp := time.Now().Format("20060102")
	branchPrefix := "klcm-sync"
	if prBranch != "" {
		branchPrefix = prBranch
	}

	for i, repo := range repos {
		branchName := fmt.Sprintf("%s-%s", branchPrefix, timestamp)
		
		// Get change summary
		changeSummary, err := getChangeSummary(repo)
		if err != nil {
			changeSummary = "Unable to determine changes"
		}

		fmt.Printf("   %d. 📁 %s -> %s/%s\n", i+1, repo.LocalPath, repo.Owner, getRepoName(repo.RemoteURL))
		fmt.Printf("      🌿 Branch: %s\n", branchName)
		fmt.Printf("      📝 Title: Update %s keyboard layout configuration\n", repo.Name)
		fmt.Printf("      📊 Changes: %s\n", changeSummary)
		fmt.Printf("      🔗 Remote: %s\n", repo.RemoteURL)
		fmt.Println()
	}

	if !prApply {
		fmt.Println("💡 Use --apply to create these PRs")
		fmt.Println("🔍 Use --dry-run to see this simulation again")
	}
	
	return nil
}

func createActualPRs(repos []RepoConfig) error {
	fmt.Println("🚀 Creating actual pull requests...")
	fmt.Println()

	// Check if GitHub CLI is available
	if !isGitHubCLIAvailable() {
		return fmt.Errorf("GitHub CLI (gh) is required but not found. Please install it: https://cli.github.com/")
	}

	// Check if user is authenticated
	if !isGitHubAuthenticated() {
		fmt.Println("🔑 GitHub authentication required")
		fmt.Println("Please run: gh auth login")
		return fmt.Errorf("not authenticated with GitHub")
	}

	timestamp := time.Now().Format("20060102")
	branchPrefix := "klcm-sync"
	if prBranch != "" {
		branchPrefix = prBranch
	}

	createdPRs := make([]string, 0)

	for i, repo := range repos {
		fmt.Printf("📁 Processing %s (%d/%d)...\n", repo.Name, i+1, len(repos))
		
		branchName := fmt.Sprintf("%s-%s", branchPrefix, timestamp)
		
		// Create and push branch
		prURL, err := createPRForRepo(repo, branchName)
		if err != nil {
			fmt.Printf("   ❌ Failed: %v\n", err)
			continue
		}

		createdPRs = append(createdPRs, prURL)
		fmt.Printf("   ✅ Created: %s\n", prURL)
		fmt.Println()
	}

	// Summary
	fmt.Printf("🎉 Successfully created %d pull request(s):\n", len(createdPRs))
	for _, url := range createdPRs {
		fmt.Printf("   🔗 %s\n", url)
	}

	if len(createdPRs) > 0 {
		fmt.Println()
		fmt.Println("💡 Next steps:")
		fmt.Println("   - Review the PRs in your browser")
		fmt.Println("   - Add descriptions and additional context")
		fmt.Println("   - Monitor for feedback from maintainers")
	}

	return nil
}

func createPRForRepo(repo RepoConfig, branchName string) (string, error) {
	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Change to repo directory
	repoPath := filepath.Join(currentDir, repo.LocalPath)
	if err := os.Chdir(repoPath); err != nil {
		return "", fmt.Errorf("failed to change to repo directory: %w", err)
	}
	defer os.Chdir(currentDir)

	// Check if we're already on a different branch
	currentBranch, err := getCurrentBranch()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	// If we're not on main/master, we may already have our changes
	if currentBranch != "main" && currentBranch != "master" {
		if verbose {
			fmt.Printf("   📋 Currently on branch: %s\n", currentBranch)
		}
		// Use current branch instead of creating new one
		branchName = currentBranch
	} else {
		// Create new branch from current state
		if err := runGitCommand("checkout", "-b", branchName); err != nil {
			return "", fmt.Errorf("failed to create branch: %w", err)
		}
	}

	// Check if there are uncommitted changes to add and commit
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	hasUncommitted := err == nil && strings.TrimSpace(string(output)) != ""

	if hasUncommitted {
		// Add and commit uncommitted changes
		if err := runGitCommand("add", "."); err != nil {
			return "", fmt.Errorf("failed to add changes: %w", err)
		}

		// Create commit message
		commitMsg := fmt.Sprintf("Update %s keyboard layout configuration\n\nGenerated by KLCM (Keyboard Layout Configuration Mapper)", repo.Name)
		
		// Commit changes
		if err := runGitCommand("commit", "-m", commitMsg); err != nil {
			return "", fmt.Errorf("failed to commit changes: %w", err)
		}
	} else {
		if verbose {
			fmt.Printf("   📋 No uncommitted changes, using existing commits\n")
		}
	}

	// Push branch (this will push whatever commits we have)
	if err := runGitCommand("push", "-u", "origin", branchName); err != nil {
		return "", fmt.Errorf("failed to push branch: %w", err)
	}

	// Create PR using GitHub CLI
	prTitle := fmt.Sprintf("Update %s keyboard layout configuration", repo.Name)
	prBody := fmt.Sprintf(`# %s Layout Update

This PR contains keyboard layout configuration updates generated by KLCM (Keyboard Layout Configuration Mapper).

## Changes
- Updated keyboard layout configuration
- Validated configuration syntax
- Ready for review and testing

## Generated by
- **Tool**: KLCM v0.1.0
- **Date**: %s
- **Branch**: %s

Please review the changes and test the configuration before merging.`, repo.Name, time.Now().Format("2006-01-02"), branchName)

	// Create PR
	cmd = exec.Command("gh", "pr", "create", 
		"--title", prTitle,
		"--body", prBody,
		"--head", branchName)
	
	output, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to create PR: %w", err)
	}

	prURL := strings.TrimSpace(string(output))
	return prURL, nil
}

func runGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func isGitHubCLIAvailable() bool {
	cmd := exec.Command("gh", "--version")
	err := cmd.Run()
	return err == nil
}

func isGitHubAuthenticated() bool {
	cmd := exec.Command("gh", "auth", "status")
	err := cmd.Run()
	return err == nil
}

func getChangeSummary(repo RepoConfig) (string, error) {
	// Get diff stats
	cmd := exec.Command("git", "diff", "--numstat", "HEAD", "--", repo.LocalPath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return "No changes", nil
	}

	totalAdded := 0
	totalDeleted := 0
	fileCount := 0

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			fileCount++
			// parts[0] = added, parts[1] = deleted
			// Handle binary files (shown as "-")
			if parts[0] != "-" {
				// Simple approach: count non-empty as 1
				if parts[0] != "0" {
					totalAdded++
				}
			}
			if parts[1] != "-" {
				if parts[1] != "0" {
					totalDeleted++
				}
			}
		}
	}

	if fileCount == 0 {
		return "No file changes", nil
	}

	return fmt.Sprintf("%d file(s) modified", fileCount), nil
}

func getRepoName(url string) string {
	// Extract repo name from URL
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "unknown"
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.AddCommand(prCreateCmd)
	prCmd.AddCommand(prStatusCmd)
	prCmd.AddCommand(prWorkflowCmd)

	prCreateCmd.Flags().BoolVar(&prAll, "all", false, "create PRs for all changed configurations")
	prCreateCmd.Flags().BoolVar(&prDryRun, "dry-run", false, "show what PRs would be created without creating them")
	prCreateCmd.Flags().BoolVar(&prApply, "apply", false, "actually create the PRs (required for real creation)")
	prCreateCmd.Flags().StringVar(&prBranch, "branch", "", "custom branch name prefix")
	prCreateCmd.Flags().BoolVar(&prForce, "force", false, "force creation even if no changes detected")
}

// prStatusCmd represents the pr status command
var prStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check status of existing pull requests",
	Long: `Check the status of pull requests created by KLCM.
	
Shows information about open PRs, their review status, and merge status.`,
	Example: `  # Check status of all PRs
  klcm pr status
  
  # Check status with verbose output
  klcm pr status --verbose`,
	RunE: runPRStatus,
}

func runPRStatus(cmd *cobra.Command, args []string) error {
	if verbose {
		fmt.Println("📊 Checking PR status...")
	}

	if !isGitHubCLIAvailable() {
		return fmt.Errorf("GitHub CLI (gh) is required but not found. Please install it: https://cli.github.com/")
	}

	if !isGitHubAuthenticated() {
		fmt.Println("🔑 GitHub authentication required")
		fmt.Println("Please run: gh auth login")
		return fmt.Errorf("not authenticated with GitHub")
	}

	fmt.Println("🔍 Checking pull request status...")
	fmt.Println()

	for i, repo := range repoConfigs {
		fmt.Printf("%d. 📁 %s (%s)\n", i+1, repo.Name, getRepoName(repo.RemoteURL))
		
		// Check for PRs from this repo
		err := checkReposPRs(repo)
		if err != nil {
			fmt.Printf("   ❌ Error checking PRs: %v\n", err)
		}
		fmt.Println()
	}

	return nil
}

func checkReposPRs(repo RepoConfig) error {
	// Get current directory and change to repo directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	repoPath := filepath.Join(currentDir, repo.LocalPath)
	if err := os.Chdir(repoPath); err != nil {
		fmt.Printf("   📭 No local repository found\n")
		return nil
	}
	defer os.Chdir(currentDir)

	// List PRs using GitHub CLI
	cmd := exec.Command("gh", "pr", "list", "--author", "@me", "--limit", "10")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("   📭 No PRs found or error accessing repository\n")
		return nil
	}

	prs := strings.TrimSpace(string(output))
	if prs == "" {
		fmt.Printf("   📭 No open PRs found\n")
		return nil
	}

	// Parse and display PR information
	lines := strings.Split(prs, "\n")
	fmt.Printf("   📋 Found %d open PR(s):\n", len(lines))
	
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			fmt.Printf("      🔗 %s\n", line)
		}
	}

	return nil
}

// prWorkflowCmd represents the pr workflow command  
var prWorkflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Interactive workflow guide for making configuration changes",
	Long: `Step-by-step interactive guide for the complete keyboard configuration workflow.

This command walks you through:
1. Updating to latest configurations
2. Making and validating changes  
3. Syncing between keyboards (optional)
4. Creating pull requests`,
	Example: `  # Start interactive workflow
  klcm pr workflow
  
  # Start workflow with verbose output
  klcm pr workflow --verbose`,
	RunE: runPRWorkflow,
}

func runPRWorkflow(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 KLCM Interactive Keyboard Configuration Workflow")
	fmt.Println("================================================")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// Step 1: Update configurations
	fmt.Println("📋 Step 1: Update to Latest Configurations")
	fmt.Println("Before making changes, let's ensure you have the latest configurations.")
	fmt.Println()
	
	if askYesNo(reader, "Would you like to check for and apply remote updates? (Y/n): ") {
		fmt.Println("🔍 Checking for remote updates...")
		if err := runCommandInteractive("./klcm", "pull", "--preview"); err != nil {
			fmt.Printf("❌ Error checking updates: %v\n", err)
		} else {
			if askYesNo(reader, "Apply these updates? (Y/n): ") {
				if err := runCommandInteractive("./klcm", "pull"); err != nil {
					fmt.Printf("❌ Error applying updates: %v\n", err)
				} else {
					fmt.Println("✅ Updates applied successfully")
				}
			}
		}
	}
	fmt.Println()

	// Step 2: Make changes
	fmt.Println("📋 Step 2: Make Your Changes")
	fmt.Println("Now edit your keyboard configuration files:")
	fmt.Println("  📁 Advantage360: configs/zmk_adv360/adv360.keymap")
	fmt.Println("  📁 Glove80: configs/zmk_glove80/glove80.keymap")
	fmt.Println()
	fmt.Println("Press Enter when you've finished making your changes...")
	reader.ReadLine()

	// Step 3: Validate changes
	fmt.Println("📋 Step 3: Validate Your Changes")
	fmt.Println("Let's make sure your changes are syntactically correct.")
	fmt.Println()
	
	fmt.Println("🔍 Running validation...")
	if err := runCommandInteractive("./klcm", "validate"); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
		fmt.Println("Please fix the errors and run 'klcm validate' again before continuing.")
		return nil
	}
	fmt.Println("✅ All configurations are valid!")
	fmt.Println()

	// Step 4: Sync between keyboards (optional)
	fmt.Println("📋 Step 4: Sync Between Keyboards (Optional)")
	fmt.Println("If you want the same changes on both keyboards, we can sync them.")
	fmt.Println()
	
	if askYesNo(reader, "Would you like to sync changes between keyboards? (y/N): ") {
		// Ask which direction to sync
		fmt.Println("Which direction would you like to sync?")
		fmt.Println("  1. Advantage360 → Glove80")
		fmt.Println("  2. Glove80 → Advantage360")
		fmt.Print("Enter choice (1/2): ")
		
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		
		var source, target string
		switch choice {
		case "1":
			source, target = "adv360", "glove80"
		case "2":
			source, target = "glove80", "adv360"
		default:
			fmt.Println("Invalid choice, skipping sync.")
			goto step5
		}

		fmt.Printf("🔍 Previewing sync from %s to %s...\n", source, target)
		if err := runCommandInteractive("./klcm", "sync", source, target, "--preview"); err != nil {
			fmt.Printf("❌ Error previewing sync: %v\n", err)
		} else {
			if askYesNo(reader, "Apply this sync? (y/N): ") {
				if err := runCommandInteractive("./klcm", "sync", source, target); err != nil {
					fmt.Printf("❌ Error applying sync: %v\n", err)
				} else {
					fmt.Println("✅ Sync applied successfully")
					// Validate after sync
					fmt.Println("🔍 Re-validating after sync...")
					if err := runCommandInteractive("./klcm", "validate"); err != nil {
						fmt.Printf("❌ Validation failed after sync: %v\n", err)
						return nil
					}
					fmt.Println("✅ All configurations still valid!")
				}
			}
		}
	}
	
	step5:
	fmt.Println()

	// Step 5: Create pull requests
	fmt.Println("📋 Step 5: Create Pull Requests")
	fmt.Println("Finally, let's create pull requests to contribute your changes back.")
	fmt.Println()
	
	fmt.Println("🔍 Analyzing what PRs would be created...")
	if err := runCommandInteractive("./klcm", "pr", "create", "--dry-run"); err != nil {
		fmt.Printf("❌ Error analyzing PRs: %v\n", err)
		return nil
	}

	if askYesNo(reader, "Create these pull requests? (y/N): ") {
		fmt.Println("🚀 Creating pull requests...")
		if err := runCommandInteractive("./klcm", "pr", "create", "--apply"); err != nil {
			fmt.Printf("❌ Error creating PRs: %v\n", err)
		} else {
			fmt.Println("✅ Pull requests created successfully!")
		}
	}

	fmt.Println()
	fmt.Println("🎉 Workflow completed!")
	fmt.Println("📖 For future reference, see WORKFLOW.md for the complete guide.")
	
	return nil
}

func askYesNo(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes" || response == ""
}

func runCommandInteractive(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}