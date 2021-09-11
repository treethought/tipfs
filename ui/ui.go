package ui

import (
	"log"

	"code.rocketnine.space/tslocum/cbind"
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
	"github.com/treethought/tipfs/ipfs"
)

type App struct {
	ipfs         *ipfs.Client
	ui           *cview.Application
	root         *cview.Flex
	dataPanels   *cview.TabbedPanels
	focusManager *cview.FocusManager
	state        *State
	widgets      []Widget
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

func (app *App) initViews() {
	repo := NewRepoTree(app)
	info := NewFileInfo(app)
	dag := NewDagInfo(app)
	content := NewContentView(app)

	peers := NewPeerList(app)

	app.widgets = append(app.widgets, repo, info, dag, content, peers)

	dataPanels := cview.NewTabbedPanels()
	dataPanels.AddTab("files", "files", repo)
	dataPanels.AddTab("peers", "peers", peers)

	dataPanels.SetBackgroundColor(tcell.ColorDefault)
	dataPanels.SetTabBackgroundColor(tcell.ColorDefault)
	dataPanels.SetTabSwitcherDivider("", " | ", "")
	dataPanels.SetTabSwitcherAfterContent(true)
	// dataPanels.SetDirection(cview.FlexColumn)
	dataPanels.SetBorder(true)

	app.dataPanels = dataPanels
	// app.dataPanels.AddItem(app.repo, 0, 4, true)

	mid := cview.NewFlex()
	mid.SetBackgroundColor(tcell.ColorDefault)
	mid.SetDirection(cview.FlexRow)
	// mid.AddItem(app.dataPanels, 0, 4, true)
	mid.AddItem(content, 0, 4, false)
	// mid.AddItem(app.info, 0, 2, false)
	mid.AddItem(dag, 0, 2, false)

	flex := cview.NewFlex()
	flex.SetBackgroundTransparent(false)
	flex.SetBackgroundColor(tcell.ColorDefault)

	left := cview.NewFlex()
	left.SetDirection(cview.FlexRow)
	// left.AddItem(app.repo, 0, 7, false)
	left.AddItem(app.dataPanels, 0, 4, false)
	left.AddItem(info, 0, 2, false)

	flex.AddItem(left, 0, 2, false)
	flex.AddItem(mid, 0, 4, false)
	app.root = flex

	app.initInputHandler(repo, content, dag, info)

}

func (app *App) handleToggle(ev *tcell.EventKey) *tcell.EventKey {
	app.focusManager.FocusNext()
	return nil

}

func (app *App) initBindings() {
	c := cbind.NewConfiguration()
	c.SetKey(tcell.ModNone, tcell.KeyTAB, app.handleToggle)
	c.SetKey(tcell.ModNone, tcell.KeyF1, app.handleToggle)
	c.SetRune(tcell.ModNone, '1', func(ev *tcell.EventKey) *tcell.EventKey {
		app.dataPanels.SetCurrentTab("files")
		return nil
	})
	c.SetRune(tcell.ModNone, '2', func(ev *tcell.EventKey) *tcell.EventKey {
		app.dataPanels.SetCurrentTab("peers")
		return nil
	})
	app.ui.SetInputCapture(c.Capture)
}

func (app *App) initInputHandler(widgets ...cview.Primitive) {
	widgets = append(widgets, app.dataPanels)
	app.focusManager = cview.NewFocusManager(app.ui.SetFocus)
	app.focusManager.SetWrapAround(true)
	app.focusManager.Add(widgets...)
}

func (app *App) Start() {

	// Initialize application
	app.ui = cview.NewApplication()

	app.initViews()
	app.initBindings()

	app.ui.SetRoot(app.root, true)
	app.ui.SetFocus(app.dataPanels)

	app.ui.EnableMouse(true)

	err := app.ui.Run()
	if err != nil {
		log.Fatal(err)
	}
}
