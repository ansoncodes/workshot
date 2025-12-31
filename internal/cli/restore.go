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

		commandsOnly, _ := cmd.Flags().GetBool("commands")

		manager := initPluginManager()

		snap, errors := snapshot.Restore(name, manager)
		if snap == nil {
			return fmt.Errorf("failed to load snapshot '%s'", name)
		}

		// COMMAND-ONLY MODE
		if commandsOnly {
			fmt.Printf("cd %q\n", snap.WorkingDir)
			if snap.GitBranch != "" {
				fmt.Printf("git checkout %s\n", snap.GitBranch)
			}
			return nil
		}

		// Formatters (match `show`)
		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		boldCyan := color.New(color.Bold, color.FgCyan).SprintFunc()
		gray := color.New(color.FgHiBlack).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		// Header
		fmt.Printf(" %s %s\n", bold("Snapshot:"), boldCyan(name))
		fmt.Printf("   %s %s\n", bold("Created:"), gray(snap.CreatedAt.Format("2006-01-02 15:04:05")))
		fmt.Println()

		// Working Directory
		fmt.Printf(" %s\n", bold("Working Directory:"))
		fmt.Printf("   %s\n", snap.WorkingDir)
		fmt.Println()

		// Git State
		if snap.GitBranch != "" || snap.GitRemote != "" {
			fmt.Printf(" %s\n", bold("Git State:"))

			if snap.GitBranch != "" {
				fmt.Printf("   %s  %s\n", bold("Branch:"), cyan(snap.GitBranch))
			}

			if snap.GitRemote != "" {
				fmt.Printf("   %s  %s\n", bold("Remote:"), gray(snap.GitRemote))
			}

			if snap.GitDirty {
				fmt.Printf("   %s  %s\n", bold("Status:"), yellow("Modified (uncommitted changes)"))
			} else {
				fmt.Printf("   %s  Clean\n", bold("Status:"))
			}

			if gitData, ok := snap.PluginData["git"].(map[string]interface{}); ok {
				if commit, ok := gitData["commit"].(string); ok && commit != "" {
					fmt.Printf("   %s  %s\n", bold("Commit:"), gray(commit))
				}

				if stashCount, ok := gitData["stash_count"].(float64); ok && stashCount > 0 {
					fmt.Printf("   %s  %.0f\n", bold("Stashes:"), stashCount)
				}
			}

			fmt.Println()
		}

		// Recent Commands
		if termData, ok := snap.PluginData["terminal"].(map[string]interface{}); ok {
			if commandsList, ok := termData["recent_commands"].([]interface{}); ok && len(commandsList) > 0 {
				fmt.Printf(" %s\n", bold("Recent Commands:"))

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
				fmt.Println()
			}
		}

		// Commands to restore
		fmt.Printf(" %s\n", bold("Commands to restore:"))
		fmt.Printf("   cd %q\n", snap.WorkingDir)
		if snap.GitBranch != "" {
			fmt.Printf("   git checkout %s\n", snap.GitBranch)
		}

		// Warnings
		if len(errors) > 0 {
			fmt.Println()
			for _, err := range errors {
				fmt.Printf("⚠ %s %v\n", bold("Warning:"), err)
			}
		}

		return nil
	},
}

func formatAge(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	} else if d < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	} else if d < 24*time.Hour {
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	days := int(d.Hours() / 24)
	if days == 1 {
		return "1 day ago"
	}
	return fmt.Sprintf("%d days ago", days)
}
