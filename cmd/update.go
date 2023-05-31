package cmd

import (
	"fmt"

	"github.com/aureliano/caravela"
	"github.com/aureliano/caravela/provider"
	"github.com/spf13/cobra"
)

func NewUpdateCommand(upcmd func(c caravela.Conf) (*provider.Release, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates this program",
		Long:  "Checks for a newer version of this program and updates it if necessary.",
		Run: func(cmd *cobra.Command, args []string) {
			update(cmd, upcmd)
		},
	}

	return cmd
}

func update(cmd *cobra.Command, upcmd func(c caravela.Conf) (*provider.Release, error)) {
	release, err := upcmd(caravela.Conf{
		Version: project.version,
		Provider: provider.GithubProvider{
			Host:        project.scmHost,
			Ssl:         project.scmSsl,
			ProjectPath: project.scmProjectPath,
		},
	})

	w := cmd.OutOrStdout()

	if err != nil {
		shutdown(cmd, "Program update failed! %s\n", err.Error())
	}

	_, _ = w.Write([]byte(fmt.Sprintf("Release %s of %s.\n\n", release.Name,
		release.ReleasedAt.Format("02/01/2006 15:04:05"))))
	_, _ = w.Write([]byte(fmt.Sprintln(release.Description)))
	_, _ = w.Write([]byte(fmt.Sprintf("\nUpdate from version %s to %s successfully completed!\n",
		project.version, release.Name)))
}
