package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AddPublicationSuite struct {
	suite.Suite
	addCmd *cobra.Command
	getCmd *cobra.Command
	buf    *bytes.Buffer
}

func (s *AddPublicationSuite) SetupTest() {
	addPublications := AddPublicationsCmd{}
	s.addCmd = addPublications.Command()

	getPublication := GetPublicationCmd{}
	s.getCmd = getPublication.Command()

	viper.Set("port", "3999")
	viper.Set("insecure", true)

	s.buf = bytes.NewBufferString("")
	s.addCmd.SetOut(s.buf)
	s.getCmd.SetOut(s.buf)
}

// Test empty input
// Test non-JSONL input
// Test empty JSONL object

// Create new valid record

// Update existing valid record

// Test if all fields return properly
func (s *AddPublicationSuite) TestAddAndGetCompletePublications() {
	t := s.T()

	file, errFile := os.ReadFile("../../etc/fixtures/complete.book_chapter.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL := toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd := addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut := getPublication("00000000000000000000000001")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err := removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err := removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// assert.Regexp(t, `validation failed for publication .* at line .: publication.type.required\[\/type\]`, string(out))

	file, errFile = os.ReadFile("../../etc/fixtures/complete.book_editor.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000002")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// assert.Regexp(t, `validation failed for publication .* at line .: publication.type.required\[\/type\]`, string(out))

	file, errFile = os.ReadFile("../../etc/fixtures/complete.book.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000003")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// assert.Regexp(t, `validation failed for publication .* at line .: publication.type.required\[\/type\]`, string(out))

	file, errFile = os.ReadFile("../../etc/fixtures/complete.conference.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000004")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// assert.Regexp(t, `validation failed for publication .* at line .: publication.type.required\[\/type\]`, string(out))

	file, errFile = os.ReadFile("../../etc/fixtures/complete.dissertation.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000005")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// assert.Regexp(t, `validation failed for publication .* at line .: publication.type.required\[\/type\]`, string(out))

	file, errFile = os.ReadFile("../../etc/fixtures/complete.issue_editor.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000006")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// assert.Regexp(t, `validation failed for publication .* at line .: publication.type.required\[\/type\]`, string(out))

	file, errFile = os.ReadFile("../../etc/fixtures/complete.journal_article.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000007")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// assert.Regexp(t, `validation failed for publication .* at line .: publication.type.required\[\/type\]`, string(out))

	file, errFile = os.ReadFile("../../etc/fixtures/complete.miscellaneous.json")
	if errFile != nil {
		t.Fatal(errFile)
	}

	// Add a publication
	addCmdOutFile, errJSONL = toJSONL(file)
	if errJSONL != nil {
		t.Fatal(errJSONL)
	}

	errAdd = addPublication(addCmdOutFile)
	if errAdd != nil {
		t.Fatal(errAdd)
	}

	// Retrieve the publication
	getCmdOutFile, errGetCmdOut = getPublication("00000000000000000000000008")
	if errGetCmdOut != nil {
		t.Fatal(errGetCmdOut)
	}

	// Remove dynamic fields distorting the JSONEqf assertion
	ina, err = removeKey(addCmdOutFile, "snapshot_id", "date_from", "date_updated", "user")
	if err != nil {
		t.Fatal(err)
	}

	ing, err = removeKey(getCmdOutFile, "snapshot_id", "date_from", "date_updated")
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, ina, ing, "the added and retrieved records aren't equal.")

	// assert.Regexp(t, `validation failed for publication .* at line .: publication.type.required\[\/type\]`, string(out))
}

// Test validation

// Validation: type field is required
// Validation: type is invalid
// Validation: classification is required
// Validation: classification is invalid
// Validation: status is required
// Validation: status is invalid
// Validation: title field is required
// Validation: year field is required
// Validation: year field is invalid
// Validation: language code is invalid
// Validation: abstract => ID is required
// Validation: abstract => lang is required
// Validation: abstract => lang is invalid
// Validation: abstraxt => text is required
// Validation: laysummary => ID is required
// Validation: laysummary => lang is required
// Validation: laysummary => lang is invalid
// Validation: laysummary => text is required
// Validation: if status = public => at least 1 author
// Validation: if status = public, !legacy, usesAuthor(), !extern => at least 1 ugent author
// Validation: if status = public, usesEditor, !usesAuthor => at least 1 editor
// Validation: if status = public, usesEditor, !usesAuthor, !extern => at least 1 ugent editor
// Validation: if status = public, !legacy, usesSupervisor => at least 1 supervisor
// Validation: author => ID is required
// Validation: author => first name is required
// Validation: author => last name is required
// Validation: author => credit role is required
// Validation: author => credit role is invalid
// Validation: editor => ID is required
// Validation: editor => first name is required
// Validation: editor => last name is required
// Validation: editor => credit role is required
// Validation: editor => credit role is invalid
// Validation: supervisor => ID is required
// Validation: supervisor => first name is required
// Validation: supervisor => last name is required
// Validation: supervisor => credit role is required
// Validation: supervisor => credit role is invalid
// Validation: project => ID is required
// Validation: department => ID is required
// Validation: relatedDataset => ID is required
// Validation: file ........
// Validation: link => invalid relation

// Specifics

// Dissertation xx
// JournalArticle xx
// Miscellaneous xx
// book xx
// book_chapter xxx
// conference xx
// book_editor xxx
// issue_editor x

func TestAddPublicationSuite(t *testing.T) {
	suite.Run(t, new(AddPublicationSuite))
}

func addPublication(jsonl string) error {
	addPublications := AddPublicationsCmd{}
	addCmd := addPublications.Command()
	addCmdBuf := bytes.NewBufferString("")
	addCmd.SetOut(addCmdBuf)

	in := strings.NewReader(jsonl)

	addCmd.SetIn(in)
	errAdd := addCmd.Execute()
	if errAdd != nil {
		return errAdd
	}

	// addCmdOut, err := ioutil.ReadAll(addCmdBuf)
	// if err != nil {
	// 	return err
	// }
	// t.Log(string(addCmdOut))

	return nil
}

func getPublication(id string) (string, error) {
	getPublication := GetPublicationCmd{}
	getCmd := getPublication.Command()
	getCmdBuf := bytes.NewBufferString("")
	getCmd.SetOut(getCmdBuf)

	getCmd.SetArgs([]string{id})
	errGetCmd := getCmd.Execute()
	if errGetCmd != nil {
		return "", errGetCmd
	}

	getCmdOut, err := ioutil.ReadAll(getCmdBuf)
	if err != nil {
		return "", err
	}

	return string(getCmdOut), nil
}

func toJSONL(file []byte) (string, error) {
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, file); err != nil {
		return "", err
	}

	output, _ := ioutil.ReadAll(buffer)
	json := string(output)
	json = json + "\n"

	return json, nil
}

func removeKey(input string, keys ...string) (string, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		return "", err
	}

	output := input

	for _, key := range keys {
		if _, exists := m[key]; exists {
			delete(m, key)
			outputData, err := json.Marshal(m)
			if err != nil {
				return "", err
			}
			output = string(outputData)
		}
	}

	return output, nil
}
