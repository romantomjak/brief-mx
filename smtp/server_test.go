package smtp

import (
	"testing"
)

func newServer() *server {
	config := Config{
		Addr: "127.0.0.1:0",
	}

	shutdownChan := make(chan bool)

	return NewServer(config, shutdownChan)
}

func Test_Default_Server_Timeout_Values_Are_Set(t *testing.T) {
	expected := 0
	server := newServer()

	if server.config.ReadTimeout == 0 {
		t.Errorf("Expected ReadTimeout to be '%v', but got '%v'", expected, server.config.ReadTimeout)
	}

	if server.config.WriteTimeout == 0 {
		t.Errorf("Expected WriteTimeout to be '%v', but got '%v'", expected, server.config.WriteTimeout)
	}

	if server.config.nextCommandTimeout == 0 {
		t.Errorf("Expected nextCommandTimeout to be '%v', but got '%v'", expected, server.config.nextCommandTimeout)
	}
}
