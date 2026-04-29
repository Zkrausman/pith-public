package config

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page int

const (
	parsersPage page = iota
	settingsPage
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Toggle key.Binding
	Tab    key.Binding
	Save   key.Binding
	Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Toggle, k.Tab, k.Save, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Toggle, k.Tab, k.Save, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("â†‘/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("â†“/j", "down"),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("space/enter", "toggle/edit"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch page"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	tabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(lipgloss.Color("#7D56F4"))

	activeTabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#FAFAFA")).
			Bold(true)

	sectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EE6FF8")).
			Bold(true)

	checkedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777")).
			Italic(true).
			MarginLeft(2)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1).
			Margin(1).
			Width(80)
)

type configModel struct {
	config    *Config
	parsers   []string
	cursor    int
	page      page
	inputs    []textinput.Model
	focusIdx  int
	quitting  bool
	saving    bool
	err       error
	help      help.Model
	keys      keyMap
}

func initialModel(cfg *Config, availableParsers []string) configModel {
	sort.Strings(availableParsers)

	// Initialize inputs for settings
	inputs := make([]textinput.Model, 5)
	
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "MaxLines"
	inputs[0].SetValue(strconv.Itoa(cfg.MaxLines))
	
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "HeadLines"
	inputs[1].SetValue(strconv.Itoa(cfg.HeadLines))
	
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "TailLines"
	inputs[2].SetValue(strconv.Itoa(cfg.TailLines))
	
	inputs[3] = textinput.New()
	inputs[3].Placeholder = "USDRate"
	inputs[3].SetValue(fmt.Sprintf("%.2f", cfg.USDPerMillionTokens))
	
	inputs[4] = textinput.New()
	inputs[4].Placeholder = "Heuristic"
	inputs[4].SetValue(fmt.Sprintf("%.2f", cfg.TokenHeuristic))

	return configModel{
		config:  cfg,
		parsers: availableParsers,
		inputs:  inputs,
		help:    help.New(),
		keys:    keys,
	}
}

func (m configModel) Init() tea.Cmd {
	return nil
}

func (m configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.focusIdx > 0 {
			if key.Matches(msg, m.keys.Quit) || msg.Type == tea.KeyEnter {
				m.focusIdx = 0
				return m, nil
			}
			var cmd tea.Cmd
			m.inputs[m.focusIdx-1], cmd = m.inputs[m.focusIdx-1].Update(msg)
			return m, cmd
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keys.Tab):
			if m.page == parsersPage {
				m.page = settingsPage
				m.cursor = 0
			} else {
				m.page = parsersPage
				m.cursor = 0
			}

		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, m.keys.Down):
			limit := len(m.parsers)
			if m.page == settingsPage {
				limit = len(m.inputs)
			}
			if m.cursor < limit-1 {
				m.cursor++
			}

		case key.Matches(msg, m.keys.Toggle):
			if m.page == parsersPage {
				p := m.parsers[m.cursor]
				m.config.EnabledParsers[p] = !m.config.EnabledParsers[p]
			} else {
				m.focusIdx = m.cursor + 1
				return m, m.inputs[m.focusIdx-1].Focus()
			}

		case key.Matches(msg, m.keys.Save):
			m.saving = true
			if err := m.save(); err != nil {
				m.err = err
				m.saving = false
				return m, nil
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *configModel) save() error {
	maxLines, err := strconv.Atoi(m.inputs[0].Value())
	if err != nil {
		return fmt.Errorf("invalid MaxLines: %v", err)
	}
	headLines, err := strconv.Atoi(m.inputs[1].Value())
	if err != nil {
		return fmt.Errorf("invalid HeadLines: %v", err)
	}
	tailLines, err := strconv.Atoi(m.inputs[2].Value())
	if err != nil {
		return fmt.Errorf("invalid TailLines: %v", err)
	}
	usdRate, err := strconv.ParseFloat(m.inputs[3].Value(), 64)
	if err != nil {
		return fmt.Errorf("invalid USDRate: %v", err)
	}
	heuristic, err := strconv.ParseFloat(m.inputs[4].Value(), 64)
	if err != nil {
		return fmt.Errorf("invalid Heuristic: %v", err)
	}

	m.config.MaxLines = maxLines
	m.config.HeadLines = headLines
	m.config.TailLines = tailLines
	m.config.USDPerMillionTokens = usdRate
	m.config.TokenHeuristic = heuristic

	return m.config.Save()
}

func (m configModel) View() string {
	if m.quitting {
		return ""
	}
	if m.saving {
		return successStyle.Render("Configuration saved successfully!") + "\n"
	}

	var content strings.Builder

	content.WriteString(titleStyle.Render(" Pith Configuration "))
	content.WriteString("\n\n")

	// Tabs
	parsersTab := tabStyle.Render("Parsers")
	if m.page == parsersPage {
		parsersTab = activeTabStyle.Render("Parsers")
	}
	settingsTab := tabStyle.Render("Settings")
	if m.page == settingsPage {
		settingsTab = activeTabStyle.Render("Settings")
	}

	content.WriteString(parsersTab + settingsTab)
	content.WriteString("\n\n")

	if m.page == parsersPage {
		content.WriteString(sectionStyle.Render("Enable/Disable specialized parsers:"))
		content.WriteString("\n")
		for i, p := range m.parsers {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			checked := "[ ]"
			if m.config.EnabledParsers[p] {
				checked = checkedStyle.Render("[x]")
			}
			line := fmt.Sprintf("%s %s %s", cursor, checked, p)
			if m.cursor == i {
				content.WriteString(selectedStyle.Render(line))
			} else {
				content.WriteString(line)
			}
			content.WriteString("\n")
		}
	} else {
		content.WriteString(sectionStyle.Render("Configure Pith output management:"))
		content.WriteString("\n")
		settingsNames := []string{"MaxLines", "HeadLines", "TailLines", "USD/M Tokens", "Heuristic"}
		settingsDescs := []string{
			"Threshold for 'Middle-Out' truncation. Snipping occurs if output exceeds this.",
			"Number of lines to preserve at the very top of the output.",
			"Number of lines to preserve at the very bottom (critical for errors).",
			"The cost rate used to calculate savings analytics in the dashboard.",
			"Average characters per token (heuristic used for estimation).",
		}
		for i, name := range settingsNames {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			val := m.inputs[i].Value()
			if m.focusIdx == i+1 {
				val = m.inputs[i].View()
			}
			line := fmt.Sprintf("%s %-15s: %s", cursor, name, val)
			if m.cursor == i {
				content.WriteString(selectedStyle.Render(line))
				content.WriteString("\n")
				content.WriteString(descStyle.Render(settingsDescs[i]))
			} else {
				content.WriteString(line)
			}
			content.WriteString("\n")
		}
	}

	if m.err != nil {
		content.WriteString("\n" + errorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n")
	}

	content.WriteString("\n")
	content.WriteString(m.help.View(m.keys))

	return containerStyle.Render(content.String())
}

func (c *Config) InteractiveConfig(availableParsers []string) error {
	// Sync available parsers into map (enable new ones by default)
	for _, p := range availableParsers {
		if _, ok := c.EnabledParsers[p]; !ok {
			c.EnabledParsers[p] = true
		}
	}

	p := tea.NewProgram(initialModel(c, availableParsers))
	_, err := p.Run()
	return err
}
