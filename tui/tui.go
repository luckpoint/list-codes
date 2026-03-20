package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	root       *TreeNode
	visible    []*TreeNode
	cursor     int
	offset     int
	width      int
	height     int
	rootPath   string
	configPath string
	saved      bool
	err        error
	statusMsg  string
	showHelp   bool
}

func NewModel(rootPath, configPath string, noConfig bool, opts BuildTreeOpts) (Model, error) {
	root, err := BuildTree(rootPath, opts)
	if err != nil {
		return Model{}, err
	}

	SetInitialState(root, opts.IncludePatterns, opts.ExcludePatterns)

	// Try loading existing config
	if configPath != "" && !noConfig {
		cfg, loadErr := LoadConfig(configPath)
		if loadErr == nil {
			ApplyConfig(root, cfg)
		}
	}

	m := Model{
		root:       root,
		rootPath:   rootPath,
		configPath: configPath,
		width:      80,
		height:     24,
	}
	m.visible = FlattenVisible(root)
	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "?":
			m.showHelp = true

		case "up", "k", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
				m.ensureVisible()
			}

		case "down", "j", "ctrl+n":
			if m.cursor < len(m.visible)-1 {
				m.cursor++
				m.ensureVisible()
			}

		case " ":
			if m.cursor < len(m.visible) {
				node := m.visible[m.cursor]
				Toggle(node)
				m.visible = FlattenVisible(m.root)
			}

		case "enter", "l", "right":
			if m.cursor < len(m.visible) {
				node := m.visible[m.cursor]
				if node.IsDir {
					ToggleExpand(node)
					m.visible = FlattenVisible(m.root)
				}
			}

		case "h", "left":
			if m.cursor < len(m.visible) {
				node := m.visible[m.cursor]
				if node.IsDir && node.Expanded {
					ToggleExpand(node)
					m.visible = FlattenVisible(m.root)
				} else if node.Parent != nil {
					// Move cursor to parent
					for i, n := range m.visible {
						if n == node.Parent {
							m.cursor = i
							m.ensureVisible()
							break
						}
					}
				}
			}

		case "a":
			SetAllState(m.root, Checked)
			m.visible = FlattenVisible(m.root)

		case "n":
			SetAllState(m.root, Unchecked)
			m.visible = FlattenVisible(m.root)

		case "s", "w":
			includes, excludes := GeneratePatterns(m.root)
			cfg := &Config{
				Include: includes,
				Exclude: excludes,
			}
			if err := SaveConfig(m.configPath, cfg); err != nil {
				m.statusMsg = fmt.Sprintf("Save error: %v", err)
				return m, nil
			}
			m.saved = true
			return m, tea.Quit

		case "pgup":
			viewHeight := m.viewHeight()
			m.cursor -= viewHeight
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.ensureVisible()

		case "pgdown":
			viewHeight := m.viewHeight()
			m.cursor += viewHeight
			if m.cursor >= len(m.visible) {
				m.cursor = len(m.visible) - 1
			}
			m.ensureVisible()
		}
	}
	return m, nil
}

func (m *Model) ensureVisible() {
	viewHeight := m.viewHeight()
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+viewHeight {
		m.offset = m.cursor - viewHeight + 1
	}
}

func (m Model) viewHeight() int {
	// header (2 lines) + footer (2 lines) = 4
	h := m.height - 4
	if h < 1 {
		h = 1
	}
	return h
}

func (m Model) View() string {
	if m.showHelp {
		return "Keybindings:\n\n" +
			"  j/k, C-n/C-p, ↑/↓  Move cursor\n" +
			"  h/l, ←/→, Enter     Collapse/Expand directory\n" +
			"  Space               Toggle check\n" +
			"  a                   Check all\n" +
			"  n                   Uncheck all\n" +
			"  PgUp/PgDn           Page scroll\n" +
			"  s/w                 Save config\n" +
			"  q, C-c              Quit\n" +
			"\nPress any key to close help"
	}

	var b strings.Builder

	// Header
	b.WriteString(fmt.Sprintf("list-codes select - %s\n", m.rootPath))
	b.WriteString("?: help | s: save | q: quit\n")

	viewHeight := m.viewHeight()
	end := m.offset + viewHeight
	if end > len(m.visible) {
		end = len(m.visible)
	}

	for i := m.offset; i < end; i++ {
		node := m.visible[i]

		// Cursor indicator
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		// Checkbox
		checkbox := "[ ]"
		switch node.State {
		case Checked:
			checkbox = "[x]"
		case Partial:
			checkbox = "[-]"
		}

		// Indentation
		indent := strings.Repeat("  ", node.Depth)

		// Expand indicator for directories
		suffix := ""
		if node.IsDir {
			if node.Expanded {
				suffix = " ▼"
			} else {
				suffix = " ▶"
			}
			b.WriteString(fmt.Sprintf("%s%s %s%s/\n", cursor, checkbox, indent, node.Name+suffix))
		} else {
			b.WriteString(fmt.Sprintf("%s%s %s%s\n", cursor, checkbox, indent, node.Name))
		}
	}

	// Footer
	selected, total := CountSelected(m.root)
	b.WriteString(fmt.Sprintf("\n%d/%d files selected", selected, total))
	if m.statusMsg != "" {
		b.WriteString(" | " + m.statusMsg)
	}

	return b.String()
}

func (m Model) Saved() bool {
	return m.saved
}

func (m Model) ConfigPath() string {
	return m.configPath
}

func RunTUI(rootPath, configPath string, noConfig bool, opts BuildTreeOpts) error {
	m, err := NewModel(rootPath, configPath, noConfig, opts)
	if err != nil {
		return err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	fm := finalModel.(Model)
	if fm.Saved() {
		fmt.Printf("Config saved to %s\n", fm.ConfigPath())
	}
	return nil
}
