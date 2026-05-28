package main

import (
	"context"
	"errors"
	"jsson/internal/lsp"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	logFile, err := os.OpenFile("jsson-lsp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logFile)
	log.Println("Starting JSSON Language Server...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	server := lsp.NewServer(os.Stdin, os.Stdout)

	if err := server.Start(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("Server error: %v", err)
		}
	}

	logFile.Close()
}
