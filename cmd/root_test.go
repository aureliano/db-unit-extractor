package cmd_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/aureliano/db-unit-extractor/cmd"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCommandHelp(t *testing.T) {
	c := cmd.NewRootCommand()
	shortDoc := "Database extractor for unit testing."
	longDoc := fmt.Sprintf("%s\nGo to https://github.com/aureliano/db-unit-extractor/issues "+
		"in order to report a bug or make any suggestion.", shortDoc)

	assert.Equal(t, "db-unit-extractor", c.Use)
	assert.Equal(t, shortDoc, c.Short)
	assert.Equal(t, longDoc, c.Long)

	output := new(bytes.Buffer)
	c.SetOut(output)
	err := c.Execute()
	assert.Nil(t, err)

	txt := output.String()
	assert.Contains(t, txt, "Database extractor for unit testing.")
	assert.Contains(t, txt, "Usage:\n  db-unit-extractor [flags]")
	assert.Contains(t, txt, "Flags:\n")
	assert.Contains(t, txt, "  -h, --help      help for db-unit-extractor")
	assert.Contains(t, txt, "  -v, --version   Print db-unit-extractor version")
}

func TestNewRootCommandVersion(t *testing.T) {
	c := cmd.NewRootCommand()
	shortDoc := "Database extractor for unit testing."
	longDoc := fmt.Sprintf("%s\nGo to https://github.com/aureliano/db-unit-extractor/issues "+
		"in order to report a bug or make any suggestion.", shortDoc)

	assert.Equal(t, "db-unit-extractor", c.Use)
	assert.Equal(t, shortDoc, c.Short)
	assert.Equal(t, longDoc, c.Long)

	output := new(bytes.Buffer)
	c.SetArgs([]string{"-v"})
	c.SetOut(output)
	err := c.Execute()
	assert.Nil(t, err)

	txt := output.String()
	assert.Contains(t, txt, "v0.0.0-dev")
}
