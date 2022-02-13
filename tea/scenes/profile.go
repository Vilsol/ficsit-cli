package scenes

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	"github.com/satisfactorymodding/ficsit-cli/tea/components"
	"github.com/satisfactorymodding/ficsit-cli/tea/utils"
)

var _ tea.Model = (*profile)(nil)

type profile struct {
	root       components.RootModel
	list       list.Model
	parent     tea.Model
	profile    *cli.Profile
	hadRenamed bool
}

func NewEditProfile(root components.RootModel, parent tea.Model, profileData *cli.Profile) tea.Model {
	model := profile{
		root:    root,
		parent:  parent,
		profile: profileData,
	}

	items := []list.Item{
		utils.SimpleItem{
			ItemTitle: "Select",
			Activate: func(msg tea.Msg, currentModel tea.Model) (tea.Model, tea.Cmd) {
				if err := root.SetCurrentProfile(profileData); err != nil {
					panic(err) // TODO Handle Error
				}

				return currentModel.(profile).parent, nil
			},
		},
	}

	if profileData.Name != cli.DefaultProfileName {
		items = append(items,
			utils.SimpleItem{
				ItemTitle: "Rename",
				Activate: func(msg tea.Msg, currentModel tea.Model) (tea.Model, tea.Cmd) {
					newModel := NewRenameProfile(root, currentModel, profileData)
					return newModel, newModel.Init()
				},
			},
			utils.SimpleItem{
				ItemTitle: "Delete",
				Activate: func(msg tea.Msg, currentModel tea.Model) (tea.Model, tea.Cmd) {
					if err := root.GetGlobal().Profiles.DeleteProfile(profileData.Name); err != nil {
						panic(err) // TODO Handle Error
					}

					return currentModel.(profile).parent, updateProfileListCmd
				},
			},
		)
	}

	model.list = list.NewModel(items, utils.NewItemDelegate(), root.Size().Width, root.Size().Height-root.Height())
	model.list.SetShowStatusBar(false)
	model.list.SetFilteringEnabled(false)
	model.list.Title = fmt.Sprintf("Profile: %s", profileData.Name)
	model.list.Styles = utils.ListStyles
	model.list.SetSize(model.list.Width(), model.list.Height())
	model.list.StatusMessageLifetime = time.Second * 3
	model.list.DisableQuitKeybindings()

	return model
}

func (m profile) Init() tea.Cmd {
	return nil
}

func (m profile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case KeyControlC:
			return m, tea.Quit
		case "q":
			if m.parent != nil {
				m.parent.Update(m.root.Size())

				if m.hadRenamed {
					return m.parent, updateProfileNamesCmd
				}

				return m.parent, nil
			}
			return m, nil
		case KeyEnter:
			i, ok := m.list.SelectedItem().(utils.SimpleItem)
			if ok {
				if i.Activate != nil {
					newModel, cmd := i.Activate(msg, m)
					if newModel != nil || cmd != nil {
						if newModel == nil {
							newModel = m
						}
						return newModel, cmd
					}
					return m, nil
				}
			}
			return m, nil
		default:
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := lipgloss.NewStyle().Margin(2, 2).GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		m.root.SetSize(msg)
	case updateProfileNames:
		m.hadRenamed = true
		m.list.Title = fmt.Sprintf("Profile: %s", m.profile.Name)
	}

	return m, nil
}

func (m profile) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.root.View(), m.list.View())
}
