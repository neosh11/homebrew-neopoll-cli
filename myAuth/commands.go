package myAuth

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/neosh11/survey/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	AuthCmd         = &cobra.Command{Use: "auth", Short: "authentication commands"}
	LoginCmd        = &cobra.Command{Use: "login", Short: "Request OTP & verify", RunE: runLogin}
	RefreshTokenCmd = &cobra.Command{Use: "refresh-token", Short: "Refresh saved session", RunE: runRefreshToken}
	LogoutCmd       = &cobra.Command{Use: "logout", Short: "Delete saved session", RunE: runLogout}

	// flags
	flagEmail    string
	flagPassword string
	flagToken    string

	// Colors
	promptColor  = color.New(color.FgCyan, color.Bold)
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgYellow)
)

func init() {
	AuthCmd.AddCommand(LoginCmd, RefreshTokenCmd, LogoutCmd)
	LoginCmd.Flags().StringVarP(&flagEmail, "email", "e", "", "Email address")
	LoginCmd.Flags().StringVarP(&flagToken, "token", "t", "", "OTP token (if you already have it)")
}

// —— prompts
func promptEmail() error {
	if flagEmail != "" {
		return nil
	}
	promptColor.Print("Email: ")
	inp, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return err
	}
	flagEmail = strings.TrimSpace(inp)
	if flagEmail == "" {
		return errors.New("email cannot be empty")
	}
	return nil
}

func promptPassword() error {
	if flagPassword != "" {
		return nil
	}
	promptColor.Print("Password: ")
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return err
	}
	flagPassword = strings.TrimSpace(string(b))
	if flagPassword == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}

func promptToken() error {
	if flagToken != "" {
		return nil
	}
	promptColor.Print("Enter OTP: ")
	inp, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return err
	}
	flagToken = strings.TrimSpace(inp)
	if flagToken == "" {
		return errors.New("token cannot be empty")
	}
	return nil
}

//
// —— commands
//

func runLogin(cmd *cobra.Command, args []string) error {

	if err := promptEmail(); err != nil {
		return err
	}

	if flagToken == "" {
		if err := postAndPrint(config.BaseURL+"/api/login", map[string]string{"email": flagEmail}); err != nil {
			return err
		}
		// 2) verify OTP & persist session
		if err := promptToken(); err != nil {
			return err
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	body, _ := json.Marshal(map[string]string{"email": flagEmail, "token": flagToken})
	resp, err := client.Post(config.BaseURL+"/api/otp", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res struct {
		Message string          `json:"message"`
		Session json.RawMessage `json:"session"`
		Error   string          `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(res.Error)
	}

	// save session.json
	path, err := sessionFilePath()
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, res.Session, fs.FileMode(0600)); err != nil {
		return err
	}
	fmt.Println(res.Message)
	successColor.Println("✔ session saved to %s\n", path)
	return nil
}

func runRefreshToken(cmd *cobra.Command, args []string) error {
	// load existing session
	path, err := sessionFilePath()
	if err != nil {
		return err
	}
	blob, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// extract refresh_token
	var old struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(blob, &old); err != nil {
		return err
	}
	if old.RefreshToken == "" {
		return errors.New("no refresh_token in saved session")
	}

	// call refresh endpoint
	client := &http.Client{Timeout: 10 * time.Second}
	body, _ := json.Marshal(map[string]string{"refresh_token": old.RefreshToken})
	resp, err := client.Post(config.BaseURL+"/api/refresh-token", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res struct {
		Message string          `json:"message"`
		Session json.RawMessage `json:"session"`
		Error   string          `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(res.Error)
	}

	// overwrite session file
	if err := os.WriteFile(path, res.Session, fs.FileMode(0600)); err != nil {
		return err
	}
	fmt.Println(res.Message)
	successColor.Println("✔ session refreshed at %s\n", path)
	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	path, err := sessionFilePath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			warnColor.Println("⚠ no session found")
			return nil
		}
		return err
	}
	successColor.Println("✔ logged out, session deleted")
	return nil
}
