package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/neosh11/survey/myAuth"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	// define your palettes
	header  = color.New(color.FgCyan, color.Bold)
	cmdName = color.New(color.FgGreen, color.Bold)
	helpTxt = color.New(color.FgWhite)
)

var (
	PromptColor  = color.New(color.FgCyan, color.Bold)
	SuccessColor = color.New(color.FgGreen)
	WarnColor    = color.New(color.FgYellow)
)

var RootCmd = &cobra.Command{
	Use:   "neopoll",
	Short: "neopoll is a CLI poll application. Designed with simplicity in mind.",
	Long: `neopoll is a CLI poll application. Designed with simplicity in mind.

Quickstart:

  1. Authenticate
     $ neopoll login

  2. Generate a sample file
     $ neopoll generate-sample --output=my-poll.json

  3. Start the poll
     $ neopoll start my-poll.json

  4. Control the poll
     $ neopoll next
     $ neopoll reveal
     $ neopoll stop --output=results.json

Type "neopoll help <command>" for more details on a specific command.
`,
}

func init() {
	// keep the grouped `auth` namespace…
	RootCmd.AddCommand(myAuth.AuthCmd)
	RootCmd.AddCommand(
		myAuth.LoginCmd,
		myAuth.RefreshTokenCmd,
		myAuth.LogoutCmd,
	)

	RootCmd.SetHelpFunc(prettyHelp)
}
func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func prettyHelp(cmd *cobra.Command, args []string) {
	header.Println("\nUSAGE:")
	cmdName.Printf("  %s\n\n", cmd.UseLine())

	header.Println("COMMANDS:")
	for _, c := range cmd.Commands() {
		cmdName.Printf("  %-15s", c.Name())
		helpTxt.Printf(" %s\n", c.Short)
	}

	header.Println("\nFLAGS:")
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		cmdName.Printf("  --%-13s", f.Name)
		helpTxt.Printf(" %s\n", f.Usage)
		if f.DefValue != "" {
			helpTxt.Printf(" (default: %s)\n", f.DefValue)
		} else {
			helpTxt.Println()
		}
	})

	header.Println("\nQUICKSTART:")
	helpTxt.Println(`
  1. Authenticate:
     neopoll login

  2. Generate sample:
     neopoll generate-sample --output=my-poll.json

  3. Start poll:
     neopoll start my-poll.json

  4. Control poll:
     neopoll next
     neopoll reveal
     neopoll stop --output=results.json
`)
}
