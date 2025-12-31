package cli

import (
	"fmt"
	"time"

	"github.com/ansoncodes/workshot/internal/snapshot"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(restoreCmd)
	// Add a flag to output only shell commands for eval
	restoreCmd.Flags().BoolP("commands", "c", false, "Output only shell commands for eval")
}

var restoreCmd = &cobra.Command{
	Use:   "restore [name]",
	Short: "Restore a saved development context",
	Long: `Restore displays your saved working state and prints the commands to restore it.

Restore WILL:
• Display saved working directory
• Show Git state and recent commands
• Emit shell commands to change directory

Restore WON'T (due to shell limitations):
• Change your current shell's directory
• Run commands automatically
• Restore running processes

Examples:
  workshot restore my-task            # Show context and commands
  eval $(workshot restore my-task -c) # Execute restore commands
  cd $(workshot restore my-task -c)   # Just change directory`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		
		// Check if we're in command-only mode
		commandsOnly, _ := cmd.Flags().GetBool("commands")
		
		// Setup plugin manager
		manager := initPluginManager()

		// Load snapshot
		snap, errors := snapshot.Restore(name, manager)
		if snap == nil {
			return fmt.Errorf("failed to load snapshot '%s'", name)
		}

		// COMMAND-ONLY MODE: Just output shell commands for eval
		if commandsOnly {
			fmt.Printf("cd %q\n", snap.WorkingDir)
			if snap.GitBranch != "" {
				fmt.Printf("git checkout %s\n", snap.GitBranch)
			}
			return nil
		}

		// HUMAN-READABLE MODE: Show info and commands
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		fmt.Printf("↻ Restoring workshot '%s'...\n\n", cyan(name))

		// Show restored data
		fmt.Printf("✓ Working directory: %s\n", snap.WorkingDir)

		if snap.GitBranch != "" {
			fmt.Printf("%s Git branch: %s", green("✓"), snap.GitBranch)
			if snap.GitDirty {
				fmt.Printf(" %s", yellow("(uncommitted changes)"))
			}
			fmt.Println()
		}

		// Show warnings
		if len(errors) > 0 {
			fmt.Println()
			for _, err := range errors {
				fmt.Printf("⚠ Warning: %v\n", err)
			}
		}

		// Show snapshot age
		age := time.Since(snap.CreatedAt)
		fmt.Printf("\n⏱ Saved: %s\n", formatAge(age))

		// Show recent commands if present
		if termData, ok := snap.PluginData["terminal"].(map[string]interface{}); ok {
			if commandsInterface, ok := termData["recent_commands"]; ok {
				if commandsList, ok := commandsInterface.([]interface{}); ok {
					fmt.Println("\nRecent commands from this context:")
					displayCount := len(commandsList)
					if displayCount > 10 {
						displayCount = 10
					}
					start := len(commandsList) - displayCount
					for i := start; i < len(commandsList); i++ {
						if cmd, ok := commandsList[i].(string); ok {
							fmt.Printf("   %s\n", cmd)
						}
					}
				}
			}
		}

		// CRITICAL: Show shell commands to actually restore context
		fmt.Println("\n---")
		fmt.Println("Commands to restore context:")
		fmt.Printf("    cd %q\n", snap.WorkingDir)
		if snap.GitBranch != "" {
			fmt.Printf("    git checkout %s\n", snap.GitBranch)
		}
		
		fmt.Println("\nQuick restore options:")
		fmt.Printf("  PowerShell:     cd %q\n", snap.WorkingDir)
		fmt.Printf("  Bash/Zsh:       cd %s\n", snap.WorkingDir)
		fmt.Printf("  Using eval:     eval $(workshot restore %s -c)\n", name)
		
		if len(errors) == 0 {
			fmt.Printf("\n%s Context information restored\n", green("✓"))
		} else {
			fmt.Printf("\n%s Context partially restored (see warnings above)\n", yellow("⚠"))
		}

		return nil
	},
}

func formatAge(d time.Duration) string {
	// format time since snapshot was created
	if d.Minutes() < 1 {
		return "just now"
	} else if d.Hours() < 1 {
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	} else if d.Hours() < 24 {
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}