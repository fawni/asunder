package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	teakey "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/coral"
	"github.com/muesli/termenv"
	"github.com/uptrace/bun"
	"github.com/x6r/asunder/internal/config"
	"github.com/x6r/asunder/internal/database"
)

type keymap struct {
	Enter teakey.Binding
}

var (
	db  *bun.DB
	key []byte

	rootCmd = &coral.Command{
		Use:   "asunder",
		Short: "asunder is a command-line TOTP manager",
		RunE: func(cmd *coral.Command, args []string) error {
			return startModel()
		},
	}
)

func init() {
	log.SetFlags(0)
	if fileExists(config.PathDB) {
		coral.OnInitialize(connectDB)
	} else {
		err := initAsunder()
		check(err)
		os.Exit(0)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
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
	buf, err := os.ReadFile(config.PathData)
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
		log.Fatalln("password does not match")
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

	log.Println("First time setup! You will be prompted for the master password everytime you use asunder.")
	err := survey.Ask(qs, &answers)
	checkSurvey(err)
	if answers.Pass != answers.Repass {
		log.Fatalln("password does not match")
	}

	secret := answers.Pass
	if _, err := database.InitDB(); err != nil {
		return err
	}
	if err := os.WriteFile(config.PathData, []byte(fmt.Sprintf(`{"key": "%s"}`, database.Hash(secret).Text)), 0644); err != nil {
		return err
	}

	log.Printf("Database created! use %s to start adding entries to the database.", termenv.String("asunder add").Bold())
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
	styles.SelectedTitle = styles.SelectedTitle.Copy().Foreground(accentColor).BorderForeground(dimmedAccentColor)
	styles.SelectedDesc = styles.SelectedTitle.Copy().Foreground(dimmedAccentColor)
	delegate.Styles = styles
	m.list.SetDelegate(delegate)

	m.list.Styles.Title = m.list.Styles.Title.Copy().Background(secondaryColor).Foreground(textColor).Bold(true)
	m.list.Styles.FilterCursor = m.list.Styles.FilterCursor.Copy().Foreground(accentColor)

	return m, nil
}
