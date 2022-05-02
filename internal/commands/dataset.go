package commands

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/internal/models"
)

func init() {
	datasetCmd.AddCommand(datasetGetCmd)
	datasetCmd.AddCommand(datasetAllCmd)
	rootCmd.AddCommand(datasetCmd)
}

var datasetCmd = &cobra.Command{
	Use:   "dataset [command]",
	Short: "Dataset commands",
}

var datasetGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get datasets by id",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s := newStore()
		e := json.NewEncoder(os.Stdout)
		for _, id := range args {
			d, err := s.GetDataset(id)
			if err != nil {
				log.Fatal(err)
			}
			e.Encode(d)
		}
	},
}

var datasetAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Get all datasets",
	Run: func(cmd *cobra.Command, args []string) {
		s := newStore()
		e := json.NewEncoder(os.Stdout)
		s.EachDataset(func(d *models.Dataset) bool {
			e.Encode(d)
			return true
		})
	},
}
