package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	teakey "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fawni/asunder/internal/common"
	"github.com/fawni/asunder/internal/database"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

type keymap struct {
	Enter teakey.Binding
}

var (
	DB  *database.DB
	Key []byte

	rootCmd = &cobra.Command{
		Use:   "asunder",
		Short: "asunder is a command-line totp manager",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			connectDB()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return startModel()
		},
	}
)

func init() {
	if !common.FileExists(common.PathDB) {
		err := initAsunder()
		common.Check(err)
		os.Exit(0)
	}
}

func Execute() {
	output := termenv.NewOutput(os.Stdout)
	defer output.Reset()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startModel() error {
	m, err := setupModel()
	if err != nil {
		return err
	}

	program := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}

func connectDB() {
	promptPassword()
	var err error
	DB, err = database.InitDB()
	common.Check(err)
}

func promptPassword() {
	buf, err := os.ReadFile(common.PathData)
	common.Check(err)
	var k struct{ Key string }
	err = json.Unmarshal(buf, &k)
	common.Check(err)

	var password string
	err = survey.AskOne(&survey.Password{
		Message: "Enter master password »",
	}, &password)
	common.CheckSurvey(err)
	if database.Hash(password).Text != k.Key {
		fmt.Println("password does not match")
		os.Exit(1)
	}

	Key = database.Hash(password).Hash
}

func initAsunder() error {
	var qs = []*survey.Question{
		{
			Name:     "password",
			Prompt:   &survey.Password{Message: "Enter a master password »"},
			Validate: survey.Required,
		},
		{
			Name:     "repassword",
			Prompt:   &survey.Password{Message: "Re-Enter master password »"},
			Validate: survey.Required,
		},
	}

	var answers struct {
		Pass   string `survey:"password"`
		Repass string `survey:"repassword"`
	}

	fmt.Println("First time setup! You will be prompted for the master password everytime you use asunder.")
	err := survey.Ask(qs, &answers)
	common.CheckSurvey(err)
	if answers.Pass != answers.Repass {
		fmt.Println("password does not match")
		os.Exit(1)
	}

	secret := answers.Pass
	if _, err := database.InitDB(); err != nil {
		return err
	}
	if err := os.WriteFile(common.PathData, []byte(fmt.Sprintf(`{"key": "%s"}`, database.Hash(secret).Text)), 0644); err != nil {
		return err
	}

	fmt.Printf("Database created! use %s to start adding entries to the database.", termenv.String("asunder add").Bold())
	return nil
}

func setupModel() (Model, error) {
	var keymap = &keymap{
		Enter: teakey.NewBinding(
			teakey.WithKeys("enter"),
			teakey.WithHelp("enter", "copy"),
		),
	}

	items, err := getItems()
	if err != nil {
		return Model{}, err
	}

	// this is a mess and i'm not sure how to make it any better

	m := Model{List: list.New(items, list.NewDefaultDelegate(), 0, 0), Timer: timer.New(common.TTL)}
	m.List.Title = "asunder"
	termenv.NewOutput(os.Stdout).SetWindowTitle(m.List.Title)
	m.List.StatusMessageLifetime = 1200 * time.Millisecond
	if len(items) > 0 {
		m.List.AdditionalShortHelpKeys = func() []teakey.Binding {
			return []teakey.Binding{
				keymap.Enter,
			}
		}
		m.List.AdditionalFullHelpKeys = m.List.AdditionalShortHelpKeys
	}
	common.Styles.SelectedTitle = common.Styles.SelectedTitle.Copy().Foreground(common.AccentColor).BorderForeground(common.DimmedAccentColor)
	common.Styles.SelectedDesc = common.Styles.SelectedTitle.Copy().Foreground(common.DimmedAccentColor)
	common.Delegate.Styles = common.Styles
	m.List.SetDelegate(common.Delegate)

	m.List.Styles.Title = m.List.Styles.Title.Copy().Background(common.SecondaryColor).Foreground(common.TextColor).Bold(true)
	m.List.Styles.FilterCursor = m.List.Styles.FilterCursor.Copy().Foreground(common.AccentColor)

	return m, nil
}
