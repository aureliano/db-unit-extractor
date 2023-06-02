package cmd_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/aureliano/db-unit-extractor/cmd"
	"github.com/aureliano/db-unit-extractor/extractor"
	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/writer"
	"github.com/stretchr/testify/assert"
)

func TestNewExtractCommandSchemaIsRequired(t *testing.T) {
	c := cmd.NewExtractCommand()

	output := new(bytes.Buffer)
	c.SetArgs([]string{"extract", "-n", "postgres://usr:pwd@127.0.0.1:5432/test"})
	c.SetErr(output)

	err := c.Execute()
	assert.NotNil(t, err)

	txt := output.String()
	assert.Equal(t, txt, "Error: required flag(s) \"schema\" not set\n")
}

func TestNewExtractCommandDSNIsRequired(t *testing.T) {
	c := cmd.NewExtractCommand()

	output := new(bytes.Buffer)
	c.SetArgs([]string{"extract", "-s", "schema.yml"})
	c.SetErr(output)

	err := c.Execute()
	assert.NotNil(t, err)

	txt := output.String()
	assert.Equal(t, txt, "Error: required flag(s) \"data-source-name\" not set\n")
}

func TestNewExtractCommandInvalidDSN(t *testing.T) {
	patches := gomonkey.ApplyFunc(os.Exit, func(int) {
		panic("os.Exit called")
	})
	defer patches.Reset()

	c := cmd.NewExtractCommand()
	output := new(bytes.Buffer)
	c.SetArgs([]string{"extract", "-s", "schema.yml", "-n", "driver://invalid-dsn"})
	c.SetOut(output)

	assert.PanicsWithValue(t, "os.Exit called", func() {
		err := c.Execute()
		assert.Nil(t, err)
	}, "os.Exit was not called")
	txt := output.String()
	assert.Equal(t, txt, "Extract error: invalid DSN 'driver://invalid-dsn'\n")
}

func TestNewExtractCommandSchemaFileDoesNotExist(t *testing.T) {
	patches := gomonkey.ApplyFunc(os.Exit, func(int) {
		panic("os.Exit called")
	})
	defer patches.Reset()

	c := cmd.NewExtractCommand()
	output := new(bytes.Buffer)
	c.SetArgs([]string{"extract", "-s", "schema.yml", "-n", "postgres://usr:pwd@127.0.0.1:5432/test"})
	c.SetOut(output)

	assert.PanicsWithValue(t, "os.Exit called", func() {
		err := c.Execute()
		assert.Nil(t, err)
	}, "os.Exit was not called")
	txt := output.String()
	assert.Equal(t, txt, "Extract error: file not found 'schema.yml'\n")
}

func TestNewExtractCommandSchemaFileIsDirectory(t *testing.T) {
	patches := gomonkey.ApplyFunc(os.Exit, func(int) {
		panic("os.Exit called")
	})
	defer patches.Reset()

	c := cmd.NewExtractCommand()
	output := new(bytes.Buffer)
	c.SetArgs([]string{"extract", "-s", os.TempDir(), "-n", "postgres://usr:pwd@127.0.0.1:5432/test"})
	c.SetOut(output)

	assert.PanicsWithValue(t, "os.Exit called", func() {
		err := c.Execute()
		assert.Nil(t, err)
	}, "os.Exit was not called")
	txt := output.String()
	assert.Equal(t, txt, fmt.Sprintf("Extract error: %s is a directory\n", os.TempDir()))
}

func TestNewExtractCommandInvalidOutputType(t *testing.T) {
	patches := gomonkey.ApplyFunc(os.Exit, func(int) {
		panic("os.Exit called")
	})
	defer patches.Reset()

	c := cmd.NewExtractCommand()
	output := new(bytes.Buffer)
	c.SetArgs([]string{
		"extract", "-s", "../test/unit/schema_test.yml", "-n", "postgres://usr:pwd@127.0.0.1:5432/test", "-t", "xls"},
	)
	c.SetOut(output)

	assert.PanicsWithValue(t, "os.Exit called", func() {
		err := c.Execute()
		assert.Nil(t, err)
	}, "os.Exit was not called")
	txt := output.String()
	assert.Equal(t, txt, "Extract error: unsupported output type 'xls'\n")
}

