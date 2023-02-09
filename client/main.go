package main

import (
	"log"

	"github.com/ugent-library/biblio-backoffice/client/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
