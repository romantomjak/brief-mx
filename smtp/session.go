package smtp

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type State int

const (
	GREET State = iota // Waiting for HELO
	MAIL               // Got helo, waiting for sender identification
	RCPT               // Sender identified, waiting for recipient(-s)
	DATA               // Got data, waiting for "."
	QUIT               // Close session
)

type Message struct {
	from		string
	recipients	list.List
}

type Session struct {
	server     	*Server
	state      	State
	transport	*Transport
	messages	[]Message
}

func (s *Session) start() {
	log.Info("Starting SMTP session")

	s.greet()

	for s.state != QUIT {
		cmd, err := s.transport.ReadCommand()
		if err == nil {
			if cmd, ok := s.parseCmd(line); ok {
				if !commands[cmd] {
					s.send(fmt.Sprintf("500 Syntax error, %v command unrecognized", cmd))
					continue
				}

				switch cmd {
				case "NOOP":
					s.send("250 I have sucessfully done nothing")
					continue
				case "QUIT":
					s.send("221 Goodnight and good luck")
					s.state = QUIT
					continue
				}
			}
		}
	}

	s.transport.Close()

	log.Info("SMTP session ended")
}

func (s *Session) greet() {
	s.transport.send("220 smtp.briefmx.com ESMTP BriefMX")
}
