package smtp

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// Represents the state of an SMTP session
type State int

// SMTP session states
const (
	INVALID State = -1
	ESTABLISH State = iota
	QUIT
)

// Stores SMTP session state and messages
type Session struct {
	conn       	net.Conn
	state      	State
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		conn: conn,
		state:INVALID,
	}
}

func (s *Session) start() {
	log.Info("Starting an SMTP session")

	s.replyGreeting()
	s.state = ESTABLISH

	//reader := bufio.NewReader(s.conn)

	//for s.state != QUIT {
	//	cmd, err := s.transport.ReadCommand()
	//	if err == nil {
	//		if cmd, ok := s.parseCmd(line); ok {
	//			if !commands[cmd] {
	//				s.send(fmt.Sprintf("500 Syntax error, %v command unrecognized", cmd))
	//				continue
	//			}
	//
	//			switch cmd {
	//			case "NOOP":
	//				s.send("250 I have sucessfully done nothing")
	//				continue
	//			case "QUIT":
	//				s.send("221 Goodnight and good luck")
	//				s.state = QUIT
	//				continue
	//			}
	//		}
	//	}
	//}
	//
	//s.transport.Close()

	log.Info("SMTP session ended")
}

func (s *Session) send(msg string) (err error) {
	deadline := time.Now().Add(time.Duration(10) * time.Second)

	if err := s.conn.SetWriteDeadline(deadline); err != nil {
		log.Errorf("Network send error: %v", err)
		return err
	}

	if _, err := fmt.Fprint(s.conn, msg+"\r\n"); err != nil {
		log.Errorf("Failed to send: \"%v\"", msg)
		return err
	}

	log.Debugf("Sent: \"%v\"", msg)

	return nil
}

func (s *Session) replyGreeting() {
	if err := s.send("220 smtp.briefmx.com ESMTP BriefMX"); err != nil {
		log.Errorf("Failed to send: %v", err)
	}
}
