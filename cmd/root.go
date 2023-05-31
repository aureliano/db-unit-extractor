package cmd

import (
	"fmt"
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
		Long:  "Database extractor for unit testing.",
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

	cmd.Flags().BoolP("version", "v", false, fmt.Sprintf("Print %s version", project.name))

	return cmd
}

func printVersion(cmd *cobra.Command) {
	goVersion := runtime.Version()
	osName := runtime.GOOS
	osArch := runtime.GOARCH

	w := cmd.OutOrStdout()
	_, _ = w.Write([]byte(fmt.Sprintf("Version:       %s\n", version)))
	_, _ = w.Write([]byte(fmt.Sprintf("Go version:    %s\n", goVersion)))
	_, _ = w.Write([]byte(fmt.Sprintf("OS/Arch:       %s/%s\n", osName, osArch)))
}

func shutdown(cmd *cobra.Command, msg string, params ...any) {
	msgBytes := []byte(fmt.Sprintf(msg, params...))
	_, _ = cmd.OutOrStdout().Write(msgBytes)

	os.Exit(1)
}
