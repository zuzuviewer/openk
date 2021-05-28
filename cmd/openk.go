package main

import (
	"log"
	"os"

	"github.com/zuzuviewer/openk/cmd/command"
)

func main() {
	if err := command.RootOpenkCmd.Execute(); err != nil {
		log.Printf("execute err %v", err)
		os.Exit(1)
	}
}
