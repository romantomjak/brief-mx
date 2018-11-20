package smtp

import (
	"net"
	"testing"
)

func TestCanReadALine(t *testing.T) {
	client, server := net.Pipe()
	transport := NewTransport(server)

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

func TestCanSendALine(t *testing.T) {
	client, server := net.Pipe()
	transport := NewTransport(server)

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

// TestWriteTimeout
// TestReadTimeout
