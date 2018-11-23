package smtp

import (
	"net"
	"testing"
)

func TestReadsBytesToString(t *testing.T) {
	client, server := net.Pipe()
	transport := NewTransport(server, DefaultTransportTimeouts)

	go func() {
		_, err := client.Write([]byte("HELLO SERVER\r\n"))
		if err != nil {
			t.Errorf("Failed to write to a pipe: %v", err)
		}
	}()

	expected := "HELLO SERVER"

	line, err := transport.ReadLine()
	if err != nil {
		t.Errorf("Failed to read from a pipe: %v", err)
	}

	if line != expected {
		t.Errorf("Expected '%s' to match '%s'", expected, line)
	}
}

func TestWritesStringToBytes(t *testing.T) {
	client, server := net.Pipe()
	transport := NewTransport(server, DefaultTransportTimeouts)

	go func() {
		err := transport.SendLine("HELLO CLIENT")
		if err != nil {
			t.Errorf("Failed to write to a pipe: %v", err)
		}
	}()

	expected := "HELLO CLIENT\r\n"

	line := make([]byte, len(expected))
	_, err := client.Read(line)
	if err != nil {
		t.Errorf("Failed to read from a pipe: %v", err)
	}

	if string(line[:]) != expected {
		t.Errorf("Expected '%s' to match '%s'", expected, line)
	}
}

func TestReadTimeout(t *testing.T) {
	_, server := net.Pipe()

	config := &TransportTimeouts{
		read:  1,
		write: 0,
	}
	transport := NewTransport(server, config)

	_, err := transport.ReadLine()
	if err == nil {
		t.Errorf("Expected: %v", err)
	}

	if nerr, ok := err.(net.Error); ok {
		if !nerr.Timeout() {
			t.Errorf("err.Timeout() = false, want true")
		}
	} else {
		t.Errorf("got %T, want net.Error", err)
	}
}

func TestWriteTimeout(t *testing.T) {
	_, server := net.Pipe()

	config := &TransportTimeouts{
		read:  0,
		write: 1,
	}
	transport := NewTransport(server, config)

	err := transport.SendLine("YOLO")
	if err == nil {
		t.Errorf("Expected: %v", err)
	}

	if nerr, ok := err.(net.Error); ok {
		if !nerr.Timeout() {
			t.Errorf("err.Timeout() = false, want true")
		}
	} else {
		t.Errorf("got %T, want net.Error", err)
	}
}

func TestClientMaximumLineLength(t *testing.T) {
	client, server := net.Pipe()
	transport := NewTransport(server, DefaultTransportTimeouts)

	go func() {
		b := make([]byte, 1001)
		b[999] = '\r'
		b[1000] = '\n'
		_, err := client.Write(b)
		if err != nil {
			t.Errorf("Failed to write to a pipe: %v", err)
		}
	}()

	_, err := transport.ReadLine()
	if err != ErrClientLineTooLong {
		t.Errorf("got %T, want smtp.ErrClientLineTooLong", err)
	}
}

func TestServerMaximumLineLength(t *testing.T) {
	_, server := net.Pipe()
	transport := NewTransport(server, DefaultTransportTimeouts)

	b := make([]byte, 513)
	err := transport.SendLine(string(b))

	if err != ErrServerLineTooLong {
		t.Errorf("got %T, want smtp.ErrServerLineTooLong", err)
	}
}
