package cmd

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/x6r/asunder/internal/common"
	"github.com/x6r/asunder/internal/database"
)

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add an entry",
		Long:  `The add command adds entries to the database`,
		PreRun: func(cmd *cobra.Command, args []string) {
			connectDB()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
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
			Prompt:   &survey.Input{Message: "Enter your username »"},
			Validate: survey.Required,
		},
		{
			Name:     "issuer",
			Prompt:   &survey.Input{Message: "Enter the issuer name »"},
			Validate: survey.Required,
		},
		{
			Name:     "secret",
			Prompt:   &survey.Password{Message: "Enter TOTP secret »"},
			Validate: survey.Required,
		},
	}

add:
	var answers database.Entry
	err := survey.Ask(qs, &answers)
	common.CheckSurvey(err)
	username, err := database.Encrypt(Key, answers.Username)
	common.Check(err)
	issuer, err := database.Encrypt(Key, answers.Issuer)
	common.Check(err)
	secret, err := database.Encrypt(Key, answers.Secret)
	common.Check(err)

	ctx := context.Background()
	entry := &database.Entry{Username: username, Issuer: issuer, Secret: secret}
	_, err = DB.NewInsert().Model(entry).Exec(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Done!")

	var again bool
	err = survey.AskOne(&survey.Confirm{Message: "Add another entry"}, &again)
	common.CheckSurvey(err)
	if again {
		goto add
	}
	return nil
}
