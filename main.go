package main

import (
	"github.com/Livenux/virt-mon/cmd"
	tea "github.com/charmbracelet/bubbletea"
	"libvirt.org/go/libvirt"
	"log"
	"time"
)

func main() {
	conn, err := libvirt.NewConnect("test:///default")
	if err != nil {
		panic(err)
	}
	p := tea.NewProgram(cmd.NewModel(conn, time.Second))

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
