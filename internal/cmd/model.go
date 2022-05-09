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

type Item struct {
	ID       int
	Username string
	Issuer   string
	Code     string
}

func (i Item) Title() string { return i.Code }
func (i Item) Description() string {
	return fmt.Sprintf("%s (%s)", strings.Title(i.Issuer), i.Username)
}
func (i Item) FilterValue() string {
	if i.Code != common.InvalidCode {
		return i.Issuer + i.Username
	}
	return ""
}

type Model struct {
	List      list.Model
	Timer     timer.Model
	Countdown string
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.Timer.Init(), tea.HideCursor)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.Timer, cmd = m.Timer.Update(msg)
		return m, cmd
	case timer.TimeoutMsg:
		m.Timer.Timeout = common.TTL
	case tea.KeyMsg:
		var cmd tea.Cmd
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			item := m.List.SelectedItem().(Item)
			if item.Code != common.InvalidCode {
				err := clipboard.WriteAll(item.Code)
				if err != nil {
					cmd = m.List.NewStatusMessage(common.DangerForegroundBold.Render(fmt.Sprintf("clipboard: %s", err.Error())))
					return m, cmd
				}
				cmd = m.List.NewStatusMessage(fmt.Sprintf("Copied %s to clipboard!", termenv.String(item.Code).Bold()))
				return m, cmd
			}
		}
	case tea.WindowSizeMsg:
		v, h, _, _ := common.AppStyle.GetMargin()
		h *= 4
		v *= 4
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	switch countdown() {
	case common.TTL:
		items, err := getItems()
		if err != nil {
			cmd := m.List.NewStatusMessage(common.DangerForegroundBold.Render(err.Error()))
			return m, cmd
		}
		m.List.SetItems(items)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	m.Countdown = m.renderStatus()
	return m.Countdown + "\n" + m.List.View()
}

func countdown() int {
	i := time.Now().Second()
	if i >= common.TTL {
		i -= common.TTL
	}
	i = common.TTL - i
	return i
}

func (m *Model) renderStatus() string {
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
	entries, err := database.GetEntries(DB, Key)
	if err != nil {
		return []list.Item{}, err
	}
	items := make([]list.Item, 0, len(entries))
	for _, entry := range entries {
		code, err := totp.GenerateCode(entry.Secret, time.Now())
		if err != nil {
			code = common.InvalidCode
		}
		items = append(items, Item{ID: entry.ID, Code: code, Username: entry.Username, Issuer: entry.Issuer})
	}
	return items, nil
}
