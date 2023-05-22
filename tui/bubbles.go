package tui

import (
	"playground/openaicli/chatgpt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// === GENERAL ===

type errTeaMsg struct{ err error }

func (e errTeaMsg) Error() string { return e.err.Error() }

// === STYLES ===

type Styles struct {
	InputField, Menu lipgloss.Style
}

func DefaultStyles() *Styles {
	return &Styles{
		InputField: lipgloss.NewStyle().BorderForeground(ansiCyan).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80),
		Menu:       lipgloss.NewStyle().Margin(20, 1),
	}
}

// === LOADING ===

func newSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(ansiPink)
	return s
}

func newTimer() timer.Model {
	return timer.NewWithInterval(chatgpt.ChatGPTTimeout, time.Millisecond)
}

// === USER INPUT ===

func newMsgInput() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Enter Your Message..."
	ta.ShowLineNumbers = false
	ta.Focus()
	return ta
}

// === OUTPUT ===

type chat struct{ msg, cmpl string }

type chatTeaMsg int

func sendChatGPTCmd(msg string) tea.Cmd {
	return func() tea.Msg {
		trimMsg := strings.TrimSpace(msg)
		cmpl, err := OAIClient.SendChatGPT(trimMsg)
		if err != nil {
			return errTeaMsg{err}
		}

		newChatID := len(store)
		store[newChatID] = chat{msg: trimMsg, cmpl: cmpl}
		return chatTeaMsg(newChatID)
	}
}
