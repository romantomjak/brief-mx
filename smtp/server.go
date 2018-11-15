package smtp

import (
    "context"
    "net"

    log "github.com/sirupsen/logrus"
)

type Config struct {
    Addr string // 0.0.0.0:3535
}

// Holds the configuration and state of the SMTP server
type Server struct {
    config Config           // SMTP configuration
    shutdownChan chan bool  // Shuts down Brief MX
    listener net.Listener   // Listens for incoming connections
}

// Creates a new server, but does not start it
func NewServer(config Config, shutdownChan chan bool,) *Server {
    return &Server{
        config: config,
        shutdownChan: shutdownChan,
    }
}

// Starts the listener
func (s *Server) Start(ctx context.Context) {
    addr, err := net.ResolveTCPAddr("tcp4", s.config.Addr)
    if err != nil {
        log.Errorf("Failed to resolve tcp4 address: %v", err)
        s.shutdown()
        return
    }

    log.Infof("SMTP listening on %s", addr)
    s.listener, err = net.ListenTCP("tcp4", addr)

    go s.accept(ctx)

    <-ctx.Done()

    // This will stop the listener
    if err := s.listener.Close(); err != nil {
        log.Error("Failed to close SMTP listener")
    }
}

// Accepts new connections
func (s *Server) accept(ctx context.Context) {
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            log.Errorf("SMTP accept error: %v", err)
            continue
        }

        session := &Session{
            server:     s,
            state:      GREET,
            transport:  NewTransport(conn),
        }

        go session.start()
    }
}

// Triggers a shutdown
func (s *Server) shutdown() {
    close(s.shutdownChan)
}
