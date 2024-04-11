package cli

import (
	"github.com/spf13/cobra"
)

type Base struct {
	Command *cobra.Command
}

func New() *Base {
	base := &Base{
		Command: &cobra.Command{
			Use:   "main",
			Short: "all commands",
		},
	}

	var verbose int
	base.Command.PersistentFlags().IntVarP(&verbose, "verbose", "v", 0, "Enable verbose logging")
	return base
}

func (c *Base) AddCommand(cmds ...*cobra.Command) {
	c.Command.AddCommand(cmds...)
}

func (c *Base) Execute() error {
	return c.Command.Execute()
}
