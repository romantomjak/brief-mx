package smtp

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func waitSig(t *testing.T, c <-chan os.Signal, sig os.Signal) {
	select {
	case s := <-c:
		if s != sig {
			t.Fatalf("Signal was %v, want %v", s, sig)
		}

	case <-time.After(1 * time.Second):
		t.Fatalf("Timeout waiting for %v", sig)
	}
}

func Test_Server_Configuration_Defaults(t *testing.T) {
	config := Config{
		Addr: ":0",
	}
	shutdownChan := make(chan bool)

	smtpServer := NewServer(config, shutdownChan)

	var expected = 0
	if smtpServer.config.ReadTimeout == 0 {
		t.Errorf("Expected ReadTimeout to be '%v', but got '%v'", expected, smtpServer.config.ReadTimeout)
	}

	if smtpServer.config.WriteTimeout == 0 {
		t.Errorf("Expected WriteTimeout to be '%v', but got '%v'", expected, smtpServer.config.WriteTimeout)
	}

	if smtpServer.config.nextCommandTimeout == 0 {
		t.Errorf("Expected nextCommandTimeout to be '%v', but got '%v'", expected, smtpServer.config.nextCommandTimeout)
	}
}

func Test_Server_Responds_To_Signals(t *testing.T) {
	config := Config{
		Addr: ":0",
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	rootCtx, _ := context.WithCancel(context.Background())
	shutdownChan := make(chan bool)

	smtpServer := NewServer(config, shutdownChan)
	go smtpServer.Start(rootCtx)

	waitSig(t, sigChan, syscall.SIGINT)
}
