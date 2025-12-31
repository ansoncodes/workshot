package cli

import (
	"github.com/spf13/cobra"

	"github.com/ansoncodes/workshot/internal/version"
)

var rootCmd = &cobra.Command{
	Use:   "workshot",
	Short: "save and restore your development context",
	Long: `workshot helps you save your development state
and restore it later with one command

this helps when switching tasks or getting interrupted

examples
  workshot freeze my-work        save current context
  workshot restore my-work       restore saved context
  workshot list                  list all saved contexts
  workshot show my-work          view context details`,
	Version: version.Get(),
}

// run the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// show only the version number
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}
