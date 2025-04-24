package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var (
	sampleOutput   string
)

var generateSampleCmd = &cobra.Command{
	Use:   "generate-sample",
	Short: "Create a sample poll JSON file",
	RunE:  generateSample,
}

func init() {
	// generate-sample command flags
	generateSampleCmd.Flags().StringVarP(&sampleOutput, "output", "o", "", "Path to save sample JSON (default: sample_TIMESTAMP.json)")
	RootCmd.AddCommand(generateSampleCmd)
}

// generateSample creates a sample poll JSON and writes it to file
func generateSample(cmd *cobra.Command, args []string) error {
	sample := []map[string]interface{}{ 
		{
			"question": "What is your favorite programming language?",
			"options": []map[string]interface{}{ 
				{"option": "Go", "is_answer": true},
				{"option": "Rust", "is_answer": false},
				{"option": "Python", "is_answer": true},
				{"option": "JavaScript", "is_answer": true},
			},
		},
	}
	b, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to generate sample JSON: %v", err)
	}

	var filename string
	if sampleOutput != "" {
		filename = sampleOutput
	} else {
		timestamp := time.Now().Format("20060102150405")
		filename = fmt.Sprintf("neopoll-%s.json", timestamp)
	}

	if err := os.WriteFile(filename, b, 0600); err != nil {
		return fmt.Errorf("failed to save sample file: %v", err)
	}

	abs, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}
	fmt.Print("Sample JSON generated at: ")
	SuccessColor.Println(abs)
	return nil
}
