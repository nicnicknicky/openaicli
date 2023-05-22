package tui

import (
	"playground/openaicli/chatgpt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	teaProgram *tea.Program
	windowSize tea.WindowSizeMsg
	OAIClient  *chatgpt.OpenaiClient
	store      = make(map[int]chat)

	menuOptions = []list.Item{
		menuOption{
			title:    "[ Send ChatGPT A Message ]",
			desc:     "Queries the ChatGPT API with a message and receive a completion.",
			menuType: askChatGPT,
		},
		menuOption{
			title:    "[ Chat History ]",
			desc:     "View historical queries and their corresponding completions.",
			menuType: selectChatHistory,
		},
	}

	// Styling
	headerStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	footerStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return headerStyle.Copy().BorderStyle(b)
	}()
)
