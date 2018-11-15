package smtp

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var commands = map[string]bool{
	"HELO": true,
	"MAIL": true,
	"RCPT": true,
	"DATA": true,
	"RSET": true,
	"SEND": true,
	"SOML": true,
	"SAML": true,
	"VRFY": true,
	"EXPN": true,
	"HELP": true,
	"NOOP": true,
	"QUIT": true,
	"TURN": true,
}

type Command struct {
	name 	string
}

type Transport struct {
	conn       net.Conn
	reader     *bufio.Reader
}

func NewTransport(conn net.Conn) *Transport {
	return &Transport{
		conn: conn,
	}
}

func (t *Transport) ReadCommand() (cmd *Command, err error) {
	line, err := t.readLine()
	if err != nil {
		log.Errorf("Failed to read line from the client: %v", err)
	}

	cmd, args, err := t.parseCommand(line)
	if err != nil {
		log.Errorf("Failed to parse command: %v", err)
	}

	command := &Command{
		name: cmd,
	}

	return command, nil
}

func (t *Transport) readLine() (line string, err error) {
	deadline := time.Now().Add(time.Duration(10) * time.Second)
	if err = t.conn.SetReadDeadline(deadline); err != nil {
		return "", err
	}
	line, err = t.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	log.Debugf("Read: \"%v\"", strings.TrimRight(line, "\r\n"))
	return line, nil
}

func (t *Transport) parseCommand(line string) (cmd string, args string, err error) {
	line = strings.TrimRight(line, "\r\n")
	l := len(line)
	switch {
	case l == 0:
		return "", true
	case l < 4:
		log.Errorf("Command too short: \"%v\"", line)
		return "", false
	case l == 4:
		return strings.ToUpper(line), true
	case l == 5:
		// Too long to be only command, too short to have args
		log.Errorf("Mangled command: \"%v\"", line)
		return "", false
	}
	return "", false
}

func (t *Transport) send(reply string) {
	deadline := time.Now().Add(time.Duration(10) * time.Second)

	if err := t.conn.SetWriteDeadline(deadline); err != nil {
		log.Errorf("Network send error: %v", err)
		t.server.shutdown()
		return
	}
	if _, err := fmt.Fprint(s.conn, msg+"\r\n"); err != nil {
		log.Errorf("Failed to send: \"%v\"", msg)
		t.server.shutdown()
		return
	}
	log.Debugf("Sent: \"%v\"", msg)
}

func (t *Transport) Close() {
	if err := t.conn.Close(); err != nil {
		log.Warningf("Error while closing connection: %v", err)
	}
}
