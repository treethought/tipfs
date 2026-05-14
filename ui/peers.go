package ui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PeerTable struct {
	*tview.Table
	app *App
}

func NewPeerTable(app *App) *PeerTable {
	m := &PeerTable{
		Table: tview.NewTable(),
		app:   app,
	}
	m.SetBorder(true)
	m.SetTitle("peers")
	m.SetBackgroundColor(tcell.ColorDefault)
	m.SetBorders(true)

	peers, err := app.ipfs.GetPeers()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	headers := []string{"peer", "address", "latency"}
	for col, h := range headers {
		m.SetCell(0, col, tview.NewTableCell(h).SetTextColor(tcell.ColorYellow))
	}
	m.SetFixed(1, 0)

	for row, p := range peers.Peers {
		m.SetCell(row+1, 0, tview.NewTableCell(p.Peer))
		m.SetCell(row+1, 1, tview.NewTableCell(p.Addr))
		m.SetCell(row+1, 2, tview.NewTableCell(p.Latency))
	}

	m.SetSelectable(true, false)

	return m
}

func (r *PeerTable) Update() {}
