package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UpdateDatasetSuite struct {
	suite.Suite
}

func (s *UpdateDatasetSuite) SetupSuite() {
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
}

// Test empty input
func (s *UpdateDatasetSuite) TestUpdateEmptyInput() {
	t := s.T()

	jsonl := ""

	_, _, err := updateDataset(jsonl)

	assert.Equal(t, "could not read from stdin: EOF", err.Error())
}

// Test invalid JSON input
func (s *UpdateDatasetSuite) TestUpdateInvalidJSON() {
	t := s.T()

	jsonl := "invalid\n"

	stdOut, _, err := updateDataset(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "could not read json input: invalid character 'i' looking for beginning of value", stdOut)
}

func (s *UpdateDatasetSuite) TestUpdateEmptyJSONLInput() {
	t := s.T()
	jsonl := "{}\n"

	stdOut, _, err := updateDataset(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `failed to validate dataset : dataset.id.required\[\/id\], dataset.status.required\[\/status\]`, string(stdOut))
}

// Test if a dataset can be succesfully updated
func (s *UpdateDatasetSuite) TestUpdateValid() {
	t := s.T()

	// Retrieve the dataset as user A
	pubA, _, err := getDataset("00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}

	updPubA, err := changeTitle(pubA, "update title")
	if err != nil {
		t.Fatal(err)
	}

	jsonl, err := toJSONL([]byte(updPubA))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = updateDataset(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	pubB, _, err := getDataset("00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	updPubA, err = removeKey(updPubA, "snapshot_id", "date_from", "date_updated", "user_id", "last_user_id")
	if err != nil {
		t.Fatal(err)
	}

	pubB, err = removeKey(pubB, "snapshot_id", "date_from", "date_updated", "user_id", "last_user_id")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, updPubA, pubB, "the added and retrieved records aren't equal.")
}

// Test an update conflict
func (s *UpdateDatasetSuite) TestUpdateConflict() {
	t := s.T()

	// Retrieve the dataset as user A
	getCmdStdOut, _, err := getDataset("00000000000000000000000001")

	if err != nil {
		t.Fatal(err)
	}

	// user B pushes a change to the dataset

	updated, err := changeTitle(getCmdStdOut, "new title")
	if err != nil {
		t.Fatal(err)
	}

	jsonl, err := toJSONL([]byte(updated))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = updateDataset(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	// // User A tries to change the dataset with the old snapshotID

	updated, err = changeTitle(getCmdStdOut, "alternate title")
	if err != nil {
		t.Fatal(err)
	}

	jsonl, err = toJSONL([]byte(updated))
	if err != nil {
		t.Fatal(err)
	}

	updCmdOut, _, err := updateDataset(jsonl)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `failed to update dataset: conflict detected for dataset\[snapshot_id: .*, id: .*\] : version conflict`, updCmdOut)
}

func (s *UpdateDatasetSuite) TearDownSuite() {
	t := s.T()

	_, _, err := purgeDataset("00000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateDatasetSuite(t *testing.T) {
	suite.Run(t, new(UpdateDatasetSuite))
}

func updateDataset(jsonl string) (string, string, error) {
	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)

	rootCmd.SetArgs([]string{"dataset", "update"})

	in := strings.NewReader(jsonl)
	rootCmd.SetIn(in)

	err := rootCmd.Execute()
	if err != nil {
		return "", "", err
	}

	updCmdOut, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	updCmdErr, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", "", err
	}

	return string(updCmdOut), string(updCmdErr), nil
}
