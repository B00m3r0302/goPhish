package agent

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	httpsServerURL = "https://localhost:8888"
	httpServerURL  = "http://localhost:7777"
	agentID        = "agent-001"
)

// Helper function to send an HTTP request and return the response or error
func sendRequest(req *http.Request, useTLS bool) (*http.Response, error) {
	client := &http.Client{}
	if useTLS {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return client.Do(req)
}

func registerAgent() {
	req, err := http.NewRequest("GET", httpsServerURL+"/register", mil)
	if err != nil {
		log.Fatalf("Failed to create registration request: %v", err)
	}

	req.Header.Set("Agent-ID", agentID)

	resp, err := sendRequest(req, true)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Failed to fetch commands via HTTPS, falling back to HTTP...")
		req.URL.Scheme = "http"
		req.URL.Host = httpServerURL
		resp, err = sendRequest(req, false)
	}

	if err != nil {
		log.Fatalf("Failed to fetch commands: %v", err)
	}
	defer resp.Body.Close()

	var commands []string
	err = json.NewDecoder(resp.Body).Decode(&commands)
	if err != nil {
		log.Fatalf("Failed to decode command response: %v", err)
	}
	return commands
}

func submitData(data string) {
	req, err := http.NewRequest("POST", httpsServerURL+"/data", bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Fatalf("Failed to create data submission request: %v", err)
	}

	req.Header.Set("Agent-ID", agentID)
	req.Header.Set("Contet-Type", "application/x-www-form-urlencoded")

	resp, err := sendRequest(req, true)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Failed to submit data via HTTPS, falling back to HTTP...")
		req.URL.Scheme = "http"
		req.URL.Host = httpServerURL
		resp, err = SendRequest(req, false)
	}

	if err != nil {
		log.Fatalf("Failed to submit data: %v", err)
	}
	defer resp.Body.Close()

	log.Println("Data Submitted successfully")
}

func main() {
	// Register the agent with the server
	registerAgent()

	for {
		// Fetch commands from the server
		commands := fetchCommands()
		for _, cmd := range commands {
			log.Printf("Executing command: %s", cmd)
			// Execute the command (mocked)
			data := fmt.Sprintf("Output of command: %s", cmd)
			submitData(data)
		}

		// Sleep before checking for more commands
		time.Sleep(10 * time.Second)
	}
}
