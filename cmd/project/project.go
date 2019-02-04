package project

import (
	"io"

	"github.com/spf13/cobra"
)

type projectOptions struct {
}

var (
	projectShort   = ""
	projectLong    = ""
	projectExample = ""
)

func NewProjectCmd(out io.Writer) *cobra.Command {
	options := projectOptions{}
	cmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"projects", "proj"},
		Example: projectExample,
	}

	// create subcommands
	cmd.AddCommand(NewGetCmd(options))
	cmd.AddCommand(NewListCmd(options))
	//cmd.AddCommand(NewDeleteCmd(options))
	//cmd.AddCommand(NewSaveCmd(options))
	return cmd
}
