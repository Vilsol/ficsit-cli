package mods

// cspell:disable

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/satisfactorymodding/ficsit-cli/ficsit"
	"github.com/satisfactorymodding/ficsit-cli/tea/components"
	"github.com/satisfactorymodding/ficsit-cli/tea/scenes/keys"
	"github.com/satisfactorymodding/ficsit-cli/tea/utils"
)

// cspell:enable

var _ tea.Model = (*modVersionMenu)(nil)

type modInfo struct {
	root           components.RootModel
	parent         tea.Model
	modData        chan ficsit.GetModMod
	modDataCache   ficsit.GetModMod
	modError       chan string
	error          *components.ErrorComponent
	help           help.Model
	keys           modInfoKeyMap
	viewport       viewport.Model
	spinner        spinner.Model
	ready          bool
	compatViewMode bool
}

type modInfoKeyMap struct {
	Up         key.Binding
	UpHalf     key.Binding
	UpPage     key.Binding
	Down       key.Binding
	DownHalf   key.Binding
	DownPage   key.Binding
	Help       key.Binding
	Back       key.Binding
	CompatInfo key.Binding
}

func (k modInfoKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Back, k.Up, k.Down, k.CompatInfo}
}

func (k modInfoKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.UpHalf, k.UpPage},
		{k.Down, k.DownHalf, k.DownPage},
		{k.CompatInfo},
		{k.Help, k.Back},
	}
}

func NewModInfo(root components.RootModel, parent tea.Model, mod utils.Mod) tea.Model {
	model := modInfo{
		root:     root,
		viewport: viewport.Model{},
		spinner:  spinner.New(),
		parent:   parent,
		modData:  make(chan ficsit.GetModMod),
		modError: make(chan string),
		ready:    false,
		help:     help.New(),
		keys: modInfoKeyMap{
			Up:         key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
			UpHalf:     key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "up half page")),
			UpPage:     key.NewBinding(key.WithKeys("pgup", "b"), key.WithHelp("pgup/b", "page up")),
			Down:       key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
			DownHalf:   key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "down half page")),
			DownPage:   key.NewBinding(key.WithKeys("pgdn", "f"), key.WithHelp("pgdn/f", "page down")),
			Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
			Back:       key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "back")),
			CompatInfo: key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "toggle compatibility info view")),
		},
	}

	model.spinner.Spinner = spinner.MiniDot
	model.help.Width = root.Size().Width

	go func() {
		fullMod, err := root.GetProvider().GetMod(context.TODO(), mod.Reference)
		if err != nil {
			model.modError <- err.Error()
			return
		}

		if fullMod == nil {
			model.modError <- "unknown error (mod is nil)"
			return
		}

		model.modData <- fullMod.Mod
	}()

	return model
}

func (m modInfo) Init() tea.Cmd {
	return tea.Batch(utils.Ticker(), m.spinner.Tick)
}

func (m modInfo) CalculateSizes(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	if m.viewport.Width == 0 {
		return m, nil
	}

	bottomPadding := 2
	if m.help.ShowAll {
		bottomPadding = 4
	}

	top, right, bottom, left := lipgloss.NewStyle().Margin(m.root.Height(), 3, bottomPadding).GetMargin()
	m.viewport.Width = msg.Width - left - right
	m.viewport.Height = msg.Height - top - bottom
	m.root.SetSize(msg)

	m.help.Width = m.viewport.Width

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m modInfo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case keys.KeyControlC:
			return m, tea.Quit
		case "q":
			if m.parent != nil {
				m.parent.Update(m.root.Size())
				return m.parent, nil
			}
			return m, tea.Quit
		case "?":
			m.help.ShowAll = !m.help.ShowAll
			return m.CalculateSizes(m.root.Size())
		case "i":
			m.compatViewMode = !m.compatViewMode
			m.viewport = m.newViewport()
			m.viewport.SetContent(m.renderModInfo())
			return m.CalculateSizes(m.root.Size())
		default:
			break
		}

		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		return m.CalculateSizes(msg)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case utils.TickMsg:
		select {
		case mod := <-m.modData:
			m.modDataCache = mod
			m.viewport = m.newViewport()
			m.viewport.SetContent(m.renderModInfo())
			break
		case err := <-m.modError:
			errorComponent, _ := components.NewErrorComponent(err, time.Second*5)
			m.error = errorComponent
			break
		default:
			// skip
			break
		}
		return m, utils.Ticker()
	}

	return m, nil
}

