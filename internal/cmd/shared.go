package cmd

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/muesli/termenv"
)

const ttl = 30

var (
	ctx         = context.Background()
	invalidCode = dangerForeground.Render("TOTP secret is invalid")
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkSurvey(err error) {
	if err != nil {
		if err == terminal.InterruptErr {
			log.Fatalln("Interrupted")
		} else {
			log.Fatalln(err)
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func renderEntriesTable() {
	entries, err := getEntries()
	check(err)
	if len(entries) < 1 {
		log.Fatalln("No entries found.")
	}

	termenv.ClearScreen()
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Issuer", "Username"})
	log.Println("Found the following entries:")
	for _, entry := range entries {
		t.AppendRow(table.Row{entry.ID, strings.Title(entry.Issuer), entry.Username})
	}
	t.Render()
	log.Println()
}
