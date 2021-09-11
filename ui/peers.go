package ui

import (
	"fmt"
	"os"

	"code.rocketnine.space/tslocum/cbind"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

type PeerList struct {
	*cview.List
	app          *App
	inputHandler *cbind.Configuration
}

func NewPeerList(app *App) *PeerList {
	m := &PeerList{
		List: cview.NewList(),
		app:  app,
	}
	m.SetBorder(false)
	m.SetBorderAttributes(tcell.AttrDim)
	m.SetPadding(1, 1, 1, 1)
	m.SetTitle("peers")
	m.SetBackgroundColor(tcell.ColorDefault)

	m.inputHandler = cbind.NewConfiguration()
	m.initBindings()

	peers, err := app.ipfs.GetPeers()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, p := range peers.Peers {
		item := cview.NewListItem(p.Peer)
		item.SetSecondaryText(p.Addr)
		item.SetReference(p)
		m.AddItem(item)
	}

	return m
}

func (r *PeerList) Update() {}

func (r *PeerList) handleSelect(ev *tcell.EventKey) *tcell.EventKey {
	// item := r.GetCurrentItem()
	// ref := item.GetReference()
	// connInf, _ := ref.(*api.SwarmConnInfo)

	// r.app.info.SetItem(e.path, e.entry)
	// r.app.dag.SetItem(e.entry)
	// r.app.content.SetItem(e.path, e.entry)
	return nil
}

func (t *PeerList) initBindings() {
	// t.inputHandler.SetKey(tcell.ModNone, tcell.KeyEnter, t.handleSelect)
	// t.inputHandler.SetRune(tcell.ModNone, 'o', t.handleOpen)
	t.SetInputCapture(t.inputHandler.Capture)
}
