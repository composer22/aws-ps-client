package cmd

import (
	"fmt"
	"os"

	"github.com/composer22/aws-ps-client/client"
	"github.com/spf13/cobra"
)

// versionCmd returns the version of the application
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of the application",
	Long:  "Returns the version of the application",
	Run: func(cmd *cobra.Command, args []string) {
		result := client.New("", "", "").Version()
		printVersion(result)
	},
	Example: `aws-ps-client version --format=json
aws-ps-client version -f bash
aws-ps-client version
`,
}

func printVersion(result string) {
	switch format {
	case "bash":
		fmt.Printf("echo \"%s\"\n", result)
	case "json":
		fmt.Printf("{\"version\":\"%s\"}\n", result)
	case "text":
		fmt.Printf("%s\n", result)
	default:
		client.PrintErr("Invalid format flag")
		os.Exit(0)
	}
}

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.SetUsageTemplate(versionUsageTemplate())
}

// Override help template.
func versionUsageTemplate() string {
	return `Usage:{{if .Runnable}}
  {{if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "[flags]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
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
