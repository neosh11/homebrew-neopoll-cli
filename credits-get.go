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

var getCreditsCmd = &cobra.Command{
	Use:   "credits-get",
	Short: "Get credits for the current user",
	RunE:  getCredits,
}

func init() {
	RootCmd.AddCommand(getCreditsCmd)
}

func getCredits(cmd *cobra.Command, args []string) error {
	accessToken, err := myAuth.GetAccessToken()
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	reqBody := map[string]interface{}{
		"token": accessToken,
	}
	b, _ := json.Marshal(reqBody)

	endpoint := config.BaseURL + "/api/get-remaining-credits"
	resp, err := client.Post(endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var res struct {
		Credits int64  `json:"credits"`
		Error   string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error: %s", res.Error)
	}
	fmt.Print("Remaining credits: ")
	SuccessColor.Println(res.Credits)
	return nil
}
