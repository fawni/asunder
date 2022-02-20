package main

import (
	"log"

	"github.com/x6r/asunder/internal/cmd"
)

func init() {
	log.SetFlags(0)
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
