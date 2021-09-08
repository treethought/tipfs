package ui

import (
	"fmt"

	api "github.com/ipfs/go-ipfs-api"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

type Buffer struct {
	*cview.Box
	text *cview.TextView
	tree *cview.TreeView

	app   *App
	entry *api.MfsLsEntry
}

func NewBuffer(app *App) *Buffer {
	m := &Buffer{
		text: cview.NewTextView(),
		tree: cview.NewTreeView(),
		app:  app,
	}
	// m.SetBorder(true)
	// m.SetPadding(1, 1, 1, 1)
	// m.SetTitle("info")
	// m.SetBackgroundColor(tcell.ColorDefault)

	return m
}

func (i *Buffer) Draw(screen tcell.Screen) {
	if i.entry.Type == api.TDirectory {
		i.tree.Draw(screen)
		return
	}
	i.text.Draw(screen)
}

func (i *Buffer) SetItem(entry *api.MfsLsEntry) {
	i.entry = entry
	info := fmt.Sprintf("%+v", entry)

	stat, err := i.app.client.StatEntry(entry)
	if err != nil {
		panic(err)
	}

	// if entry.Type == api.TDirectory {
	// 	children, err := i.app.client.ListFiles(entry.Name)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	// for _, e := range children {

	// 	// }

	// }

	info = fmt.Sprintf("%s\n%s", entry.Name, stat)

	i.text.SetText(info)
}
