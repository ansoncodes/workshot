package cli

import (
	"fmt"

	"github.com/ansoncodes/workshot/internal/snapshot"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	forceOverwrite bool
)

func init() {
	freezeCmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "Overwrite if exists")
	rootCmd.AddCommand(freezeCmd)
}

var freezeCmd = &cobra.Command{
	Use:   "freeze [name]",
	Short: "Save your current development context",
	Long: `Freeze saves your current working state including:
  • Working directory
  • Git branch, remote, and dirty state
  • Active editor
  • Recent terminal commands
  • Open files (if detectable)

The snapshot is saved to ~/.workshot/shots/ as human-readable JSON.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		yellow := color.New(color.FgYellow).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		fmt.Printf("%s Freezing workshot '%s'...\n", yellow(""), cyan(name))

		// setup plugin manager
		manager := initPluginManager()

		// save snapshot
		if err := snapshot.Freeze(name, manager); err != nil {
			return err
		}

		fmt.Printf("%s Workshot '%s' saved successfully!\n", green("✓"), cyan(name))
		fmt.Printf("   Restore it anytime with: %s\n", cyan(fmt.Sprintf("workshot restore %s", name)))

		return nil
	},
}
