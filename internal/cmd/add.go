package cmd

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/muesli/coral"
	"github.com/x6r/asunder/internal/database"
)

var (
	addCmd = &coral.Command{
		Use:   "add",
		Short: "Add an entry",
		Long:  `The add command adds entries to the database`,
		RunE: func(cmd *coral.Command, args []string) error {
			return addEntry()
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}

func addEntry() error {
	var qs = []*survey.Question{
		{
			Name:     "username",
			Prompt:   &survey.Input{Message: "Enter your username ›"},
			Validate: survey.Required,
		},
		{
			Name:     "issuer",
			Prompt:   &survey.Input{Message: "Enter the issuer name ›"},
			Validate: survey.Required,
		},
		{
			Name:     "secret",
			Prompt:   &survey.Password{Message: "Enter TOTP secret ›"},
			Validate: survey.Required,
		},
	}

add:
	var answers database.Entry
	err := survey.Ask(qs, &answers)
	checkSurvey(err)
	username, err := database.Encrypt(key, answers.Username)
	check(err)
	issuer, err := database.Encrypt(key, answers.Issuer)
	check(err)
	secret, err := database.Encrypt(key, answers.Secret)
	check(err)

	entry := &database.Entry{Username: username, Issuer: issuer, Secret: secret}
	_, err = db.NewInsert().Model(entry).Exec(ctx)
	if err != nil {
		return err
	}
	log.Println("Done!")

	var again bool
	err = survey.AskOne(&survey.Confirm{Message: "Add another entry"}, &again)
	checkSurvey(err)
	if again {
		goto add
	}
	return nil
}
