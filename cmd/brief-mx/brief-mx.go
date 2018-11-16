package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/romantomjak/brief-mx/smtp"
)

func main() {
	log.Info("Starting server")

	config := smtp.Config{
		Addr: "127.0.0.1:3535",
	}

	rootCtx, rootCancel := context.WithCancel(context.Background())
	shutdownChan := make(chan bool)

	smtpServer := smtp.NewServer(config, shutdownChan)
	go smtpServer.Start(rootCtx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

WAIT:
	for {
		select {
		case sig := <-sigChan:
			switch sig {
			case syscall.SIGINT:
				log.Info("Received SIGINT, shutting down")
				close(shutdownChan)
			case syscall.SIGTERM:
				log.Info("Received SIGTERM, shutting down")
				close(shutdownChan)
			}
		case <-shutdownChan:
			rootCancel()
			break WAIT
		}
	}
}
