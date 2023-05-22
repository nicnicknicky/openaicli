package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type MainModel struct {
	// Loading
	spinner    spinner.Model
	timer      timer.Model
	loadingMsg string

	// Menu - Main and Chats
	menu list.Model

	// Chats - Msg and Completion
	msgInput     textarea.Model
	cmplView     viewport.Model
	showMsgInput bool
	showChat     *int

	// Misc.
	styles *Styles
	err    error
}

func InitMain(menuType *menuType) tea.Model {
	menuTitle := mainMenuTitle
	menuItems := menuOptions
	if menuType != nil && *menuType == selectChatHistory {
		menuTitle = chatHistoryMenuTitle
		menuItems = GetChatMsgMenuItems()
	}

	m := MainModel{
		styles:  DefaultStyles(),
		spinner: newSpinner(),
		timer:   newTimer(),
		menu:    list.NewModel(menuItems, list.NewDefaultDelegate(), 0, 0),
	}

	h, v := m.styles.Menu.GetFrameSize()
	m.menu.SetSize(windowSize.Width-h, windowSize.Height-v)
	m.menu.Title = menuTitle

	return m
}

func (m MainModel) Init() tea.Cmd {
	return tea.Batch(spinner.Tick, tea.ClearScreen)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		windowSize = msg
		// TODO: remove duplicate in InitMain()
		h, v := m.styles.Menu.GetFrameSize()
		m.menu.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			mo := m.menu.SelectedItem().(menuOption)
			switch mt := mo.menuType; mt {
			case askChatGPT:
				m.UpdateMsgInput()
			case selectChatHistory:
				return InitMain(&mo.menuType), tea.ClearScreen
			case viewChatHistory:
				if err := m.newCmplView(m.menu.Index()); err != nil {
					m.err = err
					return m, nil
				}
			}
		case tea.KeyCtrlS:
			content := m.msgInput.Value()
			if content == "" {
				return m, nil
			}
			m.loadingMsg = "Waiting for ChatGPT's response..."
			return m, tea.Batch(spinner.Tick, m.timer.Init(), sendChatGPTCmd(content))
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			return InitMain(nil), tea.ClearScreen
		}
	case chatTeaMsg:
		chatID := int(msg)
		if err := m.newCmplView(chatID); err != nil {
			m.err = err
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)
	case errTeaMsg:
		m.err = msg
		return m, nil
	}

	m.menu, cmd = m.menu.Update(msg)
	cmds = append(cmds, cmd)
	m.msgInput, cmd = m.msgInput.Update(msg)
	cmds = append(cmds, cmd)
	m.cmplView, cmd = m.cmplView.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	if windowSize.Width == 0 || windowSize.Height == 0 {
		return fmt.Sprintf("%s Loading...\n", m.spinner.View())
	}

	if m.loadingMsg != "" {
		return lipgloss.Place(
			windowSize.Width,
			windowSize.Height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Left,
				fmt.Sprintf("%s %s\n", m.spinner.View(), m.loadingMsg),
				fmt.Sprintf("%s %s", "Timeout in", m.timer.View()),
				fmt.Sprintf("%s", helperString),
			),
		)
	}

	if m.showMsgInput {
		return lipgloss.Place(
			windowSize.Width,
			windowSize.Height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Left,
				"Start a conversation with ChatGPT...",
				m.styles.InputField.Render(m.msgInput.View()),
				"( Ctrl + S to Send | Ctrl + C to Quit )",
			),
		)
	}

	if m.showChat != nil {
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.showChat), m.cmplView.View(), m.footerView())
	}

	return lipgloss.Place(
		windowSize.Width,
		windowSize.Height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.Menu.Render(m.menu.View()),
		),
	)
}

// === MESSAGES ( textarea ) ===

func (m *MainModel) UpdateMsgInput() {
	m.msgInput = newMsgInput()
	m.showMsgInput = true
}

// === COMPLETIONS ( viewport ) ===

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m MainModel) headerView(showChat *int) string {
	var titleStr string
	if showChat != nil {
		msg := store[*showChat].msg
		maxCharLen := int(float64(m.cmplView.Width) * 0.8)
		if len(msg) > maxCharLen {
			maxCharLen = len(msg)
			msg = fmt.Sprintf("%s...", msg[:maxCharLen])
		}
		titleStr = msg
	}
	title := headerStyle.Render(titleStr)
	line := strings.Repeat("─", max(0, m.cmplView.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m MainModel) footerView() string {
	percent := footerStyle.Render(fmt.Sprintf("%3.f%%", m.cmplView.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.cmplView.Width-lipgloss.Width(percent)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, percent)
}

func (m *MainModel) newCmplView(chatID int) error {
	headerHeight := lipgloss.Height(m.headerView(&chatID))
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight

	m.cmplView = viewport.New(windowSize.Width, windowSize.Height-verticalMarginHeight)
	m.cmplView.YPosition = headerHeight + 1
	m.cmplView.HighPerformanceRendering = useHighPerformanceRenderer

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(windowSize.Width),
	)
	if err != nil {
		return err
	}

	cmpl := store[chatID].cmpl
	cmpl += "\n\n---"
	cmpl += fmt.Sprintf("\n\n# %s", helperString)
	rCmpl, err := renderer.Render(cmpl)
	if err != nil {
		return err
	}

	m.cmplView.SetContent(rCmpl)
	m.loadingMsg = ""
	m.showMsgInput = false
	m.showChat = &chatID

	return nil
}