func (m modInfo) newViewport() viewport.Model {
	bottomPadding := 2
	if m.help.ShowAll {
		bottomPadding = 4
	}

	top, right, bottom, left := lipgloss.NewStyle().Margin(m.root.Height(), 3, bottomPadding).GetMargin()
	return viewport.Model{Width: m.root.Size().Width - left - right, Height: m.root.Size().Height - top - bottom}
}

func (m modInfo) renderModInfo() string {
	mod := m.modDataCache

	title := lipgloss.NewStyle().Padding(0, 2).Render(utils.TitleStyle.Render(mod.Name)) + "\n"
	title += lipgloss.NewStyle().Padding(0, 3).Render("("+string(mod.Mod_reference)+")") + "\n"

	sidebar := ""
	sidebar += utils.LabelStyle.Render("Views: ") + strconv.Itoa(mod.Views) + "\n"
	sidebar += utils.LabelStyle.Render("Downloads: ") + strconv.Itoa(mod.Downloads) + "\n"
	sidebar += "\n"
	sidebar += utils.LabelStyle.Render("EA  Compat: ") + m.renderCompatInfo(mod.Compatibility.EA.State) + "\n"
	sidebar += utils.LabelStyle.Render("EXP Compat: ") + m.renderCompatInfo(mod.Compatibility.EXP.State) + "\n"
	sidebar += "\n"
	sidebar += utils.LabelStyle.Render("Authors:") + "\n"

	converter := md.NewConverter("", true, nil)
	converter.AddRules(md.Rule{
		Filter: []string{"#text"},
		Replacement: func(content string, selection *goquery.Selection, options *md.Options) *string {
			text := selection.Text()
			return &text
		},
	})

	for _, author := range mod.Authors {
		sidebar += "\n"
		sidebar += utils.LabelStyle.Render(author.User.Username) + " - " + author.Role
	}

	description := ""
	if m.compatViewMode {
		a := ""
		a += "Compatibility information is maintained by the community." + "\n"
		a += "If you encounter issues with a mod, please report it on the Discord." + "\n"
		a += "Learn more about what compatibility states mean on ficsit.app" + "\n\n"

		description = m.renderDescriptionText(a, converter)

		description += "  " + utils.TitleStyle.Render("Early Access Branch Compatibility Note") + "\n"
		description += m.renderDescriptionText(mod.Compatibility.EA.Note, converter)
		description += "\n\n"
		description += "  " + utils.TitleStyle.Render("Experimental Branch Compatibility Note") + "\n"
		description += m.renderDescriptionText(mod.Compatibility.EXP.Note, converter)
	} else {
		description += m.renderDescriptionText(mod.Full_description, converter)
	}

	bottomPart := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, strings.TrimSpace(description))

	return lipgloss.JoinVertical(lipgloss.Left, title, bottomPart)
}

func (m modInfo) renderDescriptionText(text string, converter *md.Converter) string {
	text = strings.TrimSpace(text)
	if text == "" {
		text = "(No notes provided)"
	}

	markdownDescription, err := converter.ConvertString(text)
	if err != nil {
		slog.Error("failed to convert html to markdown", slog.Any("err", err))
		markdownDescription = text
	}

	description, err := glamour.Render(markdownDescription, "dark")
	if err != nil {
		slog.Error("failed to render markdown", slog.Any("err", err))
		description = text
	}

	return description
}

func (m modInfo) renderCompatInfo(state ficsit.CompatibilityState) string {
	stateText := string(state)
	switch state {
	case ficsit.CompatibilityStateWorks:
		return utils.CompatWorksStyle.Render(stateText)
	case ficsit.CompatibilityStateDamaged:
		return utils.CompatDamagedStyle.Render(stateText)
	case ficsit.CompatibilityStateBroken:
		return utils.CompatBrokenStyle.Render(stateText)
	default:
		return utils.CompatUntestedStyle.Render("Unknown")
	}
}

func (m modInfo) View() string {
	if m.error != nil {
		helpBar := lipgloss.NewStyle().Padding(1, 2).Render(m.help.View(m.keys))
		return lipgloss.JoinVertical(lipgloss.Left, m.root.View(), m.error.View(), m.viewport.View(), helpBar)
	}

	if m.viewport.Height == 0 {
		spinnerView := lipgloss.NewStyle().Padding(0, 2, 1).Render(m.spinner.View() + " Loading...")
		return lipgloss.JoinVertical(lipgloss.Left, m.root.View(), spinnerView)
	}

	helpBar := lipgloss.NewStyle().Padding(1, 2).Render(m.help.View(m.keys))
	return lipgloss.JoinVertical(lipgloss.Left, m.root.View(), m.viewport.View(), helpBar)
}
