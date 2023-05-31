package cmd

import (
	"fmt"
	"os"

	"github.com/aureliano/db-unit-extractor/extractor"
	"github.com/spf13/cobra"

	"regexp"
)

const (
	defaultMaxOpenConn = 3
	defaultMaxIdleConn = 2
)

var (
	refRegExp      = regexp.MustCompile(`^(\w+)\s*=\s*(.+)$`)
	dsnRegExp      = regexp.MustCompile(`^(\w+)://(\w+):(\w+)@([\w.]+):(\d+)/(\w+)\??(\w+=\w+)*$`)
	supportedTypes = []string{"console"}
)

func NewExtractCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract",
		Short: "Extract data-set from database",
		Long:  "Extract data-set from a database to any supported file.",
		Example: fmt.Sprintf(`  # Extract data-set from PostgreSQL and write to the console.
            %s extract -s /path/to/schema.yml -n postgres://usr:pwd@127.0.0.1:5432/test

            # Pass parameter expected in schema file.
            %s extract -s /path/to/schema.yml -n postgres://usr:pwd@127.0.0.1:5432/test -r customer_id=4329`,
			project.binName, project.binName),
		Run: func(cmd *cobra.Command, args []string) {
			extract(cmd)
		},
	}

	cmd.Flags().StringP("schema", "s", "", "Path to the file with the data schema to be extracted.")
	cmd.Flags().StringP("data-source-name", "n", "",
		"Data source name (aka connection string: <driver>://<username>:<password>@<host>:<port>/<database>).")
	cmd.Flags().Int("max-open-conn", defaultMaxOpenConn, "Set the maximum number of concurrently open connections")
	cmd.Flags().Int("max-idle-conn", defaultMaxIdleConn, "Set the maximum number of concurrently idle connections")
	cmd.Flags().StringArrayP("output-type", "t", supportedTypes,
		"Extracted data output format type. Expected: console")
	cmd.Flags().BoolP("formatted-output", "f", false, "Whether the output should be formatted.")
	cmd.Flags().StringP("directory", "d", ".", "Output directory.")
	cmd.Flags().StringArrayP("references", "r", nil, "Expected input parameter in 'schema' file. Expected: name=value")

	_ = cmd.MarkFlagRequired("schema")
	_ = cmd.MarkFlagRequired("data-source-name")

	return cmd
}

func extract(cmd *cobra.Command) {
	conf := extractor.Conf{}
	conf.SchemaPath, _ = cmd.Flags().GetString("schema")
	conf.DSN, _ = cmd.Flags().GetString("data-source-name")
	conf.MaxOpenConn, _ = cmd.Flags().GetInt("max-open-conn")
	conf.MaxIdleConn, _ = cmd.Flags().GetInt("max-idle-conn")
	conf.OutputTypes, _ = cmd.Flags().GetStringArray("output-type")
	conf.FormattedOutput, _ = cmd.Flags().GetBool("formatted-output")
	conf.OutputDir, _ = cmd.Flags().GetString("directory")
	refs, _ := cmd.Flags().GetStringArray("references")

	if err := validateConf(conf); err != nil {
		shutdown(cmd, "Extract error: %s\n", err.Error())
	}

	var err error
	conf.References, err = mapReferences(refs)
	if err != nil {
		shutdown(cmd, "Extract error: %s\n", err.Error())
	}

	if err = extractor.Extract(conf, nil, nil); err != nil {
		shutdown(cmd, "Extract error: %s\n", err.Error())
	}

	_, _ = cmd.OutOrStdout().Write([]byte(extractionSuccessMessage(conf)))
}

func mapReferences(refs []string) (map[string]interface{}, error) {
	mp := make(map[string]interface{})
	for _, ref := range refs {
		matches := refRegExp.FindAllStringSubmatch(ref, -1)
		if matches != nil {
			key := matches[0][1]
			value := matches[0][2]

			mp[key] = value
		} else {
			return nil, fmt.Errorf("invalid reference '%s'", ref)
		}
	}

	return mp, nil
}

func validateConf(conf extractor.Conf) error {
	if !dsnRegExp.MatchString(conf.DSN) {
		return fmt.Errorf("invalid DSN '%s'", conf.DSN)
	}

	info, err := os.Stat(conf.SchemaPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file not found '%s'", conf.SchemaPath)
	} else if info.IsDir() {
		return fmt.Errorf("%s is a directory", conf.SchemaPath)
	}

	for _, in := range conf.OutputTypes {
		if !supportedType(in) {
			return fmt.Errorf("unsupported output type '%s'", in)
		}
	}

	info, err = os.Stat(conf.OutputDir)
	info.Mode().IsDir()
	if !(os.IsNotExist(err) || info.IsDir()) {
		return fmt.Errorf("%s is not a directory", conf.OutputDir)
	}

	return nil
}

func supportedType(tp string) bool {
	for _, st := range supportedTypes {
		if tp == st {
			return true
		}
	}

	return false
}

func extractionSuccessMessage(conf extractor.Conf) string {
	return fmt.Sprintf("Extraction is done!\nAssets generated in the directory %s\n", conf.OutputDir)
}
