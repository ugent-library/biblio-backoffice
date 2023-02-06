package cmd

import (
	"bytes"
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

// Test whether all fields are present

// Update existing valid record

// Test whether all fields are present

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

// Dissertation
// JournalArticle
// Miscellaneous
// book
// book_chapter
// conference
// book_editor
// issue_editor

// Test if all fields return properly

// Test validation
func (s *AddPublicationSuite) TestAddAndGetValidPublication() {
	t := s.T()

	// Add a publication
	// addPublications := AddPublicationsCmd{}
	// addCmd := addPublications.Command()
	// addCmdBuf := bytes.NewBufferString("")
	// addCmd.SetOut(addCmdBuf)

	addCmdOutFile, errJsonl := toJSONL("../../etc/fixtures/complete.publication.json")
	if errJsonl != nil {
		t.Fatal(errJsonl)
	}

	// in := strings.NewReader(addCmdOutFile)

	// addCmd.SetIn(in)
	// errAddCmd := addCmd.Execute()
	// if errAddCmd != nil {
	// 	t.Fatal(errAddCmd)
	// }

	// _, err := ioutil.ReadAll(addCmdBuf)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(addCmdOut))

	// Retrieve the publication
	getPublication := GetPublicationCmd{}
	getCmd := getPublication.Command()
	getCmdBuf := bytes.NewBufferString("")
	getCmd.SetOut(getCmdBuf)

	getCmd.SetArgs([]string{"01GMQZQJZRH678JAYAW1J5VXSY"})
	errGetCmd := getCmd.Execute()
	if errGetCmd != nil {
		t.Fatal(errGetCmd)
	}

	getCmdOut, err := ioutil.ReadAll(getCmdBuf)
	if err != nil {
		t.Fatal(err)
	}

	getCmdOutFile := string(getCmdOut)
	t.Log(getCmdOutFile)

	assert.JSONEqf(t, addCmdOutFile, getCmdOutFile, "error message %s", "formatted")

	// assert.Regexp(t, `validation failed for publication .* at line .: publication.type.required\[\/type\]`, string(out))
}

func TestAddPublicationSuite(t *testing.T) {
	suite.Run(t, new(AddPublicationSuite))
}

func toJSONL(f string) (string, error) {
	file, err := os.ReadFile(f)
	if err != nil {
		return "", err
	}
	json := string(file)

	json = strings.ReplaceAll(json, "  ", "")
	json = strings.ReplaceAll(json, "\n", "")
	json = json + "\n"

	return json, nil
}
