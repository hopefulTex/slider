// A slider util that lets users pick an int from a given range, good for phone numbers
package slider

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	Decrement key.Binding
	Increment key.Binding
}

var DefaultKeyMap = KeyMap{
	Increment: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "Increment"),
	),
	Decrement: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "Decrement"),
	),
}

type Model struct {
	Min          int  // minimum number
	Max          int  // maximum number
	steps        int  // number incremented per tick
	Value        int  // current value selected
	width        int  // width of the bar
	handleRune   rune // symbol for the slider
	barRune      rune // symbol for the slide
	ShowHiLo     bool // show the min/max values allowed
	ShowTop      bool // show min/max on top of the slider
	ShowValue    bool // show the current value
	ShowBottom   bool // show min/max under the slider
	ShowValLeft  bool // show current value to the left of the slider
	ShowValRight bool // show current value to the right of the slider

	fillingColor bool
	handleColor  lipgloss.Color
	barColor     lipgloss.Color
	textColor    lipgloss.Color
	border       lipgloss.Border
}

func New() Model {
	return Model{
		Min:          0,
		Max:          10,
		steps:        1,
		Value:        0,
		width:        30,
		handleRune:   '┃',
		barRune:      '─',
		ShowHiLo:     false,
		ShowTop:      true,
		ShowBottom:   true,
		ShowValRight: false,
		ShowValLeft:  false,
		ShowValue:    true,
		fillingColor: false, // bar colors change only if behind the slider
		handleColor:  lipgloss.Color("#894593"),
		barColor:     lipgloss.Color("#5a32e2"),
		textColor:    lipgloss.Color(""),
		border:       lipgloss.NormalBorder(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Increment):
			m.Value = (min(m.Value+m.steps, m.Max))
		case key.Matches(msg, DefaultKeyMap.Decrement):
			m.Value = (max(m.Value-m.steps, m.Min))
		}
	}

	return m, nil
}

func (m Model) View() string {
	var textStyle = lipgloss.NewStyle().Foreground(m.textColor).Padding(0).Margin(0)
	var barStyle = textStyle.Copy().Foreground(m.barColor)
	var handleStyle = textStyle.Copy().Foreground(m.handleColor)
	var rangeString strings.Builder
	var bar strings.Builder
	var slider strings.Builder
	var value strings.Builder
	var padding int = max(len(fmt.Sprintf("%d", m.Max)), len(fmt.Sprintf("%d", m.Min)))

	if m.ShowValue {
		value.WriteString(fmt.Sprintf("%d", m.Value))
		var tmp int = padding - value.Len()
		for i := 0; i < tmp; i++ {
			value.WriteRune(' ')
		}
	}

	if m.ShowHiLo {
		if m.ShowValLeft && m.ShowValue {
			for i := 0; i <= padding; i++ {
				rangeString.WriteRune(' ')
			}
		}
		rangeString.WriteString(fmt.Sprintf("%d", m.Min))
		var tmp int = m.width - len(fmt.Sprintf("%d", m.Min)) - len(fmt.Sprintf("%d", m.Max))
		for i := 0; i <= tmp; i++ {
			rangeString.WriteRune(' ')
		}
		rangeString.WriteString(fmt.Sprintf("%d\n", m.Max))
	}

	if m.ShowTop {
		slider.WriteString(rangeString.String())
	}

	if m.ShowValLeft {
		slider.WriteString(fmt.Sprintf("%s ", value.String()))

	}

	var decSteps float64 = float64(m.width) / float64((m.Max - m.Min))
	for i := 0; i <= m.width; i++ {
		if i == int(decSteps*float64(m.Value)) {
			bar.WriteString(handleStyle.Render(string(m.handleRune)))
		} else {
			if m.fillingColor {
				if i <= int(decSteps*float64(m.Value)) {
					bar.WriteString(
						barStyle.
							Foreground(m.barColor).
							Render(string(m.barRune)))
				} else {
					bar.WriteString(
						barStyle.UnsetForeground().Render(string(m.barRune)))
				}
			} else {
				bar.WriteString(barStyle.Render(string(m.barRune)))
			}
		}
	}
	slider.WriteString(bar.String())

	if m.ShowValRight {
		slider.WriteString(fmt.Sprintf(" %s", value.String()))
	}

	if m.ShowBottom && m.ShowHiLo {
		slider.WriteRune('\n')
		slider.WriteString(rangeString.String())
	}

	return lipgloss.NewStyle().Border(m.border).Render(slider.String())
}

func min(nums ...int) int {
	var smallest = nums[0]
	for _, x := range nums {
		if x < smallest {
			smallest = x
		}
	}

	return smallest
}

func max(nums ...int) int {
	var largest = nums[0]
	for _, x := range nums {
		if x > largest {
			largest = x
		}
	}

	return largest
}
