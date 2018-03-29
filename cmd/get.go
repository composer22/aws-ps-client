package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/composer22/aws-ps-client/client"
	"github.com/spf13/cobra"
)

// getCmd represents the get command for retrieving a single k/v.
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a value from the AWS Parameter Store for a key",
	Long:  "Given a key, retrieve the value from the AWS Parameter Store",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("AWS Parameter Store key is mandatory.")
			os.Exit(0)
		}
		k := args[0]
		v := cmd.Flag("version").Value.String()
		param, err := client.New(awsAccessKey, awsAccessSecret,
			awsRegion).Query(k, v)
		if err != nil {
			client.PrintErr(err.Error())
			os.Exit(0)
		}
		printGetCmd(param, k)
	},
	Example: `aws-ps-client get /path/with/KEY1 --aws-access-key bar --aws-access-secret letmein --aws-region us-west-2
aws-ps-client get /path/with/KEY1  -k foo -s letmein --format=json
aws-ps-client get /path/with/KEY1  -k /path/to/token.1line.file -s /path/to/secret.1line.file
aws-ps-client get /path/with/KEY1  --version 12 -f text
aws-ps-client get /path/with/KEY1  -v 12 -f bash
`,
}

func printGetCmd(param *client.Parameter, prefix string) {
	switch format {
	case "bash":
		name := filepath.Base(param.Name)
		fmt.Printf("export %s=\"%s\"\n", name, param.Value)
	case "json":
		b, err := json.MarshalIndent(param, "", "\t")
		if err != nil {
			client.PrintErr(err.Error())
			os.Exit(0)
		}
		fmt.Printf("{\"result\":\"%s\"}\n", string(b))
	case "text":
		fmt.Printf("\"%s\"=\"%s\"\n", param.Name, param.Value)
	default:
		client.PrintErr("Invalid format flag")
		os.Exit(0)
	}
	return
}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.SetUsageTemplate(getUsageTemplate())
	getCmd.Flags().StringP("version", "v", "",
		"Version of the AWS Parameter Store key/value to retrieve")
}

// Override help template.
func getUsageTemplate() string {
	return `Usage:{{if .Runnable}}
  {{if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "[flags] KEY"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
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
