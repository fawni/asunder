package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fawni/asunder/internal/common"
	"github.com/fawni/asunder/internal/database"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete an entry",
		Long:  `The delete command deletes an entry from the database`,
		PreRun: func(cmd *cobra.Command, args []string) {
			connectDB()
		},
		Run: func(cmd *cobra.Command, args []string) {
			deleteEntry()
		},
	}
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteEntry() {
	renderEntriesTable()
	var id int
	err := survey.AskOne(&survey.Input{
		Message: "Enter ID to delete »",
	}, &id)
	common.CheckSurvey(err)

	common.Check(err)
	var delete bool
	err = survey.AskOne(&survey.Confirm{Message: fmt.Sprintf("Delete [%s]", getDescriptionById(id))}, &delete)
	common.CheckSurvey(err)

	if delete {
		ctx := context.Background()
		entry := &database.Entry{ID: id}
		_, err := DB.NewDelete().Model(entry).WherePK().Exec(ctx)
		common.Check(err)
		fmt.Println("Done!")
	} else {
		fmt.Println("Cancelled!")
		os.Exit(0)
	}
}

func getDescriptionById(id int) string {
	ctx := context.Background()
	entry := new(database.Entry)
	err := DB.NewSelect().Model(entry).Where("id = ?", id).Scan(ctx)
	common.Check(err)
	issuer, err := database.Decrypt(Key, entry.Issuer)
	common.Check(err)
	username, err := database.Decrypt(Key, entry.Username)
	common.Check(err)
	return fmt.Sprintf("%s - %s", strings.Title(issuer), username)
}

func renderEntriesTable() {
	entries, err := database.GetEntries(DB, Key)
	common.Check(err)
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
