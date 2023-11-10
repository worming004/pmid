package main

import (
	"bubblepod/pkg/log"
	"bubblepod/pkg/podman"
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var defaultItemAssert list.DefaultItem = item{}
var itemAssert list.Item = item{}

type item struct {
	title, id string
}

func toItem(img podman.Image) item {
	return item{title: img.Name, id: img.Id}
}
func toItems(imgs []podman.Image) []list.Item {
	var result = make([]list.Item, len(imgs))
	for i, img := range imgs {
		result[i] = toItem(img)
	}
	return result
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.id }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "d" {
			index := m.list.Index()
			m.list.RemoveItem(index)
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	closer := log.SetupDefaultLogger()
	defer closer.Close()

	podmanImages, err := podman.PodmanCommands{}.GetImages()
	if err != nil {
		slog.Error(err.Error())
	}
	var items []list.Item = toItems(podmanImages)

	slog.Info(items[0].(item).Title())

	// m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.KeyMap.NextPage.SetKeys("rigth", "l", "pgdown", "f")

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
