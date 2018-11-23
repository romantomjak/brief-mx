package smtp

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// The underlying transport for Session objects
type transport struct {
	conn	net.Conn

	timeouts *TransportTimeouts

	reader  *bufio.Reader
	writer *bufio.Writer
}

// Transport layer timeouts
type TransportTimeouts struct {
	read  time.Duration // Socket timeout for read operations
	write time.Duration // Socket timeout for write operations
}

var DefaultTransportTimeouts = &TransportTimeouts{
	read:  30 * time.Second,
	write: 30 * time.Second,
}

var ErrServerLineTooLong = errors.New("Server line too long. Maxiumum 512b")
var ErrClientLineTooLong = errors.New("Client line too long. Maxiumum 1000b")

func NewTransport(conn net.Conn, config *TransportTimeouts) *transport {
	return &transport{
		conn:     conn,
		timeouts: config,
		reader:   bufio.NewReader(conn),
		writer:   bufio.NewWriter(conn),
	}
}

// Reads a line off the wire
func (t *transport) ReadLine() (line string, err error) {
	err = t.conn.SetReadDeadline(time.Now().Add(t.timeouts.read))
	if err != nil {
		return "", err
	}
	line, err = t.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if len(line) > 1000 {
		return "", ErrClientLineTooLong
	}
	trimmedLine := strings.TrimRight(line, "\r\n")
	log.Debugf("Read: \"%v\"", trimmedLine)
	return trimmedLine, nil
}

// Sends a line onto the wire
func (t *transport) SendLine(line string) (err error) {
	if len(line) > 512 {
		return ErrServerLineTooLong
	}

	err = t.conn.SetWriteDeadline(time.Now().Add(t.timeouts.write))
	if err != nil {
		return err
	}
	_, err = t.writer.Write([]byte(line + "\r\n"))
	if err != nil {
		return err
	}
	err = t.writer.Flush()
	if err != nil {
		return err
	}
	log.Debugf("Wrote: \"%v\"", line)
	return nil
}
