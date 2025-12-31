package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ansoncodes/workshot/internal/storage"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	deleteForce bool
)

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Skip confirmation")
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:     "delete [name]",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a saved workshot",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		store, err := storage.New()
		if err != nil {
			return err
		}

		// check if snapshot exists
		if !store.Exists(name) {
			return fmt.Errorf("workshot '%s' not found", name)
		}

		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		// ask before delete unless force is set
		if !deleteForce {
			fmt.Printf("%s Are you sure you want to delete workshot '%s'? (y/N): ",
				yellow("⚠"), cyan(name))

			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		// delete snapshot
		if err := store.Delete(name); err != nil {
			return err
		}

		fmt.Printf("%s Deleted workshot '%s'\n", green("✓"), cyan(name))
		return nil
	},
}
