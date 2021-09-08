package ui

import (
	"github.com/gdamore/tcell/v2"
	api "github.com/ipfs/go-ipfs-api"
	"gopkg.in/yaml.v2"

	"code.rocketnine.space/tslocum/cview"
)

type DagInfo struct {
	*cview.TextView

	app   *App
	entry *api.MfsLsEntry
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

func (i *DagInfo) SetItem(entry *api.MfsLsEntry) {
	i.entry = entry

	dag, err := i.app.client.GetDag(entry.Hash)
	if err != nil {
		panic(err)
	}

	data, err := yaml.Marshal(dag)
	if err != nil {
		panic(err)
	}

	i.SetText(string(data))
}
