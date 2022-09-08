package cmd

import "github.com/spf13/cobra"

type DatasetCmd struct {
}

func (c *DatasetCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dataset [command]",
		Short: "Dataset commands",
	}

	return cmd
}
