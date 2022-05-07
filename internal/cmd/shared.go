package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/muesli/termenv"
	"github.com/x6r/asunder/internal/database"
)

const ttl = 30

var (
	ctx         = context.Background()
	invalidCode = dangerForeground.Render("TOTP secret is invalid")
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkSurvey(err error) {
	if err != nil {
		if err == terminal.InterruptErr {
			fmt.Println("Interrupted")
			os.Exit(1)
		} else {
			fmt.Println(err)
			os.Exit(1)
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
	entries, err := database.GetEntries(db, key)
	check(err)
	if len(entries) < 1 {
		fmt.Println("No entries found.")
		os.Exit(1)
	}

	termenv.ClearScreen()
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Issuer", "Username"})
	fmt.Println("Found the following entries:")
	for _, entry := range entries {
		t.AppendRow(table.Row{entry.ID, strings.Title(entry.Issuer), entry.Username})
	}
	t.Render()
	fmt.Println()
}
