package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Answer struct {
	QuestionIndex int         `json:"question_index"`
	Response      interface{} `json:"response"`
}

type Session struct {
	Name       string    `json:"name"`
	StartedAt  time.Time `json:"started_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	PollPath   string    `json:"poll_path"`
	Answers    []Answer  `json:"answers"`
	CurrentIdx int       `json:"current_index"`
	Password   string    `json:"password"`
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
    return filepath.Join(cfgDir, "poll-session.json"), nil
}

func saveSession(s Session) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
    sessionFilePath, err := sessionFilePath()
    if err != nil {
        return err
    }
	return os.WriteFile(sessionFilePath, b, 0600)
}

func loadSession() (Session, error) {
	var s Session
    sessionFilePath, err := sessionFilePath()
    if err != nil {
        return s, err
    }
	data, err := os.ReadFile(sessionFilePath)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(data, &s)
	return s, err
}
