package ui

import (
	"fmt"
	"os"
	"time"

	"code.rocketnine.space/tslocum/cbind"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

type PeerTable struct {
	*cview.Table
	app          *App
	inputHandler *cbind.Configuration
}

func NewPeerTable(app *App) *PeerTable {
	m := &PeerTable{
		Table: cview.NewTable(),
		app:   app,
	}
	m.SetBorder(true)
	m.SetPadding(1, 1, 1, 1)
	m.SetTitle("peers")
	m.SetBackgroundColor(tcell.ColorDefault)

	m.inputHandler = cbind.NewConfiguration()
	m.initBindings()

	m.SetEvaluateAllRows(true)
	m.SetScrollBarVisibility(cview.ScrollBarNever)
	m.SetBorders(true)

	peers, err := app.ipfs.GetPeers()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m.SetCell(0, 0, cview.NewTableCell("id"))
	m.SetCell(0, 1, cview.NewTableCell("address"))
	m.SetCell(0, 2, cview.NewTableCell("latency"))

	m.SetFixed(1, 0)

	for r, p := range peers {
		id := cview.NewTableCell(p.ID().Pretty())
		m.SetCell(r+1, 0, id)
		addr := cview.NewTableCell(p.Address().String())
		m.SetCell(r+1, 1, addr)
		lat, err := p.Latency()
		if err != nil {
			lat = time.Duration(0)
		}
		latency := cview.NewTableCell(lat.String())
		m.SetCell(r+1, 2, latency)

	}

	m.SetSelectable(true, false)

	return m
}

func (r *PeerTable) Update() {}

func (r *PeerTable) handleSelect(ev *tcell.EventKey) *tcell.EventKey {
	return nil
}

func (t *PeerTable) initBindings() {
	// t.inputHandler.SetKey(tcell.ModNone, tcell.KeyEnter, t.handleSelect)
	// t.inputHandler.SetRune(tcell.ModNone, 'o', t.handleOpen)
	t.SetInputCapture(t.inputHandler.Capture)
}
