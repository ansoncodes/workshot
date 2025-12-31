package cli

import (
	"encoding/json"
	"fmt"

	"github.com/ansoncodes/workshot/internal/storage"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
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

		cyan := color.New(color.FgCyan).SprintFunc()

		// print snapshot as json
		data, err := json.MarshalIndent(snap, "", "  ")
		if err != nil {
			return err
		}

		fmt.Printf("Workshot: %s\n\n", cyan(name))
		fmt.Println(string(data))

		return nil
	},
}