func TestNewExtractCommandOutputDirIsNotADirectory(t *testing.T) {
	patches := gomonkey.ApplyFunc(os.Exit, func(int) {
		panic("os.Exit called")
	})
	defer patches.Reset()

	c := cmd.NewExtractCommand()
	output := new(bytes.Buffer)
	c.SetArgs([]string{
		"extract", "-s", "../test/unit/schema_test.yml", "-n", "postgres://usr:pwd@127.0.0.1:5432/test",
		"-t", "console", "-d", "../test/unit/schema_test.yml",
	})
	c.SetOut(output)

	assert.PanicsWithValue(t, "os.Exit called", func() {
		err := c.Execute()
		assert.Nil(t, err)
	}, "os.Exit was not called")
	txt := output.String()
	assert.Equal(t, txt, "Extract error: ../test/unit/schema_test.yml is not a directory\n")
}

func TestNewExtractCommandInvalidReference(t *testing.T) {
	patches := gomonkey.ApplyFunc(os.Exit, func(int) {
		panic("os.Exit called")
	})
	defer patches.Reset()

	c := cmd.NewExtractCommand()
	output := new(bytes.Buffer)
	c.SetArgs([]string{
		"extract", "-s", "../test/unit/schema_test.yml", "-n", "postgres://usr:pwd@127.0.0.1:5432/test",
		"-t", "console", "-d", os.TempDir(), "-r", "test=",
	})
	c.SetOut(output)

	assert.PanicsWithValue(t, "os.Exit called", func() {
		err := c.Execute()
		assert.Nil(t, err)
	}, "os.Exit was not called")
	txt := output.String()
	assert.Equal(t, txt, "Extract error: invalid reference 'test='\n")
}

func TestNewExtractCommandExtractError(t *testing.T) {
	patches := gomonkey.ApplyFunc(os.Exit, func(int) {
		panic("os.Exit called")
	}).ApplyFunc(extractor.Extract, func(extractor.Conf, reader.DBReader, []writer.FileWriter) error {
		return fmt.Errorf("extract error")
	})
	defer patches.Reset()

	c := cmd.NewExtractCommand()
	output := new(bytes.Buffer)
	c.SetArgs([]string{
		"extract", "-s", "../test/unit/schema_test.yml", "-n", "postgres://usr:pwd@127.0.0.1:5432/test",
		"-t", "console", "-d", os.TempDir(), "-r", "test=123",
	})
	c.SetOut(output)

	assert.PanicsWithValue(t, "os.Exit called", func() {
		err := c.Execute()
		assert.Nil(t, err)
	}, "os.Exit was not called")
	txt := output.String()
	assert.Equal(t, txt, "Extract error: extract error\n")
}

func TestNewExtractCommand(t *testing.T) {
	patches := gomonkey.ApplyFunc(extractor.Extract, func(extractor.Conf, reader.DBReader, []writer.FileWriter) error {
		return nil
	})
	defer patches.Reset()

	c := cmd.NewExtractCommand()
	assert.Equal(t, "extract", c.Use)
	assert.Equal(t, "Extract data-set from database", c.Short)
	assert.Equal(t, "Extract data-set from a database to any supported file.", c.Long)

	exampleLines := strings.Split(c.Example, "\n")
	assert.Equal(t, "# Extract data-set from PostgreSQL and write to the console.", strings.TrimLeft(exampleLines[0], " "))
	assert.Equal(t, "db-unit-extractor extract -s /path/to/schema.yml -n postgres://usr:pwd@127.0.0.1:5432/test",
		strings.TrimLeft(exampleLines[1], " "))
	assert.Equal(t, "", exampleLines[2])
	assert.Equal(t, "# Pass parameter expected in schema file.", strings.TrimLeft(exampleLines[3], " "))
	assert.Equal(t, "db-unit-extractor extract -s /path/to/schema.yml -n postgres://usr:pwd@127.0.0.1:5432/test "+
		"-r customer_id=4329", strings.TrimLeft(exampleLines[4], " "))

	output := new(bytes.Buffer)
	c.SetArgs([]string{
		"extract", "-s", "../test/unit/schema_test.yml", "-n", "postgres://usr:pwd@127.0.0.1:5432/test",
		"-t", "console", "-d", os.TempDir(), "-r", "test=123",
	})
	c.SetOut(output)

	err := c.Execute()
	assert.Nil(t, err)

	txt := output.String()
	assert.Equal(t, txt, fmt.Sprintf("Extraction is done!\nAssets generated in the directory %s\n", os.TempDir()))
}
