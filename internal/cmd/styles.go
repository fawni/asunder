package cmd

import (
	"github.com/charmbracelet/bubbles/list"
	lg "github.com/charmbracelet/lipgloss"
)

const (
	accent      = "#ccbcff"
	accentDark  = "#8977C2"
	accentLight = "#E9E2FF"
	secondary   = "#FCB6D0"
	text        = "#F8F8F0"
	danger      = "#FA827B"
)

var (
	styles   = list.NewDefaultItemStyles()
	delegate = list.NewDefaultDelegate()

	dimmedAccentColor = lg.AdaptiveColor{Light: accentLight, Dark: accentDark}
	accentColor       = lg.Color(accent)
	secondaryColor    = lg.Color(secondary)
	textColor         = lg.Color(text)
	dangerColor       = lg.Color(danger)

	appStyle = lg.NewStyle().Margin(1, 2)

	dangerForeground     = lg.NewStyle().Foreground(dangerColor)
	dangerForegroundBold = lg.NewStyle().Bold(true).Foreground(dangerColor)
	accentForegroundBold = lg.NewStyle().Bold(true).Foreground(accentColor)
)
