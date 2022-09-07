package cmd

import "github.com/spf13/cobra"

type FileCmd struct {
}

func (c *FileCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file [command]",
		Short: "File commands",
	}

	return cmd
}
