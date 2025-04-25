package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/neosh11/survey/myAuth"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

var logoAscii = `
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%*++++%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%      %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%      %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%      %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%     %%-%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%@@@@@%%    %@  %%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%@     @%   %@   %%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%@     @%  %@   %%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%@     @% %@   %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%     %@     @%%@   %%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%@     %@     @%@   % %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%@     %%%    %@   %  %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%@    %% %%  %@   %   %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%@    %   %%%@   %    %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%@    @%   %@   %     %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%@     %%      %      %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%@     %%%   .%%      %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%@     %@-%:*%@%      %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%@     %@  %% @%      %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%     %@     %%      %%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%% %@ %%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%% %@ %%%%%%%%%%%%%
%%%%%%%%%%%%% #  %%%%   %%%%   %%% %  %%%%   %%% %@ %%%%%%%%%%%%%
%%%%%%%%%%%%%     %%  %  %%  -  %%     %%  -  %% %@ %%%%%%%%%%%%%
%%%%%%%%%%%%%  %% #% %%% %  %%% %% %%% *@ %%%  % %@ %%%%%%%%%%%%%
%%%%%%%%%%%%% *%% #*     @  %%% #% %%%  : %%%  % %@ %%%%%%%%%%%%%
%%%%%%%%%%%%% *%% #@ %%%%%  %%% %% %%%  @ %%%  % %@ %%%%%%%%%%%%%
%%%%%%%%%%%%% *%% #%  %+ %%  @  %%  *  %%  @  %% %@ %%%%%%%%%%%%%
%%%%%%%%%%%%% +%% #%%   %%%%   %%% @  %%%%   %%% %@ %%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%% %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

`

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

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `To load completions:

Bash:

  $ source <(neopoll completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ neopoll completion bash > /etc/bash_completion.d/neopoll
  # macOS:
  $ neopoll completion bash > /usr/local/etc/bash_completion.d/neopoll

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  $ neopoll completion zsh > "${fpath[1]}/_neopoll"

Fish:

  $ neopoll completion fish | source

  # To load completions for each session, execute once:
  $ neopoll completion fish > ~/.config/fish/completions/neopoll.fish

PowerShell:

  PS> neopoll completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> neopoll completion powershell > neopoll.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			_ = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			_ = cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

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
	RootCmd.AddCommand(completionCmd)
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

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil && width >= 66 {
		WarnColor.Println(logoAscii)
	}

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
