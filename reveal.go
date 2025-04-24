package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/neosh11/survey/config"
	"github.com/neosh11/survey/myAuth"
	"github.com/spf13/cobra"
)

var revealCmd = &cobra.Command{
	Use:   "reveal",
	Short: "Reveal the answers to the current question",
	RunE:  reveal,
}

func init() {
	RootCmd.AddCommand(revealCmd)
}

func reveal(cmd *cobra.Command, args []string) error {
    accessToken, err := myAuth.GetAccessToken()
	if err != nil {
		return err
	}

	session, err := loadSession()
	if err != nil {
		return fmt.Errorf("no active session: %v", err)
	}

	if session.Password == "" {
		return fmt.Errorf("session password is empty")
	}
	if time.Now().After(session.ExpiresAt) {
		return fmt.Errorf("session expired at %s", session.ExpiresAt)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	reqBody := map[string]interface{}{
		"token":    accessToken,
		"name":     session.Name,
		"password": session.Password,
	}
	b, _ := json.Marshal(reqBody)
	endpoint := config.BaseURL + "/api/reveal"
	resp, err := client.Post(endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var res struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	fmt.Println(res.Message)
	return saveSession(session)
}
