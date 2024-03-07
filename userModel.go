package main

import (
	"fmt"

	"gh-bubrls/messages"
	"gh-bubrls/structs"
	"gh-bubrls/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh/v2/pkg/api"
)

type UserModel struct {
	UnAuthenticated bool
	User            structs.User
	SelectedOrgUrl  string
	list            list.Model
	loaded          bool
	width           int
	height          int
}

func NewUserModel() UserModel {
	return UserModel{
		list: list.New(
			[]list.Item{},
			list.NewDefaultDelegate(),
			0,
			0,
		),
	}
}

func (m UserModel) Init() tea.Cmd {
	m.User
	return getOrganisations
}

func (m UserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		if !m.loaded {
			m.list.SetWidth(m.width)
			m.list.SetHeight(m.height)
			m.loaded = true
		}
		return m, nil

	case messages.AuthenticationErrorMsg:
		m.UnAuthenticated = true
		return m, nil

	case messages.OrgListMsg:
		m.list = buildOrgListModel(msg.Organisations, m.width, m.height, m.User)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
			// case "enter", " ":
			// 	MainModel[consts.UserModelName] = m
			// 	item := m.list.SelectedItem()
			// 	orgModel := NewOrgModel(item.(structs.ListItem).Title(), m.width, m.height)

			// 	MainModel[consts.OrganisationModelName] = orgModel

			// 	return orgModel, orgModel.GetRepositories
		}
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m UserModel) View() string {
	if m.UnAuthenticated {
		return fmt.Sprintln("You are not authenticated try running `gh auth login`. Press q to quit.")
	}

	return style.App.Render(m.list.View())
}

func buildOrgListModel(organisations []structs.Organisation, width, height int, user structs.User) list.Model {
	items := make([]list.Item, len(organisations))
	for i, org := range organisations {
		items[i] = structs.NewListItem(org.Login, org.Url)
	}

	list := list.New(items, style.DefaultDelegate, width, height-2)

	list.Title = "User: " + user.Name
	list.SetStatusBarItemName("Organisation", "Organisations")
	list.Styles.Title = style.Title
	list.SetShowTitle(true)

	return list
}

func getOrganisations() tea.Msg {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return messages.AuthenticationErrorMsg{Err: err}
	}
	response := []structs.Organisation{}

	err = client.Get("user/orgs", &response)
	if err != nil {
		fmt.Println(err)
		return messages.ErrMsg{Err: err}
	}

	return messages.OrgListMsg{Organisations: response}
}
