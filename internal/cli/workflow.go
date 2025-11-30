package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// workflowCmd represents the workflow command
var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Interactive guide for the complete keyboard configuration workflow",
	Long: `Step-by-step interactive guide for making, validating, and contributing keyboard configuration changes.

This command walks you through the complete process:
1. Pull latest configurations from remote repositories
2. Make and validate your changes
3. Sync changes between keyboards (optional)
4. Create pull requests to contribute back

This is perfect for new users or when you want a guided experience.`,
	Example: `  # Start the interactive workflow
  klcm workflow
  
  # Start workflow with verbose output
  klcm workflow --verbose`,
	RunE: runWorkflow,
}

func runWorkflow(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 KLCM Interactive Keyboard Configuration Workflow")
	fmt.Println("===================================================")
	fmt.Println()
	fmt.Println("This guide will walk you through the complete process of making")
	fmt.Println("keyboard configuration changes and contributing them back.")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// Step 1: Update configurations
	fmt.Println("📋 Step 1: Update to Latest Configurations")
	fmt.Println("-----------------------------------------------")
	fmt.Println("Before making changes, let's ensure you have the latest configurations.")
	fmt.Println()
	
	if askYesNoWorkflow(reader, "Check for and preview remote updates? (Y/n): ") {
		fmt.Println()
		fmt.Println("🔍 Checking for remote updates...")
		if err := runCLICommand("pull", "--preview"); err != nil {
			fmt.Printf("❌ Error checking updates: %v\n", err)
		} else {
			fmt.Println()
			if askYesNoWorkflow(reader, "Apply these updates? (Y/n): ") {
				if err := runCLICommand("pull"); err != nil {
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
	fmt.Println("----------------------------")
	fmt.Println("Now edit your keyboard configuration files:")
	fmt.Println()
	fmt.Println("  📁 Advantage360: configs/zmk_adv360/adv360.keymap")
	fmt.Println("  📁 Glove80:      configs/zmk_glove80/glove80.keymap")
	fmt.Println("  📁 Pillz Mod:    configs/zmk_adv_mod/pillzmod_pro.keymap")
	fmt.Println()
	fmt.Println("💡 Tips:")
	fmt.Println("   - Make small, focused changes")
	fmt.Println("   - Test one change at a time")
	fmt.Println("   - Keep notes of what you're changing and why")
	fmt.Println()
	fmt.Println("Press Enter when you've finished making your changes...")
	reader.ReadLine()

	// Step 3: Validate changes
	fmt.Println("📋 Step 3: Validate Your Changes")
	fmt.Println("---------------------------------")
	fmt.Println("Let's make sure your changes are syntactically correct.")
	fmt.Println()
	
	fmt.Println("🔍 Running validation...")
	if err := runCLICommand("validate"); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
		fmt.Println()
		fmt.Println("Please fix the errors shown above and run 'klcm validate' again.")
		fmt.Println("Then restart this workflow with 'klcm workflow'")
		return nil
	}
	fmt.Println("✅ All configurations are valid!")
	fmt.Println()

	// Step 4: Sync between keyboards (optional)
	fmt.Println("📋 Step 4: Sync Between Keyboards (Optional)")
	fmt.Println("--------------------------------------------")
	fmt.Println("If you want the same changes on both keyboards, we can sync them.")
	fmt.Println()
	
	if askYesNoWorkflow(reader, "Sync changes between keyboards? (y/N): ") {
		// Ask which direction to sync
		fmt.Println()
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

		fmt.Printf("\n🔍 Previewing sync from %s to %s...\n", source, target)
		if err := runCLICommand("sync", source, target, "--preview"); err != nil {
			fmt.Printf("❌ Error previewing sync: %v\n", err)
		} else {
			fmt.Println()
			if askYesNoWorkflow(reader, "Apply this sync? (y/N): ") {
				if err := runCLICommand("sync", source, target); err != nil {
					fmt.Printf("❌ Error applying sync: %v\n", err)
				} else {
					fmt.Println("✅ Sync applied successfully")
					// Validate after sync
					fmt.Println()
					fmt.Println("🔍 Re-validating after sync...")
					if err := runCLICommand("validate"); err != nil {
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
	fmt.Println("--------------------------------")
	fmt.Println("Finally, let's create pull requests to contribute your changes back.")
	fmt.Println()
	
	fmt.Println("🔍 Analyzing what PRs would be created...")
	if err := runCLICommand("pr", "create", "--dry-run"); err != nil {
		fmt.Printf("❌ Error analyzing PRs: %v\n", err)
		return nil
	}

	fmt.Println()
	if askYesNoWorkflow(reader, "Create these pull requests? (y/N): ") {
		fmt.Println()
		fmt.Println("🚀 Creating pull requests...")
		fmt.Println("⚠️  Note: This requires GitHub CLI (gh) and authentication")
		fmt.Println()
		if err := runCLICommand("pr", "create", "--apply"); err != nil {
			fmt.Printf("❌ Error creating PRs: %v\n", err)
			fmt.Println()
			fmt.Println("💡 You can also create PRs manually:")
			fmt.Println("   1. Commit your changes: git add . && git commit -m 'Update keyboard layout'")
			fmt.Println("   2. Push to a branch: git push origin feature/my-layout")
			fmt.Println("   3. Create PR on GitHub website")
		} else {
			fmt.Println("✅ Pull requests created successfully!")
		}
	}

	fmt.Println()
	fmt.Println("🎉 Workflow completed successfully!")
	fmt.Println()
	fmt.Println("📚 Resources:")
	fmt.Println("   📖 Complete guide: WORKFLOW.md")
	fmt.Println("   🆘 Command help: klcm [command] --help")
	fmt.Println("   🔍 Check changes: git status && git diff")
	fmt.Println()
	fmt.Println("Thank you for contributing to the keyboard community! 🙏")
	
	return nil
}

func askYesNoWorkflow(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes" || response == ""
}

func runCLICommand(args ...string) error {
	cmd := exec.Command("./klcm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(workflowCmd)
}