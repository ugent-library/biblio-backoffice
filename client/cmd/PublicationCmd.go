package cmd

import "github.com/spf13/cobra"

type PublicationCmd struct {
}

func (c *PublicationCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publication [command]",
		Short: "Publication commands",
	}

	return cmd
}
