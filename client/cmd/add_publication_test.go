package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AddPublicationSuite struct {
	suite.Suite
	cmd *cobra.Command
	buf *bytes.Buffer
}

func (s *AddPublicationSuite) SetupTest() {
	addPublications := AddPublicationsCmd{}
	s.cmd = addPublications.Command()

	viper.Set("port", "3999")
	viper.Set("insecure", true)

	s.buf = bytes.NewBufferString("")
	s.cmd.SetOut(s.buf)
}

func (s *AddPublicationSuite) TestAddPublication() {
	t := s.T()

	in := strings.NewReader("{}\n")
	s.cmd.SetIn(in)

	s.cmd.Execute()

	out, err := ioutil.ReadAll(s.buf)
	if err != nil {
		t.Fatal(err)
	}

	// t.Log(string(out))
	assert.Regexp(t, `validation failed for publication .* at line .: publication.type.required\[\/type\]`, string(out))
}

func TestAddPublicationSuite(t *testing.T) {
	suite.Run(t, new(AddPublicationSuite))
}
