package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type VerificationResult struct {
	Timestamp time.Time `json:"timestamp"`
	AppPath   string    `json:"app_path"`
	Verified  bool      `json:"verified"`
	Output    string    `json:"output"`
	Error     string    `json:"error,omitempty"`
	Hostname  string    `json:"hostname"`
}

func verifyApp(path string) VerificationResult {
	cmd := exec.Command("codesign", "--verify", "--deep", "--strict", "--verbose=2", path)
	output, err := cmd.CombinedOutput()
	hostname, _ := os.Hostname() // Captura o hostname da workstation
	result := VerificationResult{
		Timestamp: time.Now(),
		AppPath:   path,
		Verified:  err == nil,
		Output:    string(output),
		Hostname:  hostname,
	}

	if err != nil {
		result.Error = fmt.Sprintf("verification failed: %s", err)
		log.Printf("Verification failed for %s: %v\nOutput: %s\n", path, err, string(output))
	} else {
		log.Printf("Verification successful for %s: %s\n", path, string(output))
	}

	return result
}

func logResultToJSON(result VerificationResult) {
	logData, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}

	// Imprime o log em formato JSON, que será capturado pelo FluentBit
	fmt.Println(string(logData))
}

func main() {
	appsToVerify := []string{
		"/Applications/Google Chrome.app",
		"/Applications/Safari.app",
	}

	for _, appPath := range appsToVerify {
		result := verifyApp(appPath)
		logResultToJSON(result)
	}
}
