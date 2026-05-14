package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/treethought/tipfs/ipfs"
)

type App struct {
	ipfs      *ipfs.Client
	ui        *tview.Application
	root      *tabbedPanels
	focusList []tview.Primitive
	focusIdx  int
	state     *State
	widgets   []Widget
}

type Widget interface {
	Update()
}

type State struct {
	app         *App
	currentFile TreeEntry
}

func NewState(app *App) *State {
	return &State{
		app:         app,
		currentFile: TreeEntry{path: "/", entry: nil},
	}
}

func (s *State) SetItem(e TreeEntry) {
	s.currentFile = e
	for _, w := range s.app.widgets {
		w.Update()
	}
}

func New() *App {
	app := &App{
		ipfs:    ipfs.NewClient(""),
		widgets: make([]Widget, 0),
	}
	app.state = NewState(app)
	return app
}

// tabbedPanels is a minimal tab container: a Pages for content and a one-line
// TextView at the bottom for the tab bar.
type tabbedPanels struct {
	*tview.Flex
	pages   *tview.Pages
	tabBar  *tview.TextView
	tabs    []struct{ id, label string }
	current string
}

func newTabbedPanels() *tabbedPanels {
	tp := &tabbedPanels{
		Flex:   tview.NewFlex().SetDirection(tview.FlexRow),
		pages:  tview.NewPages(),
		tabBar: tview.NewTextView().SetDynamicColors(true),
	}
	tp.pages.SetBackgroundColor(tcell.ColorDefault)
	tp.tabBar.SetBackgroundColor(tcell.ColorDefault)
	tp.Flex.AddItem(tp.pages, 0, 1, true)
	tp.Flex.AddItem(tp.tabBar, 1, 0, false)
	return tp
}

func (tp *tabbedPanels) addTab(id, label string, content tview.Primitive) {
	first := len(tp.tabs) == 0
	tp.tabs = append(tp.tabs, struct{ id, label string }{id, label})
	tp.pages.AddPage(id, content, true, first)
	if first {
		tp.current = id
	}
	tp.renderTabBar()
}

func (tp *tabbedPanels) setCurrentTab(id string) {
	tp.current = id
	tp.pages.SwitchToPage(id)
	tp.renderTabBar()
}

func (tp *tabbedPanels) renderTabBar() {
	var b strings.Builder
	for i, tab := range tp.tabs {
		if i > 0 {
			b.WriteString(" | ")
		}
		if tab.id == tp.current {
			fmt.Fprintf(&b, "[::b]%s[::-]", tab.label)
		} else {
			b.WriteString(tab.label)
		}
	}
	tp.tabBar.SetText(b.String())
}

func (app *App) focusNext() {
	if len(app.focusList) == 0 {
		return
	}
	app.focusIdx = (app.focusIdx + 1) % len(app.focusList)
	app.ui.SetFocus(app.focusList[app.focusIdx])
}

func (app *App) initFilesLayout() *tview.Flex {
	repo := NewRepoTree(app)
	info := NewFileInfo(app)
	dag := NewDagInfo(app)
	content := NewContentView(app)
	app.widgets = append(app.widgets, repo, info, dag, content)
	app.focusList = append(app.focusList, repo, content, dag, info)

	side := tview.NewFlex().SetDirection(tview.FlexRow)
	side.SetBackgroundColor(tcell.ColorDefault)
	side.AddItem(repo, 0, 2, true)
	side.AddItem(info, 0, 1, false)

	mid := tview.NewFlex().SetDirection(tview.FlexRow)
	mid.SetBackgroundColor(tcell.ColorDefault)
	mid.AddItem(content, 0, 2, false)
	mid.AddItem(dag, 0, 1, false)

	flex := tview.NewFlex()
	flex.SetBackgroundColor(tcell.ColorDefault)
	flex.AddItem(side, 0, 2, true)
	flex.AddItem(mid, 0, 4, false)

	return flex
}

func (app *App) initPeersLayout() *tview.Flex {
	peers := NewPeerTable(app)
	app.widgets = append(app.widgets, peers)
	app.focusList = append(app.focusList, peers)

	flex := tview.NewFlex()
	flex.SetBackgroundColor(tcell.ColorDefault)
	flex.AddItem(peers, 0, 1, true)
	return flex
}

func (app *App) initViews() {
	filesFlex := app.initFilesLayout()
	peersFlex := app.initPeersLayout()

	dataPanels := newTabbedPanels()
	dataPanels.addTab("files", "files", filesFlex)
	dataPanels.addTab("peers", "peers", peersFlex)
	dataPanels.setCurrentTab("files")
	dataPanels.SetBackgroundColor(tcell.ColorDefault)

	app.root = dataPanels
}

func (app *App) initBindings() {
	app.ui.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyTAB:
			app.focusNext()
			return nil
		}
		switch ev.Rune() {
		case '1':
			app.root.setCurrentTab("files")
			return nil
		case '2':
			app.root.setCurrentTab("peers")
			return nil
		}
		return ev
	})
}

func (app *App) Start() {
	app.ui = tview.NewApplication()
	app.initViews()
	app.initBindings()
	app.ui.SetRoot(app.root, true)
	app.ui.EnableMouse(true)
	if err := app.ui.Run(); err != nil {
		log.Fatal(err)
	}
}
