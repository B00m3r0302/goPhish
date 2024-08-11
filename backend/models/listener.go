package models

import (
	"net"
	"sync"
	"time"
)

// Listener represents a network listener that waits for connections from agents
type Listener struct {
	ID 				uint 		'json:"id" gorm:"primaryKey"'
	Address 		string 		'json:"address" gorm:"not null"'
	Protocol 		string 		'json:"protocol" gorm:"not null"'
	MaxConnections 	int 		'json:"max_connections" gorm:"not null"'
	Timeout 		int 		'json:"timeout"'
	CreatedAt 		time.Time 	'json"created_at"'
	UpdatedAt 		time.Time 	'json:"updated_at"'

	listener net.Listener // Actual net.Listener object
	active 	 bool 		  // Whether the listener is currently active
	mu 		 sync.Mutex   // Mutex to ensure thread-safe operations
}

// Start initializes and starts the listener, making it ready to accept connections
func (l *Listener) Start() error {
	l.mv.Lock()
	defer l.mu.Unlock()

	ln, err := net.Listen(l.Protocol, l.Address)
	if err != nil {
		return err
	}
	l.Listener = ln
	l.active = true

	go l.acceptConnections()

	return nil
}

// Stop cloeses the listener, stopping it from accepting any more connections
func (l *Listener) Stop() effof {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.listener != nil {
		err := l.listener.Close()
		l.active = false
		return err
	}
	return nil
}

// acceptConnections continuously accepts connections from agents
func (l *Listener) acceptConnections() {
	for {
		conn, err := l.listener.Accept()
		if err != nil {
			if !l.active {
				return // Listener has been stopped
			}
			continue // Handle the error appropriately
		}

		// Handle the connection in a new goroutine
		go l.handleConnection(conn)
	}
}

// handleConnection processes an incoming connection from an agent
func (l *Listener) handleConnection(conn net.Conn) {
	defer conn.Close()

	// set a timeout for the connection if configured
	if l.Timeout > 0 {
		conn.SetDeadline(time.Now().Add(time.Duration(l.Timeout) * time.Second))
	}

	// Placeholder for actual connection handling logic 
	// You could implement handshake, authentication, command processing, etc.
}

// IsActive checks whether the listener is currently active
func (l *Listener) IsActive()
 bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.active
 }
