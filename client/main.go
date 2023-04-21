package main

import (
	"os"

	"github.com/ugent-library/biblio-backoffice/client/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
