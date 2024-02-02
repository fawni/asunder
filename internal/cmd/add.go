package cmd

import (
	"context"

	"github.com/charmbracelet/huh"
	"github.com/fawni/asunder/internal/common"
	"github.com/fawni/asunder/internal/database"
	"github.com/spf13/cobra"
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

	username string
	issuer   string
	secret   string
)

func init() {
	rootCmd.AddCommand(addCmd)
}

func addEntry() error {
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter your username").Prompt("? ").Value(&username),
			huh.NewInput().Title("Enter your issuer name").Prompt("? ").Value(&issuer),
			huh.NewInput().Title("Enter your TOTP secret").Prompt("? ").Value(&secret).Password(true),
		),
	).WithTheme(huh.ThemeCatppuccin()).Run()

	if err != nil {
		return err
	}

	username, err := database.Encrypt(Key, username)
	common.Check(err)
	issuer, err := database.Encrypt(Key, issuer)
	common.Check(err)
	secret, err := database.Encrypt(Key, secret)
	common.Check(err)

	ctx := context.Background()
	entry := &database.Entry{Username: username, Issuer: issuer, Secret: secret}
	_, err = DB.NewInsert().Model(entry).Exec(ctx)
	if err != nil {
		return err
	}

	var again bool
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Would you like to add another entry?").
				Value(&again)),
	).WithTheme(huh.ThemeCatppuccin()).Run()

	if err != nil {
		return err
	}

	if again {
		username, issuer, secret = "", "", ""
		return addEntry()
	}

	return nil
}
