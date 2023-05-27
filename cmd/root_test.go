package cmd_test

import (
	"bytes"
	"testing"

	"github.com/aureliano/db-unit-extractor/cmd"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {
	c := cmd.NewRootCommand()
	assert.Equal(t, "db-unit-extractor", c.Use)
	assert.Equal(t, "Database extractor for unit testing.", c.Short)
	assert.Equal(t, "Database extractor for unit testing.", c.Long)

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
