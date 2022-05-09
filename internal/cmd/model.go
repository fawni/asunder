package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"github.com/pquerna/otp/totp"
	"github.com/x6r/asunder/internal/common"
	"github.com/x6r/asunder/internal/database"
)

type item struct {
	id       int
	username string
	issuer   string
	code     string
}

func (i item) Title() string { return i.code }
func (i item) Description() string {
	return fmt.Sprintf("%s (%s)", strings.Title(i.issuer), i.username)
}
func (i item) FilterValue() string {
	if i.code != invalidCode {
		return i.issuer + i.username
	}
	return ""
}

type model struct {
	list      list.Model
	timer     timer.Model
	countdown string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.timer.Init(), tea.HideCursor)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case timer.TimeoutMsg:
		m.timer.Timeout = ttl
	case tea.KeyMsg:
		var cmd tea.Cmd
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			item := m.list.SelectedItem().(item)
			if item.code != invalidCode {
				err := clipboard.WriteAll(item.code)
				if err != nil {
					cmd = m.list.NewStatusMessage(common.DangerForegroundBold.Render(fmt.Sprintf("clipboard: %s", err.Error())))
					return m, cmd
				}
				cmd = m.list.NewStatusMessage(fmt.Sprintf("Copied %s to clipboard!", termenv.String(item.code).Bold()))
				return m, cmd
			}
		}
	case tea.WindowSizeMsg:
		v, h, _, _ := common.AppStyle.GetMargin()
		h *= 4
		v *= 4
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	switch countdown() {
	case ttl:
		items, err := getItems()
		if err != nil {
			cmd := m.list.NewStatusMessage(common.DangerForegroundBold.Render(err.Error()))
			return m, cmd
		}
		m.list.SetItems(items)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	m.countdown = m.renderStatus()
	return m.countdown + "\n" + m.list.View()
}

func countdown() int {
	i := time.Now().Second()
	if i >= ttl {
		i -= ttl
	}
	i = ttl - i
	return i
}

func (m *model) renderStatus() string {
	sec := countdown()
	var ttl, status string
	if sec <= 7 {
		ttl = common.DangerForegroundBold.Render(strconv.Itoa(sec) + "s")
	} else {
		ttl = common.AccentForegroundBold.Render(strconv.Itoa(sec) + "s")
	}
	status = common.AppStyle.Render("Expiration: " + ttl)
	return status
}

func getItems() ([]list.Item, error) {
	entries, err := database.GetEntries(db, key)
	if err != nil {
		return []list.Item{}, err
	}
	items := make([]list.Item, 0, len(entries))
	for _, entry := range entries {
		code, err := totp.GenerateCode(entry.Secret, time.Now())
		if err != nil {
			code = invalidCode
		}
		items = append(items, item{id: entry.ID, code: code, username: entry.Username, issuer: entry.Issuer})
	}
	return items, nil
}
