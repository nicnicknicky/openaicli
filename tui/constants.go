package tui

import "github.com/charmbracelet/lipgloss"

const (
	mainMenuTitle        = "ChatGPT CLI Client"
	chatHistoryMenuTitle = "Chat History"
	helperString         = "\U0001f44b ( Esc to Menu | Ctrl + C to Quit )"

	askChatGPT menuType = iota
	selectChatHistory
	viewChatHistory

	// Styling
	useHighPerformanceRenderer = false
	ansiCyan                   = lipgloss.Color("36")
	ansiPink                   = lipgloss.Color("205")
)
