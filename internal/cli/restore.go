package cli

import (
	"fmt"
	"time"

	"github.com/YOUR_USERNAME/workshot/internal/snapshot"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(restoreCmd)
}

var restoreCmd = &cobra.Command{
	Use:   "restore [name]",
	Short: "Restore a saved development context",
	Long: `Restore brings back your saved working state including:
  • Working directory (cd to it)
  • Git branch (checkout)
  • Open files (in your editor)
  • Display recent commands for reference`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		fmt.Printf("↻ Restoring workshot '%s'...\n\n", cyan(name))

		// setup plugin manager
		manager := initPluginManager()

		// restore snapshot
		snap, errors := snapshot.Restore(name, manager)
		if snap == nil {
			return fmt.Errorf("failed to load snapshot")
		}

		// show restored data
		fmt.Printf("✓ Working directory: %s\n", snap.WorkingDir)

		if snap.GitBranch != "" {
			fmt.Printf("%s Git branch: %s", green("✓"), snap.GitBranch)
			if snap.GitDirty {
				fmt.Printf(" %s", yellow("(uncommitted changes)"))
			}
			fmt.Println()
		}

		// show warnings
		if len(errors) > 0 {
			fmt.Println()
			for _, err := range errors {
				fmt.Printf("⚠ Warning: %v\n", err)
			}
		}

		// show snapshot age
		age := time.Since(snap.CreatedAt)
		fmt.Printf("\n⏱ Saved: %s\n", formatAge(age))

		// show recent commands if present
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

		if len(errors) == 0 {
			fmt.Printf("\n%s Context restored\n", green("✓"))
		} else {
			fmt.Printf("\n%s Context partially restored see warnings above\n", yellow("⚠"))
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
