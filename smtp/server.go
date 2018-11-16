package smtp

import (
	"context"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Addr         string        // 0.0.0.0:3535
	ReadTimeout  time.Duration // 10 seconds is a good idea
	WriteTimeout time.Duration // 10 seconds is a good idea

	nextCommandTimeout time.Duration // RFC 5321: 5 minutes
}

// Holds the configuration and state of the SMTP server
type server struct {
	config       Config       // SMTP configuration
	shutdownChan chan bool    // Shuts down Brief MX
	listener     net.Listener // Listens for incoming connections
}

// Creates a new server, but does not start it
func NewServer(config Config, shutdownChan chan bool) *server {
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 10 * time.Second
	}

	if config.WriteTimeout == 0 {
		config.WriteTimeout = 10 * time.Second
	}

	config.nextCommandTimeout = 5 * time.Minute

	return &server{
		config:       config,
		shutdownChan: shutdownChan,
	}
}

// Starts the listener
func (s *server) Start(ctx context.Context) {
	addr, err := net.ResolveTCPAddr("tcp4", s.config.Addr)
	if err != nil {
		log.Errorf("Failed to resolve tcp4 address: %v", err)
		s.shutdown()
		return
	}

	log.Infof("SMTP listening on %s", addr)
	s.listener, err = net.ListenTCP("tcp4", addr)

	go s.listenAndServe(ctx)

	<-ctx.Done()

	// This will stop the listener
	if err := s.listener.Close(); err != nil {
		log.Error("Failed to close SMTP listener")
	}
}

// Accepts new connections
func (s *server) listenAndServe(ctx context.Context) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// FIXME: handle temporary accept errors

			log.Errorf("SMTP accept error: %v", err)
			continue
		}

		go NewSession(conn).start()
	}
}

// Triggers a shutdown
func (s *server) shutdown() {
	close(s.shutdownChan)
}
