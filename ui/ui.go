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
	root         *cview.TabbedPanels
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

func (app *App) initFilesLayout() *cview.Flex {
	repo := NewRepoTree(app)
	info := NewFileInfo(app)
	dag := NewDagInfo(app)
	content := NewContentView(app)
	app.widgets = append(app.widgets, repo, info, dag, content)

	// side panel of file tree and info
	side := cview.NewFlex()
	side.SetBackgroundTransparent(false)
	side.SetBackgroundColor(tcell.ColorDefault)
	side.SetDirection(cview.FlexRow)
	side.AddItem(repo, 0, 2, true)
	side.AddItem(info, 0, 1, false)

	// larger main content and dag explorer
	mid := cview.NewFlex()
	mid.SetBackgroundTransparent(false)
	mid.SetBackgroundColor(tcell.ColorDefault)
	mid.SetDirection(cview.FlexRow)
	mid.AddItem(content, 0, 2, false)
	mid.AddItem(dag, 0, 1, false)

	// wrapping cotainer
	flex := cview.NewFlex()
	flex.SetBackgroundTransparent(false)
	flex.SetBackgroundColor(tcell.ColorDefault)
	flex.AddItem(side, 0, 2, true)
	flex.AddItem(mid, 0, 4, false)

	app.initInputHandler(repo, content, dag, info)
	return flex

}
func (app *App) initPeersLayout() *cview.Flex {
	peers := NewPeerTable(app)
	app.widgets = append(app.widgets, peers)

	flex := cview.NewFlex()
	flex.SetBackgroundTransparent(false)
	flex.SetBackgroundColor(tcell.ColorDefault)

	side := cview.NewFlex()
	side.SetBackgroundTransparent(false)
	side.SetBackgroundColor(tcell.ColorDefault)
	side.AddItem(peers, 0, 1, true)

	flex.AddItem(side, 0, 1, true)

	return flex
}

func (app *App) initCIDLayout() *cview.Flex {
	cidview := NewCIDView(app)
	app.widgets = append(app.widgets, cidview)

	flex := cview.NewFlex()
	flex.SetBackgroundTransparent(false)
	flex.SetBackgroundColor(tcell.ColorDefault)

	side := cview.NewFlex()
	side.SetBackgroundTransparent(false)
	side.SetBackgroundColor(tcell.ColorDefault)
	side.AddItem(cidview, 0, 1, true)

	flex.AddItem(side, 0, 1, true)

	return side
}

func (app *App) initViews() {

	filesFlex := app.initFilesLayout()
	peersFlex := app.initPeersLayout()
	cidFlex := app.initCIDLayout()

	dataPanels := cview.NewTabbedPanels()
	dataPanels.SetTitle("panels")
	dataPanels.AddTab("files", "files", filesFlex)
	dataPanels.AddTab("peers", "peers", peersFlex)
	dataPanels.AddTab("cid", "cid", cidFlex)
	dataPanels.SetCurrentTab("files")
	dataPanels.SetBorder(false)
	dataPanels.SetPadding(0, 0, 0, 0)

	dataPanels.SetBackgroundColor(tcell.ColorDefault)
	dataPanels.SetTabBackgroundColor(tcell.ColorDefault)
	dataPanels.SetTabSwitcherDivider("", " | ", "")
	dataPanels.SetTabSwitcherAfterContent(true)

	app.root = dataPanels

}

func (app *App) handleToggle(ev *tcell.EventKey) *tcell.EventKey {
	app.focusManager.FocusNext()
	return nil

}

func (app *App) initBindings() {
	c := cbind.NewConfiguration()
	c.SetKey(tcell.ModNone, tcell.KeyTAB, app.handleToggle)
	c.SetRune(tcell.ModNone, '1', func(ev *tcell.EventKey) *tcell.EventKey {
		app.root.SetCurrentTab("files")
		return nil
	})
	c.SetRune(tcell.ModNone, '2', func(ev *tcell.EventKey) *tcell.EventKey {
		app.root.SetCurrentTab("peers")
		return nil
	})
	c.SetRune(tcell.ModNone, '3', func(ev *tcell.EventKey) *tcell.EventKey {
		app.root.SetCurrentTab("cid")
		return nil
	})
	app.ui.SetInputCapture(c.Capture)
}

func (app *App) initInputHandler(widgets ...cview.Primitive) {
	widgets = append(widgets)
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

	app.ui.EnableMouse(true)

	err := app.ui.Run()
	if err != nil {
		log.Fatal(err)
	}
}
