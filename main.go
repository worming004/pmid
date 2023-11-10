package main

import (
	"fmt"
	"github.com/worming004/pmid/pkg/log"
	"github.com/worming004/pmid/pkg/podman"
	"log/slog"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var defaultItemAssert list.DefaultItem = item{}
var itemAssert list.Item = item{}

var podmanCommands = podman.Default

type item struct {
	title, id, tag string
}

func toItem(img podman.Image) item {
	return item{title: img.Name, id: img.Id, tag: img.Tag}
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
			it := m.list.SelectedItem().(item)
			err := podmanCommands.DeleteImageById(it.id)
			if err != nil {
				slog.Error(err.Error())
			}
			m.list.SetItems(getList())
			return m, nil
		}
		if msg.String() == "t" {
			it := m.list.SelectedItem().(item)
			err := podmanCommands.DeleteImageByTag(it.title, it.tag)
			if err != nil {
				slog.Error(err.Error())
			}
			m.list.SetItems(getList())
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

	items := getList()

	// m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func getList() []list.Item {
	podmanImages, err := podmanCommands.GetImages()
	if err != nil {
		slog.Error(err.Error())
	}
	var items []list.Item = toItems(podmanImages)
	return items
}
