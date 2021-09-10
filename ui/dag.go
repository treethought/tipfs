package ui

import (
	"github.com/gdamore/tcell/v2"
	"gopkg.in/yaml.v2"

	"code.rocketnine.space/tslocum/cview"
)

type DagInfo struct {
	*cview.TextView
	app *App
}

func NewDagInfo(app *App) *DagInfo {
	m := &DagInfo{
		TextView: cview.NewTextView(),
		app:      app,
	}
	m.SetBorder(true)
	m.SetPadding(1, 1, 1, 1)
	m.SetTitle("dag")
	m.SetBackgroundColor(tcell.ColorDefault)

	return m
}

func (i *DagInfo) Update() {
	entry := i.app.state.currentItem.entry
	i.SetText("loading...")

	go i.app.ui.QueueUpdateDraw(func() {

		dag, err := i.app.client.GetDag(entry.Hash)
		if err != nil {
			panic(err)
		}

		data, err := yaml.Marshal(dag)
		if err != nil {
			panic(err)
		}
		i.SetText(string(data))
	})
}
