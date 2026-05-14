package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FileInfo struct {
	*tview.TextView
	app *App
}

func NewFileInfo(app *App) *FileInfo {
	m := &FileInfo{
		TextView: tview.NewTextView(),
		app:      app,
	}
	m.SetBorder(true)
	m.SetTitle("info")
	m.SetBackgroundColor(tcell.ColorDefault)
	return m
}

func (i *FileInfo) Update() {
	current := i.app.state.currentFile
	if current.entry == nil {
		return
	}

	go func() {
		stat, err := i.app.ipfs.StatFile(current.path, current.entry)
		i.app.ui.QueueUpdateDraw(func() {
			if err != nil {
				i.SetText(fmt.Sprintf("%s\n%v", current.path, err))
				return
			}
			i.SetText(fmt.Sprintf("%s\n%s", current.entry.Name, stat))
		})
	}()
}
