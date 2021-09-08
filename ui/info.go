package ui

import (
	"fmt"

	api "github.com/ipfs/go-ipfs-api"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

type Info struct {
	*cview.TextView
	app *App
}

func NewInfo(app *App) *Info {
	m := &Info{
		TextView: cview.NewTextView(),
		app:      app,
	}
	m.SetBorder(true)
	m.SetPadding(1, 1, 1, 1)
	m.SetTitle("info")
	m.SetBackgroundColor(tcell.ColorDefault)

	return m
}

func (i *Info) SetItem(entry *api.MfsLsEntry) {
	info := fmt.Sprintf("%+v", entry)

	stat, err := i.app.client.StatEntry(entry)
	if err != nil {
		panic(err)
	}
	info = fmt.Sprintf("%s\n%s", entry.Name, stat)

	i.SetText(info)
}
