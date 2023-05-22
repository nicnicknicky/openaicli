package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func StartTUI() error {
	teaProgram = tea.NewProgram(InitMain(nil), tea.WithAltScreen())
	if _, err := teaProgram.Run(); err != nil {
		return err
	}

	return nil
}
