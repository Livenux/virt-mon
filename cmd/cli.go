package cmd

import (
	"fmt"
	"github.com/Livenux/virt-mon/pkg/virt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"libvirt.org/go/libvirt"
	"log"
	"strings"
	"time"
)

const (
	columnName     = "name"
	columnMem      = "mem"
	columnMemUsage = "mem_usage"
	columnVCpus    = "vCpus"
	columnCpuUsage = "cpu_usage"
)

var (
	styleCritical = lipgloss.NewStyle().Foreground(lipgloss.Color("#f00"))
	styleStable   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0"))
	styleGood     = lipgloss.NewStyle().Foreground(lipgloss.Color("#0f0"))
)

type Model struct {
	table       table.Model
	updateDelay time.Duration
	data        []*virt.DomainStat
	conn        *libvirt.Connect
}

func NewModel(conn *libvirt.Connect, d time.Duration) Model {
	return Model{
		conn:        conn,
		updateDelay: d,
	}
}

func (m Model) refreshDataCmd() tea.Msg {
	msg, err := virt.AllDomainStat(m.conn)
	if err != nil {
		log.Fatalln(err)
	}
	return msg
}

func generateColumns() []table.Column {
	return []table.Column{
		table.NewColumn(columnName, "NAME", 16),
		table.NewColumn(columnMem, "MemoryUsage", 8),
		table.NewColumn(columnMemUsage, "Memory_Usage", 8),
		table.NewColumn(columnVCpus, "vCPUS", 8),
		table.NewColumn(columnCpuUsage, "CPU%", 8),
	}
}

func (m Model) Init() tea.Cmd {
	return m.refreshDataCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)

		case "up", "j":
			if m.updateDelay < time.Second {
				m.updateDelay *= 10
			}

		case "down", "k":
			if m.updateDelay > time.Millisecond*1 {
				m.updateDelay /= 10
			}
		}

		// Reapply the new data and the new columns based on critical count
		m.table = m.table.
			WithRows(m.generateRows()).
			WithColumns(generateColumns())

		// This can be from any source, but for demo purposes let's party!
		delay := m.updateDelay
		cmds = append(cmds, func() tea.Msg {
			time.Sleep(delay)
			return m.refreshDataCmd()
		})
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString(
		fmt.Sprintf(
			"Table demo with updating data!  Updating every %v\nPress up/down to update faster/slower\nPress q or ctrl+c to quit\n",
			m.updateDelay,
		))

	pad := lipgloss.NewStyle().Padding(1)

	body.WriteString(pad.Render(m.table.View()))

	return body.String()
}

func (m Model) generateRows() []table.Row {
	var rows []table.Row
	for _, entry := range m.data {
		row := table.NewRow(table.RowData{
			columnName:     entry.Name,
			columnVCpus:    entry.VCpu,
			columnCpuUsage: entry.CpuUsage,
			columnMem:      entry.Memory,
			columnMemUsage: entry.MemoryUsage,
		})

		if entry.MemoryUsage >= 90 || entry.CpuUsage >= 90 {
			row = row.WithStyle(styleCritical)
		} else {
			row = row.WithStyle(styleStable)
		}

		if entry.MemoryUsage <= 70 && entry.CpuUsage <= 70 {
			row = row.WithStyle(styleGood)
		}

		rows = append(rows, row)
	}

	return rows

}
