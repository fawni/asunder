package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/fawni/asunder/internal/common"
	"github.com/fawni/asunder/internal/database"
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
	entries, err := database.GetEntries(DB, Key)
	common.Check(err)
	if len(entries) < 1 {
		fmt.Println("No entries found.")
		os.Exit(1)
	}

	var options []huh.Option[int]
	for _, entry := range entries {
		options = append(options, huh.NewOption(fmt.Sprintf("%s, %s", entry.Issuer, entry.Username), entry.ID))
	}

	var id int
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select an entry to delete").
				Options(options...).
				Value(&id),
		)).WithTheme(huh.ThemeCatppuccin()).Run()
	common.Check(err)

	var sure bool
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Are you sure you want to delete \"%s\"?", getDescriptionById(id))).
				Value(&sure),
		),
	).WithTheme(huh.ThemeCatppuccin()).Run()
	common.Check(err)

	if sure {
		ctx := context.Background()
		entry := &database.Entry{ID: id}
		_, err := DB.NewDelete().Model(entry).WherePK().Exec(ctx)
		common.Check(err)
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

	return fmt.Sprintf("%s, %s", issuer, username)
}
