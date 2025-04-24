package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/neosh11/survey/config"
	"github.com/neosh11/survey/myAuth"
	"github.com/spf13/cobra"
)

var (
	outputPath string
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the poll session and save the results",
	RunE:  stopPoll,
}

func init() {
	// Add the --output (-o) flag to specify a custom output file path
	stopCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to save results JSON (default: results_TIMESTAMP.json)")
	RootCmd.AddCommand(stopCmd)
}

func stopPoll(cmd *cobra.Command, args []string) error {
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
	endpoint := config.BaseURL + "/api/stop-poll"
	resp, err := client.Post(endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res struct {
		Message string `json:"message"`
		Data    string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	fmt.Println(res.Message)

	if res.Data != "" {
		// Save data, using custom output if provided
		err = saveDataToFile(res.Data, outputPath)
		if err != nil {
			return fmt.Errorf("failed to save data: %v", err)
		}
	}

	return saveSession(session)
}

// saveDataToFile writes the JSON data to the specified outputPath, or defaults to a timestamped file
func saveDataToFile(data string, output string) error {
	var filename string
	if output != "" {
		// Use the provided path
		filename = output
	} else {
		// Generate default filename
		timestamp := time.Now().Format("20060102150405")
		filename = fmt.Sprintf("results_%s.json", timestamp)
	}

	err := os.WriteFile(filename, []byte(data), 0600)
	if err != nil {
		return fmt.Errorf("failed to save results: %v", err)
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}
	fmt.Printf("Results saved to %s\n", absPath)
	return nil
}
