package cmd

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/aureliano/caravela"
	"github.com/spf13/cobra"
)

type projectInfo struct {
	name           string
	version        string
	binName        string
	scmHost        string
	scmSsl         bool
	scmProjectPath string
}

var version = "v0.0.0-dev"

var project = projectInfo{
	name:           "db-unit-extractor",
	version:        version,
	binName:        "db-unit-extractor",
	scmHost:        "api.github.com",
	scmSsl:         true,
	scmProjectPath: "aureliano/db-unit-extractor",
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   project.binName,
		Short: "Database extractor for unit testing.",
		Long: "Database extractor for unit testing.\nGo to https://github.com/aureliano/db-unit-extractor/issues " +
			"in order to report a bug or make any suggestion.",
		Run: func(cmd *cobra.Command, args []string) {
			version, _ := cmd.Flags().GetBool("version")
			if version {
				printVersion(cmd)
			} else {
				_ = cmd.Help()
			}
		},
	}

	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.AddCommand(NewUpdateCommand(caravela.Update))
	cmd.AddCommand(NewExtractCommand(caravela.CheckUpdates))

	cmd.Flags().BoolP("version", "v", false, fmt.Sprintf("Print %s version", project.name))

	return cmd
}

func printVersion(cmd *cobra.Command) {
	goVersion := runtime.Version()
	osName := runtime.GOOS
	osArch := runtime.GOARCH

	w := cmd.OutOrStdout()
	write(w, "Version:       %s\n", version)
	write(w, "Go version:       %s\n", goVersion)
	write(w, "OS/Arch:       %s/%s\n", osName, osArch)
}

func shutdown(cmd *cobra.Command, msg string, params ...any) {
	write(cmd.OutOrStdout(), msg, params...)
	os.Exit(1)
}

func write(w io.Writer, msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	_, _ = w.Write([]byte(message))
}
