package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ansoncodes/workshot/internal/storage"
	"github.com/ansoncodes/workshot/pkg/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	showCmd.Flags().BoolP("json", "j", false, "Output raw JSON")
	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:     "show [name]",
	Aliases: []string{"info"},
	Short:   "Show detailed information about a workshot",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		store, err := storage.New()
		if err != nil {
			return err
		}

		snap, err := store.Load(name)
		if err != nil {
			return err
		}

		// JSON MODE
		jsonOut, _ := cmd.Flags().GetBool("json")
		if jsonOut {
			data, err := json.MarshalIndent(snap, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		printSnapshot(name, snap)
		return nil
	},
}

func printSnapshot(name string, snap *types.Snapshot) {
	// Formatters
	bold := color.New(color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	boldCyan := color.New(color.Bold, color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	// Header
	fmt.Printf(" %s %s\n", bold("Snapshot:"), boldCyan(name))
	fmt.Printf("   %s %s\n", bold("Created:"), gray(snap.CreatedAt.Format("2006-01-02 15:04:05")))
	fmt.Println()

	// Working Directory
	fmt.Printf(" %s\n", bold("Working Directory:"))
	fmt.Printf("   %s\n", green(snap.WorkingDir))
	fmt.Println()

	// Git State
	if snap.GitBranch != "" {
		fmt.Printf(" %s\n", bold("Git State:"))
		fmt.Printf("   %s  %s\n", bold("Branch:"), cyan(snap.GitBranch))

		status := "Clean"
		if snap.GitDirty {
			status = "Dirty (uncommitted changes)"
		}
		fmt.Printf("   %s  %s\n", bold("Status:"), status)

		if snap.GitRemote != "" {
			fmt.Printf("   %s  %s\n", bold("Remote:"), gray(snap.GitRemote))
		}

		if gitData, ok := snap.PluginData["git"].(map[string]interface{}); ok {
			if commit, ok := gitData["commit"].(string); ok && commit != "" {
				fmt.Printf("   %s  %s\n", bold("Commit:"), gray(commit))
			}
			if stash, ok := gitData["stash_count"].(float64); ok && stash > 0 {
				fmt.Printf("   %s  %.0f\n", bold("Stashes:"), stash)
			}
		}
		fmt.Println()
	}

	// Editor
	if editorData, ok := snap.PluginData["editor"].(map[string]interface{}); ok {
		if detected, ok := editorData["detected"].(string); ok && detected != "" {
			fmt.Printf(" %s\n", bold("Editor:"))
			fmt.Printf("   %s\n\n", cyan(detected))
		}
	}

	// Recent Commands
	if terminalData, ok := snap.PluginData["terminal"].(map[string]interface{}); ok {
		if commands, ok := terminalData["recent_commands"].([]interface{}); ok && len(commands) > 0 {
			fmt.Printf(" %s\n", bold("Recent Commands:"))

			max := 10
			start := 0
			if len(commands) > max {
				start = len(commands) - max
			}

			for _, c := range commands[start:] {
				if cmd, ok := c.(string); ok {
					colorFn := white
					if strings.HasPrefix(cmd, "git ") {
						colorFn = cyan
					}
					fmt.Printf("   %s\n", colorFn(cmd))
				}
			}
			fmt.Println()
		}
	}

	// Metadata
	fmt.Printf(" %s\n", bold("Metadata:"))
	fmt.Printf("   %s %d\n", bold("Schema Version:"), snap.SchemaVersion)
	fmt.Printf("   %s %s\n", bold("Age:"), formatDuration(time.Since(snap.CreatedAt)))
	fmt.Printf("   %s %d active\n", bold("Plugins:"), len(snap.PluginData))
}

func formatDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	}
}
