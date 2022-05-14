package common

import (
	"github.com/charmbracelet/bubbles/list"
	lg "github.com/charmbracelet/lipgloss"
)

const (
	Accent      = "#DDB6F2"
	AccentDark  = "#A689B5"
	AccentLight = "#E5C8F5"
	Secondary   = "#F28FAD"
	Text        = "#F8F8F0"
	Danger      = "#FA827B"
)

var (
	Styles   = list.NewDefaultItemStyles()
	Delegate = list.NewDefaultDelegate()

	DimmedAccentColor = lg.AdaptiveColor{Light: AccentLight, Dark: AccentDark}
	AccentColor       = lg.Color(Accent)
	SecondaryColor    = lg.Color(Secondary)
	TextColor         = lg.Color(Text)
	DangerColor       = lg.Color(Danger)

	AppStyle = lg.NewStyle().Margin(1, 2)

	DangerForeground     = lg.NewStyle().Foreground(DangerColor)
	DangerForegroundBold = lg.NewStyle().Bold(true).Foreground(DangerColor)
	AccentForegroundBold = lg.NewStyle().Bold(true).Foreground(AccentColor)
)
