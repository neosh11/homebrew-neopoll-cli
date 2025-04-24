package myAuth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/neosh11/survey/config"
)

// ErrLoginRequired means both your access_token expired and refresh was rejected.
var ErrLoginRequired = errors.New("session invalid or expired; please run `neopoll auth login`")

// sessionFilePath() should be the same helper you already have.
// func sessionFilePath() (string, error) { … }

// GetAccessToken returns a valid access_token, refreshing if needed.
func GetAccessToken() (string, error) {
    s, err:= _getAccessToken()
    if err != nil {
        return "", fmt.Errorf("Failed to get access token, please login again using `neopoll login`.")
    }
    if s == "" {
        return "", fmt.Errorf("Failed to get access token, please login again using `neopoll login`.")
    }
    return s, nil
}

func _getAccessToken() (string, error) {
    path, err := sessionFilePath()
    if err != nil {
        return "", err
    }

    blob, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }

    // match Supabase session shape
    var sess struct {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
        ExpiresAt    int64  `json:"expires_at"`
    }
    if err := json.Unmarshal(blob, &sess); err != nil {
        return "", err
    }

    // still valid?
    if time.Now().Unix() < sess.ExpiresAt {
        return sess.AccessToken, nil
    }

    // expired → try refresh
    return refreshSession(sess.RefreshToken)

}



// refreshSession calls your /api/refresh-token endpoint, saves the new session, and returns the new access_token.
func refreshSession(refreshToken string) (string, error) {
    client := &http.Client{Timeout: 10 * time.Second}
    payload := map[string]string{"refresh_token": refreshToken}
    body, _ := json.Marshal(payload)

    resp, err := client.Post(config.BaseURL+"/api/refresh-token", "application/json", bytes.NewReader(body))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // backend uses 403 to signal “you must re-login”
    if resp.StatusCode == http.StatusForbidden {
        return "", ErrLoginRequired
    }
    if resp.StatusCode >= 400 {
        var e struct{ Error string `json:"error"` }
        _ = json.NewDecoder(resp.Body).Decode(&e)
        return "", fmt.Errorf("refresh failed: %s", e.Error)
    }

    // decode new session blob
    var res struct {
        Session json.RawMessage `json:"session"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
        return "", err
    }

    // overwrite on disk
    path, err := sessionFilePath()
    if err != nil {
        return "", err
    }
    if err := os.WriteFile(path, res.Session, fs.FileMode(0600)); err != nil {
        return "", err
    }

    // extract the new access token
    var newSess struct {
        AccessToken string `json:"access_token"`
        ExpiresAt   int64  `json:"expires_at"`
    }
    if err := json.Unmarshal(res.Session, &newSess); err != nil {
        return "", err
    }
    return newSess.AccessToken, nil
}

func postAndPrint(url string, payload interface{}) error {
    client := &http.Client{Timeout: 10 * time.Second}
    body, _ := json.Marshal(payload)
    resp, err := client.Post(url, "application/json", bytes.NewReader(body))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var res struct {
        Message string `json:"message"`
        Error   string `json:"error"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
        return err
    }
    if resp.StatusCode >= 400 {
        return fmt.Errorf("error: %s", res.Error)
    }
    fmt.Println(res.Message)
    return nil
}

func sessionFilePath() (string, error) {
    dir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    cfgDir := filepath.Join(dir, ".neopoll")
    if err := os.MkdirAll(cfgDir, 0700); err != nil {
        return "", err
    }
    return filepath.Join(cfgDir, "session.json"), nil
}
