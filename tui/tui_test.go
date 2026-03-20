package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestModel(t *testing.T) Model {
	t.Helper()
	dir := createTestProject(t)
	m, err := NewModel(dir, "", false, BuildTreeOpts{})
	require.NoError(t, err)
	return m
}

func TestNewModel(t *testing.T) {
	m := newTestModel(t)
	assert.NotNil(t, m.root)
	assert.True(t, len(m.visible) > 0)
	assert.Equal(t, 0, m.cursor)
}

func TestModel_NavigateDown(t *testing.T) {
	m := newTestModel(t)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)
	assert.Equal(t, 1, m.cursor)
}

func TestModel_NavigateUp(t *testing.T) {
	m := newTestModel(t)

	// Move down first
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)
	assert.Equal(t, 1, m.cursor)

	// Move up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(Model)
	assert.Equal(t, 0, m.cursor)
}

func TestModel_NavigateUp_AtTop(t *testing.T) {
	m := newTestModel(t)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(Model)
	assert.Equal(t, 0, m.cursor)
}

func TestModel_Toggle(t *testing.T) {
	m := newTestModel(t)

	// Move to first item and toggle
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = updated.(Model)

	node := m.visible[0]
	assert.NotEqual(t, Unchecked, node.State)
}

func TestModel_ExpandCollapse(t *testing.T) {
	m := newTestModel(t)

	// Root is already expanded, find a collapsed dir
	initialLen := len(m.visible)

	// Move to a directory and expand it
	for i, node := range m.visible {
		if node.IsDir && !node.Expanded && len(node.Children) > 0 {
			m.cursor = i
			updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
			m = updated.(Model)
			assert.True(t, len(m.visible) > initialLen)
			break
		}
	}
}

func TestModel_SelectAll(t *testing.T) {
	m := newTestModel(t)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(Model)

	selected, total := CountSelected(m.root)
	assert.Equal(t, total, selected)
}

func TestModel_SelectNone(t *testing.T) {
	m := newTestModel(t)

	// Select all first
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(Model)

	// Deselect all
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m = updated.(Model)

	selected, _ := CountSelected(m.root)
	assert.Equal(t, 0, selected)
}

func TestModel_View(t *testing.T) {
	m := newTestModel(t)
	m.width = 80
	m.height = 24

	view := m.View()
	assert.Contains(t, view, "list-codes select")
	assert.Contains(t, view, "files selected")
}

func TestModel_Quit(t *testing.T) {
	m := newTestModel(t)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.NotNil(t, cmd)
}

func TestModel_WindowResize(t *testing.T) {
	m := newTestModel(t)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(Model)
	assert.Equal(t, 120, m.width)
	assert.Equal(t, 40, m.height)
}
