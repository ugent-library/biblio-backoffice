package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/ugent-library/biblio-backoffice/internal/models"
)

type RelateSuite struct {
	suite.Suite
}

func (s *RelateSuite) SetupSuite() {
	viper.Set("port", "3999")
	viper.Set("insecure", true)

	t := s.T()

	file, err := os.ReadFile("../../etc/fixtures/complete.dataset.json")
	if err != nil {
		t.Fatal(err)
	}

	addCmdOutFile, err := toJSONL(file)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = addDataset(addCmdOutFile)
	if err != nil {
		t.Fatal(err)
	}

	file, err = os.ReadFile("../../etc/fixtures/complete.book_chapter.json")
	if err != nil {
		t.Fatal(err)
	}

	rec, _ := addKey(string(file), "id", "00000000000000000000000061")

	jsonl, err := toJSONL([]byte(rec))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = addPublication(jsonl)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *RelateSuite) TestRelateValid() {
	t := s.T()

	stdOut, _, err := relate("00000000000000000000000061", "00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "related: publication[id: 00000000000000000000000061] -> dataset[id: 00000000000000000000000001]", stdOut)

	// Retrieve the publication
	stdOut, _, err = getPublication("00000000000000000000000061")
	if err != nil {
		t.Fatal(err)
	}

	p := &models.Publication{}
	err = json.Unmarshal([]byte(stdOut), p)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, p.RelatedDataset, 2, "number of related datasets isn't 2, got:", len(p.RelatedDataset))
	assert.Equal(t, "00000000000000000000000001", p.RelatedDataset[1].ID)

	// Retrieve the dataset
	stdOut, _, err = getDataset("00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}

	d := &models.Dataset{}
	err = json.Unmarshal([]byte(stdOut), d)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, d.RelatedPublication, 2, "number of related publications isn't 2, got:", len(d.RelatedPublication))
	assert.Equal(t, "00000000000000000000000061", d.RelatedPublication[1].ID)
}

func (s *RelateSuite) TestRelateNonExistingPublication() {
	t := s.T()

	stdOut, _, err := relate("notexists", "00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "could not find publication with id notexists", stdOut)

}

func (s *RelateSuite) TestRelateNonExistingDataset() {
	t := s.T()

	stdOut, _, err := relate("00000000000000000000000061", "notexists")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "could not find dataset with id notexists", stdOut)
}

func (s *RelateSuite) TearDownSuite() {
	t := s.T()

	_, _, err := purgePublication("00000000000000000000000061")
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = purgeDataset("00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRelateSuite(t *testing.T) {
	suite.Run(t, new(RelateSuite))
}

func relate(pid string, did string) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"relate-dataset", pid, did})

	errGet := rootCmd.Execute()
	if errGet != nil {
		return "", "", errGet
	}

	resOut, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	resErr, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	return string(resOut), string(resErr), nil
}
