package main

import (
	"fmt"
	"os"

	"github.com/neosh11/survey/myAuth"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "neopoll",
	Short: "neopoll is a CLI poll application. Designed with simplicity in mind.",
}

func init() {
	// keep the grouped `auth` namespace…
	RootCmd.AddCommand(myAuth.AuthCmd)
	RootCmd.AddCommand(
		myAuth.LoginCmd,
		myAuth.RefreshTokenCmd,
		myAuth.LogoutCmd,
	)
}
func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
