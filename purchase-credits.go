package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/neosh11/survey/config"
	"github.com/neosh11/survey/myAuth"
	"github.com/spf13/cobra"
)

var purchaseCreditsCmd = &cobra.Command{
	Use:   "buy-credits",
	Short: "Buy more credits for your account",
	RunE:  purchaseCredits,
}

func init() {
	RootCmd.AddCommand(purchaseCreditsCmd)
}

func purchaseCredits(cmd *cobra.Command, args []string) error {
	accessToken, err := myAuth.GetAccessToken()
	if err != nil {
		return err
	}

	// 2) Ask the user which package they want
	var choice string
	prompt := &survey.Select{
		Message: "How many credits would you like to purchase?",
		Options: []string{
			"$5  →  5000 credits",
			"$25 → 50000 credits",
		},
		Default: "$5  →  5000 credits",
	}
	if err := survey.AskOne(prompt, &choice); err != nil {
		return err
	}

	// 3) Translate the choice into a numeric credit amount
	var dollars int
	switch choice {
	case "$5  →  5000 credits":
		dollars = 5
	case "$25 → 50000 credits":
		dollars = 25
	}

	client := &http.Client{Timeout: 10 * time.Second}
	reqBody := map[string]interface{}{
		"token":   accessToken,
		"dollars": dollars,
	}
	b, _ := json.Marshal(reqBody)

	endpoint := config.BaseURL + "/api/purchase-credits"
	resp, err := client.Post(endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var res struct {
		Error     string `json:"error"`
		Message   string `json:"message"`
		SessionId string `json:"session_id"`
		Url       string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error: %s", res.Error)
	}
	fmt.Print("Click here to make your purchase: ")
	SuccessColor.Println(res.Url)
	return nil
}
