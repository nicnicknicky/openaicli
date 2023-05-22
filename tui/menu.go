package tui

import "github.com/charmbracelet/bubbles/list"

type menuType int

type menuOption struct {
	title, desc string
	menuType    menuType
}

func (mo menuOption) Title() string       { return mo.title }
func (mo menuOption) Description() string { return mo.desc }
func (mo menuOption) FilterValue() string { return mo.title }

func GetChatMsgMenuItems() []list.Item {
	if len(store) == 0 {
		return nil
	}

	var menuOptions []list.Item
	for _, chat := range store {
		menuOptions = append(menuOptions, menuOption{
			title:    chat.msg,
			menuType: viewChatHistory,
		})
	}

	return menuOptions
}
