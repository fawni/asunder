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
	"github.com/muesli/coral"
	"github.com/muesli/termenv"
	"github.com/x6r/asunder/internal/common"
	"github.com/x6r/asunder/internal/database"
)

type keymap struct {
	Enter teakey.Binding
}

var (
	db  *database.DB
	key []byte

	rootCmd = &coral.Command{
		Use:   "asunder",
		Short: "asunder is a command-line totp manager",
		CompletionOptions: coral.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PreRun: func(cmd *coral.Command, args []string) {
			connectDB()
		},
		RunE: func(cmd *coral.Command, args []string) error {
			return startModel()
		},
	}
)

func init() {
	if !fileExists(common.PathDB) {
		err := initAsunder()
		check(err)
		os.Exit(0)
	}
}

func Execute() {
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
	if err := program.Start(); err != nil {
		return err
	}

	return nil
}

func connectDB() {
	promptPassword()
	var err error
	db, err = database.InitDB()
	check(err)
}

func promptPassword() {
	buf, err := os.ReadFile(common.PathData)
	check(err)
	var k struct{ Key string }
	err = json.Unmarshal(buf, &k)
	check(err)

	var password string
	err = survey.AskOne(&survey.Password{
		Message: "Enter master password ›",
	}, &password)
	checkSurvey(err)
	if database.Hash(password).Text != k.Key {
		fmt.Println("password does not match")
		os.Exit(1)
	}

	key = database.Hash(password).Hash
	termenv.ClearLine()
}

func initAsunder() error {
	var qs = []*survey.Question{
		{
			Name:     "password",
			Prompt:   &survey.Password{Message: "Enter a master password ›"},
			Validate: survey.Required,
		},
		{
			Name:     "repassword",
			Prompt:   &survey.Password{Message: "Re-Enter master password ›"},
			Validate: survey.Required,
		},
	}

	var answers struct {
		Pass   string `survey:"password"`
		Repass string `survey:"repassword"`
	}

	fmt.Println("First time setup! You will be prompted for the master password everytime you use asunder.")
	err := survey.Ask(qs, &answers)
	checkSurvey(err)
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

func setupModel() (model, error) {
	var keymap = &keymap{
		Enter: teakey.NewBinding(
			teakey.WithKeys("enter"),
			teakey.WithHelp("enter", "copy"),
		),
	}

	items, err := getItems()
	if err != nil {
		return model{}, err
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0), timer: timer.New(ttl)}
	m.list.Title = "asunder"
	termenv.SetWindowTitle(m.list.Title)
	m.list.StatusMessageLifetime = 1200 * time.Millisecond
	if len(items) > 0 {
		m.list.AdditionalShortHelpKeys = func() []teakey.Binding {
			return []teakey.Binding{
				keymap.Enter,
			}
		}
		m.list.AdditionalFullHelpKeys = m.list.AdditionalShortHelpKeys
	}
	common.Styles.SelectedTitle = common.Styles.SelectedTitle.Copy().Foreground(common.AccentColor).BorderForeground(common.DimmedAccentColor)
	common.Styles.SelectedDesc = common.Styles.SelectedTitle.Copy().Foreground(common.DimmedAccentColor)
	common.Delegate.Styles = common.Styles
	m.list.SetDelegate(common.Delegate)

	m.list.Styles.Title = m.list.Styles.Title.Copy().Background(common.SecondaryColor).Foreground(common.TextColor).Bold(true)
	m.list.Styles.FilterCursor = m.list.Styles.FilterCursor.Copy().Foreground(common.AccentColor)

	return m, nil
}
