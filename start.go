package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/neosh11/survey/config"
	"github.com/neosh11/survey/myAuth"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [poll.json]",
	Short: "Start a new poll session",
	Args:  cobra.ExactArgs(1),
	RunE:  startPoll,
}

func init() {
	RootCmd.AddCommand(startCmd)
}

func startPoll(cmd *cobra.Command, args []string) error {
	accessToken, err := myAuth.GetAccessToken()
	if err != nil {
		return err
	}

	jsonPath := args[0]
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	PromptColor.Print("Enter poll password: ")
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	password = strings.TrimSpace(password)

	reqBody := map[string]interface{}{
		"token":    accessToken,
		"password": password,
		"poll":     json.RawMessage(data),
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	fmt.Println("Starting poll session...")

	client := &http.Client{Timeout: 10 * time.Second}
	endpoint := config.BaseURL + "/api/start-poll"
	resp, err := client.Post(endpoint, "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// If non-200, decode error JSON and report
	if resp.StatusCode != http.StatusOK {
		var errRes struct {
			Error string `json:"error"`
		}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&errRes); decodeErr == nil && errRes.Error != "" {
			return fmt.Errorf(WarnColor.Sprintf("error starting poll: %s", errRes.Error))
		}
		// fallback to raw body
		raw, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf(WarnColor.Sprintf("error starting poll: %s", strings.TrimSpace(string(raw))))
	}

	var res struct {
		Name     string `json:"name"`
		PollLink string `json:"poll_link"`
		QRLink   string `json:"qr_link"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	fmt.Print("✔ Poll link: ")
    SuccessColor.Println(res.PollLink)
	fmt.Print("✔ QR code link : ")
    SuccessColor.Println(res.QRLink)

	session := Session{
		Name:       res.Name,
		StartedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(1 * time.Hour),
		PollPath:   jsonPath,
		Answers:    []Answer{},
		CurrentIdx: 0,
		Password:   password,
	}
	return saveSession(session)
}
