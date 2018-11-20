package smtp

import (
	"bufio"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

// The underlying transport for Session objects
type Transport struct {
	conn	net.Conn
	reader	*bufio.Reader
}

func NewTransport(conn net.Conn) *Transport {
	return &Transport{
		conn: conn,
		reader: bufio.NewReader(conn),
	}
}

// Reads a line off the wire
func (t *Transport) ReadLine() (line string, err error) {
	line, err = t.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	trimmedLine := strings.TrimRight(line, "\r\n")
	log.Debugf("Read: \"%v\"", trimmedLine)
	return trimmedLine, nil
}

// Sends a line onto the wire
func (t *Transport) SendLine(line string) (err error) {
	// FIXME: probably shouldn't ignore this. maybe buffered writer?
	_, err = t.conn.Write([]byte(line + "\r\n"))
	if err != nil {
		return err
	}
	log.Debugf("Wrote: \"%v\"", line)
	return nil
}
