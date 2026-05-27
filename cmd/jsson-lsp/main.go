package main

import (
	"jsson/internal/lsp"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	logFile, err := os.OpenFile("jsson-lsp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logFile)
	log.Println("Starting JSSON Language Server...")

	server := lsp.NewServer(os.Stdin, os.Stdout)

	if err := server.Start(); err != nil {
		log.Printf("Server error: %v", err)
		logFile.Close()

		return
	}

	logFile.Close()
}
