package installation

import (
	"github.com/spf13/cobra"

	"github.com/satisfactorymodding/ficsit-cli/cli"
)

func init() {
	Cmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all installations",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		global, err := cli.InitCLI(false)
		if err != nil {
			return err
		}

		for _, install := range global.Installations.Installations {
			println(install.Path, "-", install.Profile)
		}

		return nil
	},
}
