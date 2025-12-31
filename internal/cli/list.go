package cli

import (
	"fmt"
	"time"

	"github.com/ansoncodes/workshot/internal/storage"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all saved workshots",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.New()
		if err != nil {
			return err
		}

		// load snapshot list
		metadataList, err := store.List()
		if err != nil {
			return err
		}

		cyan := color.New(color.FgCyan).SprintFunc()
		gray := color.New(color.FgHiBlack).SprintFunc()

		// no snapshots found
		if len(metadataList) == 0 {
			fmt.Println("No saved workshots found.")
			fmt.Println("\nCreate your first workshot with:")
			fmt.Printf("  %s\n", cyan("workshot freeze my-work"))
			return nil
		}

		fmt.Printf("Found %d saved workshot(s):\n\n", len(metadataList))

		for _, meta := range metadataList {
			fmt.Printf("  %s\n", cyan(meta.Name))

			age := time.Since(meta.CreatedAt)
			fmt.Printf("     %s • %s", gray(formatAge(age)), meta.WorkingDir)
			if meta.GitBranch != "" {
				fmt.Printf(" • %s", meta.GitBranch)
			}
			fmt.Println()
		}

		fmt.Println("\nRestore any workshot with:")
		fmt.Printf("  %s\n", cyan("workshot restore <name>"))

		return nil
	},
}
