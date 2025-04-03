package cmdutils

import (
	"fmt"
	"main/pkg/model"
	"main/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type MenuModel struct {
	Choices []string
	Cursor  int
	Result  string
	Err     error
	Input   string
	History []string
}

func NewMenuModel() MenuModel {
	return MenuModel{
		Choices: []string{
			"OS and OSQuery Info",
			"Applications",
			"Run All Queries",
			"Quit",
		},
		Input:   "",
		History: make([]string, 0),
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) ExecuteChoice(choice int) (MenuModel, tea.Cmd) {
	if choice < 0 || choice >= len(m.Choices) {
		m.Err = fmt.Errorf("invalid choice: %d", choice)
		m.Result = ""
		return m, nil
	}

	if choice == len(m.Choices)-1 {
		return m, tea.Quit
	}

	if choice == 2 {
		queryTypes := []model.QueryType{
			model.QueryTypeOSAndOSQuery,
			model.QueryTypeApplications,
		}

		var results []string
		var errors []error

		for _, queryType := range queryTypes {
			output, err := ExecuteQuery(queryType)
			if err != nil {
				errors = append(errors, fmt.Errorf("error executing %s: %v", queryType, err))
				continue
			}
			results = append(results, string(output))
		}

		if len(errors) > 0 {
			errorMsgs := make([]string, len(errors))
			for i, err := range errors {
				errorMsgs[i] = err.Error()
			}
			m.Err = fmt.Errorf("errors occurred:\n%s", strings.Join(errorMsgs, "\n"))
		}
		if len(results) > 0 {
			m.Result = strings.Join(results, "\n")
		}
		m.History = append(m.History, fmt.Sprintf("%d", choice+1))
		m.Input = ""
		return m, nil
	}

	var queryType model.QueryType
	switch choice {
	case 0:
		queryType = model.QueryTypeOSAndOSQuery
	case 1:
		queryType = model.QueryTypeApplications
	}

	output, err := ExecuteQuery(queryType)
	if err != nil {
		m.Err = err
		m.Result = ""
		return m, nil
	}

	m.History = append(m.History, fmt.Sprintf("%d", choice+1))
	m.Result = string(output)
	m.Err = nil
	m.Input = ""
	return m, nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "1", "2", "3", "4":
			m.Input = msg.String()
			m.Cursor = int(msg.String()[0] - '1')
			return m, nil

		case "enter":
			if m.Input != "" {
				choice := int(m.Input[0] - '1')
				return m.ExecuteChoice(choice)
			}

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
				m.Input = fmt.Sprintf("%d", m.Cursor+1)
			}

		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
				m.Input = fmt.Sprintf("%d", m.Cursor+1)
			}

		case "backspace":
			m.Input = ""
			return m, nil
		}
	}

	return m, nil
}

func (m MenuModel) View() string {
	s := styles.TitleStyle.Render("OSQuery Terminal") + "\n\n"
	s += styles.TitleStyle.Render("Help") + "\n\n"

	for i, choice := range m.Choices {
		number := styles.CommandStyle.Render(fmt.Sprintf("%d", i+1))
		desc := styles.DescriptionStyle.Render(fmt.Sprintf(": %s", choice))
		s += fmt.Sprintf("  %s %s\n", number, desc)
	}

	s += "\n"
	s += styles.CommandStyle.Render("ctrl+c") + styles.DescriptionStyle.Render(" : exit") + "\n"
	s += styles.CommandStyle.Render("↑/↓") + styles.DescriptionStyle.Render("   : navigate options") + "\n\n"

	if m.Err != nil {
		s += styles.ErrorStyle.Render(fmt.Sprintf("Error: %v\n\n", m.Err))
	}

	if m.Result != "" {
		s += styles.ResultStyle.Render(m.Result) + "\n\n"
	}

	s += styles.PromptStyle.Render("🚀 > ")
	if m.Input != "" {
		s += styles.InputStyle.Render(m.Input)
		choice := int(m.Input[0] - '1')
		if choice >= 0 && choice < len(m.Choices) {
			s += styles.PlaceholderStyle.Render(fmt.Sprintf(" (Execute %s...)", m.Choices[choice]))
		}
	} else {
		s += styles.PlaceholderStyle.Render("Execute query...")
	}
	s += "\n"

	return s
}
