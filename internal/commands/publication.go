package commands

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/internal/models"
)

func init() {
	publicationCmd.AddCommand(publicationGetCmd)
	publicationCmd.AddCommand(publicationAllCmd)
	rootCmd.AddCommand(publicationCmd)
}

var publicationCmd = &cobra.Command{
	Use:   "publication [command]",
	Short: "Publication commands",
}

var publicationGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get publications by id",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s := newStorageService()
		e := json.NewEncoder(os.Stdout)
		for _, id := range args {
			d, err := s.GetPublication(id)
			if err != nil {
				log.Fatal(err)
			}
			e.Encode(d)
		}
	},
}

var publicationAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Get all publications",
	Run: func(cmd *cobra.Command, args []string) {
		s := newStorageService()
		e := json.NewEncoder(os.Stdout)
		s.EachPublication(func(d *models.Publication) bool {
			e.Encode(d)
			return true
		})
	},
}
