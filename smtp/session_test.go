package smtp

import (
	"bufio"
	"bytes"
	"net"
	"testing"
)

func Test_Session_Starts_In_An_Invalid_State(t *testing.T) {
	server, _ := net.Pipe()

	session := NewSession(server)

	expected := INVALID
	if session.state != expected {
		t.Errorf("Expected session state to be '%v', but got '%v'", expected, session.state)
	}
}

func Test_Sends_Greeting_On_Connection(t *testing.T) {
	server, client := net.Pipe()

	session := NewSession(server)

	go func() {
		session.start()
		server.Close()
	}()

	reader := bufio.NewReader(client)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		t.Error("Could not read from connection")
	}

	expectedGreeting := []byte("220 smtp.briefmx.com ESMTP BriefMX\r\n")
	if !bytes.Equal(line, expectedGreeting) {
		t.Errorf("Expected SMTP greeting to be '%s', but got '%v'", expectedGreeting, line)
	}

	expectedState := INVALID
	if session.state != expectedState {
		t.Errorf("Expected session state to be '%v', but got '%v'", expectedState, session.state)
	}

	client.Close()
}
