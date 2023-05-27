package cmd

import (
	"fmt"

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
			_ = cmd.Help()
		},
	}

	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.Flags().BoolP("version", "v", false, fmt.Sprintf("Print %s version", project.name))

	return cmd
}
