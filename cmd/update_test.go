package cmd_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/aureliano/caravela"
	"github.com/aureliano/caravela/provider"
	"github.com/aureliano/db-unit-extractor/cmd"
	"github.com/stretchr/testify/assert"
)

func TestNewUpdateCommandError(t *testing.T) {
	fakeExit := func(int) {
		panic("os.Exit called")
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()

	upcmd := func(c caravela.Conf) (*provider.Release, error) {
		return nil, fmt.Errorf("update error")
	}

	c := cmd.NewUpdateCommand(upcmd)
	assert.Equal(t, "update", c.Use)
	assert.Equal(t, "Updates this program", c.Short)
	assert.Equal(t, "Checks for a newer version of this program and updates it if necessary.", c.Long)

	output := new(bytes.Buffer)
	c.SetArgs([]string{"update"})
	c.SetOut(output)

	assert.PanicsWithValue(t, "os.Exit called", func() { c.Execute() }, "os.Exit was not called")
	txt := output.String()
	assert.Equal(t, txt, "Program update failed! update error\n")
}

func TestNewUpdateCommand(t *testing.T) {
	now := time.Now()
	upcmd := func(c caravela.Conf) (*provider.Release, error) {
		return &provider.Release{
			Name:        "v1.0.0",
			Description: "unit test",
			ReleasedAt:  now,
		}, nil
	}

	c := cmd.NewUpdateCommand(upcmd)
	assert.Equal(t, "update", c.Use)
	assert.Equal(t, "Updates this program", c.Short)
	assert.Equal(t, "Checks for a newer version of this program and updates it if necessary.", c.Long)

	output := new(bytes.Buffer)
	c.SetArgs([]string{"update"})
	c.SetOut(output)
	err := c.Execute()
	assert.Nil(t, err)

	txt := output.String()
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Release v1.0.0 of %s.\n\n", now.Format("02/01/2006 15:04:05")))
	sb.WriteString("unit test\n")
	sb.WriteString("\nUpdate from version v0.0.0-dev to v1.0.0 successfully completed!\n")

	assert.Equal(t, sb.String(), txt)
}
