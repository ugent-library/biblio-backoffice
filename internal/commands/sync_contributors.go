package commands

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

func init() {
	rootCmd.AddCommand(syncContributors)
}

type change struct {
	PublicationID string `json:"publication_id"`
	Attribute     string `json:"type"`
	From          string `json:"from"`
	To            string `json:"to"`
}

var syncContributors = &cobra.Command{
	Use:   "sync-contributors",
	Short: "Synchronize contributors",
	Run: func(cmd *cobra.Command, args []string) {

		services := Services()
		logger := newLogger()

		personService := services.PersonService
		orgService := services.OrganizationService
		repo := services.Repository

		repo.EachPublication(func(pub *models.Publication) bool {

			changes := []change{}

			for _, role := range []string{"author", "editor", "supervisor"} {
				contributors := pub.Contributors(role)
				for _, c := range contributors {
					if c.ID == "" {
						continue
					}

					// TODO: how to handle records that are gone from the authority table?
					person, err := personService.GetPerson(c.ID)
					if err != nil {
						logger.Warnf(
							"unable to fetch person record %s with role %s from publication %s: %s",
							c.ID, role, pub.ID, err,
						)
						continue
					}
					if c.ORCID != person.ORCID {
						changes = append(changes, change{
							PublicationID: pub.ID,
							Attribute:     "contributor.orcid",
							From:          c.ORCID,
							To:            person.ORCID,
						})
						c.ORCID = person.ORCID
					}

					if !reflect.DeepEqual(c.UGentID, person.UGentID) {
						changes = append(changes, change{
							PublicationID: pub.ID,
							Attribute:     "contributor.ugent_id",
							From:          strings.Join(c.UGentID, ","),
							To:            strings.Join(person.UGentID, ","),
						})
						c.UGentID = append([]string{}, person.UGentID...)
					}

					oldDeps := make([]string, 0, len(c.Department))
					for _, dep := range c.Department {
						oldDeps = append(oldDeps, dep.ID)
					}
					sort.Strings(oldDeps)

					newDeps := make([]string, 0, len(person.Department))
					for _, dep := range person.Department {
						newDeps = append(newDeps, dep.ID)
					}
					sort.Strings(newDeps)

					if !reflect.DeepEqual(oldDeps, newDeps) {
						changes = append(changes, change{
							PublicationID: pub.ID,
							Attribute:     "department",
							From:          strings.Join(oldDeps, ","),
							To:            strings.Join(newDeps, ","),
						})
						for _, pd := range person.Department {
							newDep := models.ContributorDepartment{ID: pd.ID}
							org, orgErr := orgService.GetOrganization(pd.ID)
							if orgErr == nil {
								newDep.Name = org.Name
							}
							c.Department = append(c.Department, newDep)
						}
					}

				}
			}

			for _, change := range changes {
				bytes, _ := json.Marshal(change)
				fmt.Println(string(bytes))
			}

			return true
		})
	},
}
