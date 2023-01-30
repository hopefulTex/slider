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
	min          int		// minimum number
	max          int		// maximum number
	steps        int		// number incremented per tick
	value        int		// current value selected
	width        int		// width of the bar
	handleRune   rune		// symbol for the slider
	barRune      rune		// symbol for the slide
	showHiLo     bool		// show the min/max values allowed
	showTop      bool		// show min/max on top of the slider
	showValue    bool		// show the current value
	showBottom   bool		// show min/max under the slider
	showValLeft  bool		// show current value to the left of the slider
	showValRight bool		// show current value to the right of the slider
	
	fillingColor bool
	handleColor	lipgloss.Color
	barColor		lipgloss.Color
	textColor	lipgloss.Color
	border	lipgloss.Border
}

func New() Model {
	return Model{
		min:        0,
		max:        10,
		steps:      1,
		value:      0,
		width:      30,
		handleRune: '┃',
		barRune:    '─',
		showHiLo: false,
		showTop:    true,
		showBottom: true,
		showValRight: false,
		showValLeft: false,
		showValue: true,
		fillingColor: false,	// bar colors change only if behind the slider
		handleColor: lipgloss.Color("#894593"),
		barColor: lipgloss.Color("#5a32e2"),
		textColor: lipgloss.Color(""),
		border:      lipgloss.NormalBorder(),
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
			m.value = (min(m.value+m.steps, m.max))
		case key.Matches(msg, DefaultKeyMap.Decrement):
			m.value = (max(m.value-m.steps, m.min))
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
	var padding int = max(len(fmt.Sprintf("%d", m.max)), len(fmt.Sprintf("%d", m.min)))

	if m.showValue {
		value.WriteString(fmt.Sprintf("%d", m.value))
		var tmp int = padding - value.Len()
		for i := 0; i < tmp; i++ {
			value.WriteRune(' ')
		}
	}

	if m.showHiLo {
		if m.showValLeft && m.showValue {
			for i := 0; i <= padding; i++ {
				rangeString.WriteRune(' ')
			}
		}
		rangeString.WriteString(fmt.Sprintf("%d", m.min))
		var tmp int = m.width - len(fmt.Sprintf("%d", m.min)) - len(fmt.Sprintf("%d", m.max))
		for i := 0; i <= tmp; i++ {
			rangeString.WriteRune(' ')
		}
		rangeString.WriteString(fmt.Sprintf("%d\n", m.max))
	}

	if m.showTop {
		slider.WriteString(rangeString.String())
	}

	if m.showValLeft {
		slider.WriteString(fmt.Sprintf("%s ",value.String()))
		
	}

	var decSteps float64 = float64(m.width) / float64((m.max - m.min))
	for i := 0; i <= m.width; i++ {
		if i == int(decSteps*float64(m.value)) {
			bar.WriteString(handleStyle.Render(string(m.handleRune)))
		} else {
			if m.fillingColor {
				if i <= int(decSteps*float64(m.value)) {
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

	if m.showValRight {
		slider.WriteString(fmt.Sprintf(" %s",value.String()))
	}

	if m.showBottom && m.showHiLo {
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

