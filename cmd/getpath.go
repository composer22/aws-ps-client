package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/composer22/aws-ps-client/client"
	"github.com/spf13/cobra"
)

// getPathCmd represents the get command for retrieving variables by
// directoy structure.
var getPathCmd = &cobra.Command{
	Use:   "getpath",
	Short: "Get key/values from a directory in the AWS Parameter Store",
	Long:  "Given a directory path, retrieve the values from the AWS Parameter Store",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("AWS Parameter Store directory path is mandatory.")
			os.Exit(0)
		}
		k := args[0]
		v := cmd.Flag("version").Value.String()
		r, err := strconv.ParseBool(cmd.Flag("recursive").Value.String())
		if err != nil {
			client.PrintErr(err.Error())
			os.Exit(0)
		}

		params, err := client.New(awsAccessKey, awsAccessSecret,
			awsRegion).QueryPath(k, r, v)
		if err != nil {
			client.PrintErr(err.Error())
			os.Exit(0)
		}
		printGetPathCmd(params, k)
	},
	Example: `aws-ps-client getpath /path/to/keys/ --aws-access-key bar --aws-access-secret letmein --aws-region us-west-2
aws-ps-client getpath /path/to/keys/ -k bar -s letmein --format=json
aws-ps-client getpath /path/to/keys/ -k /path/to/token.1line.file -s /path/to/secret.1line.file
aws-ps-client getpath /path/to/keys/ --version 12 -f text --recursive false
aws-ps-client getpath /path/to/keys/ -v 12 -f bash
`,
}

func printGetPathCmd(params client.Parameters, prefix string) {
	switch format {
	case "bash":
		for _, p := range params {
			name := strings.Replace(p.Name, prefix, "", 1)
			fmt.Printf("export %s=\"%s\"\n", name, p.Value)
		}
	case "json":
		b, err := json.MarshalIndent(params, "", "\t")
		if err != nil {
			client.PrintErr(err.Error())
			os.Exit(0)
		}
		fmt.Printf("{\"result\":\"%s\"}\n", string(b))
	case "text":
		for _, p := range params {
			fmt.Printf("\"%s\"=\"%s\"\n", p.Name, p.Value)
		}
	default:
		client.PrintErr("Invalid format flag")
		os.Exit(0)
	}
	return
}

func init() {
	RootCmd.AddCommand(getPathCmd)
	getPathCmd.SetUsageTemplate(getPathUsageTemplate())
	getPathCmd.Flags().StringP("version", "v", "",
		"Version of the AWS Parameter Store key/value pairs to retrieve")
	getPathCmd.Flags().BoolP("recursive", "u", true,
		"Recurse through the subdirectories for keys")
}

// Override help template.
func getPathUsageTemplate() string {
	return `Usage:{{if .Runnable}}
  {{if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "[flags] DIR-PATH"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
  {{ .CommandPath}} [command]{{end}}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

}
