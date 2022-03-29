package main

import (
	"log"

	"github.com/x6r/asunder/internal/cmd"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd.Execute()
}
