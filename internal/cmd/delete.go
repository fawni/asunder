package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/muesli/coral"
	"github.com/x6r/asunder/internal/database"
)

var (
	deleteCmd = &coral.Command{
		Use:   "delete",
		Short: "Delete an entry",
		Long:  `The delete command deletes an entry from the database`,
		PreRun: func(cmd *coral.Command, args []string) {
			connectDB()
		},
		Run: func(cmd *coral.Command, args []string) {
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
		Message: "Enter ID to delete ›",
	}, &id)
	checkSurvey(err)

	check(err)
	var delete bool
	err = survey.AskOne(&survey.Confirm{Message: fmt.Sprintf("Delete [%s]", getDescriptionById(id))}, &delete)
	checkSurvey(err)

	if delete {
		entry := &database.Entry{ID: id}
		_, err := db.NewDelete().Model(entry).WherePK().Exec(ctx)
		check(err)
		fmt.Println("Done!")
	} else {
		fmt.Println("Cancelled!")
		os.Exit(0)
	}
}

func getDescriptionById(id int) string {
	entry := new(database.Entry)
	err := db.NewSelect().Model(entry).Where("id = ?", id).Scan(ctx)
	check(err)
	issuer, err := database.Decrypt(key, entry.Issuer)
	check(err)
	username, err := database.Decrypt(key, entry.Username)
	check(err)
	return fmt.Sprintf("%s - %s", strings.Title(issuer), username)
}
