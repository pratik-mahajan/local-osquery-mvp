package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor     = lipgloss.Color("#FF6B6B")
	textColor        = lipgloss.Color("#FFFFFF")
	promptColor      = lipgloss.Color("#87CEEB")
	placeholderColor = lipgloss.Color("#666666")
	inputColor       = lipgloss.Color("#FFFFFF")
	resultColor      = lipgloss.Color("#00FF00")

	TitleStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Bold(true)

	CommandStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	DescriptionStyle = lipgloss.NewStyle().
				Foreground(textColor)

	PromptStyle = lipgloss.NewStyle().
			Foreground(promptColor).
			Bold(true)

	PlaceholderStyle = lipgloss.NewStyle().
				Foreground(placeholderColor).
				Italic(true)

	InputStyle = lipgloss.NewStyle().
			Foreground(inputColor)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	ResultStyle = lipgloss.NewStyle().
			Foreground(resultColor)
)
