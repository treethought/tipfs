package ui

import (
	"fmt"


	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

type FileInfo struct {
	*cview.TextView
	app *App
}

func NewFileInfo(app *App) *FileInfo {
	m := &FileInfo{
		TextView: cview.NewTextView(),
		app:      app,
	}
	m.SetBorder(true)
	m.SetPadding(1, 1, 1, 1)
	m.SetTitle("info")
	m.SetBackgroundColor(tcell.ColorDefault)

	return m
}

func (i *FileInfo) Update() {
	current := i.app.state.currentItem
	info := fmt.Sprintf("%+v", current.entry)
	i.Clear()

	go i.app.ui.QueueUpdateDraw(func() {
		stat, err := i.app.client.StatFile(current.path, current.entry)
		if err != nil {
			i.SetText(fmt.Sprintf("%s\n%v", current.path, err))
			return
		}
		info = fmt.Sprintf("%s\n%s", current.entry.Name, stat)

		i.SetText(info)

	})
}
