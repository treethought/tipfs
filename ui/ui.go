package ui

import (
	"log"

	"code.rocketnine.space/tslocum/cbind"
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
	"github.com/treethought/tipfs/ipfs"
)

type App struct {
	client       *ipfs.Client
	ui           *cview.Application
	root         *cview.Flex
	menu         *RepoTree
	info         *Info
	buffer       *Buffer
	panels       *cview.Panels
	focusManager *cview.FocusManager
}

func New() *App {
	return &App{}
}

func (app *App) initViews() {
	app.menu = NewRepoTree(app)
	app.info = NewInfo(app)
	app.buffer = NewBuffer(app)

	panels := cview.NewPanels()
	app.panels = panels

	mid := cview.NewFlex()
	mid.SetBackgroundColor(tcell.ColorDefault)
	mid.SetDirection(cview.FlexRow)
	mid.AddItem(app.panels, 0, 4, true)
	mid.AddItem(app.info, 0, 4, false)

	flex := cview.NewFlex()
	flex.SetBackgroundTransparent(false)
	flex.SetBackgroundColor(tcell.ColorDefault)

	left := cview.NewFlex()
	left.SetDirection(cview.FlexRow)
	left.AddItem(app.menu, 0, 7, false)

	flex.AddItem(left, 0, 2, false)
	flex.AddItem(mid, 0, 4, false)
	app.root = flex

}

func (app *App) handleToggle(ev *tcell.EventKey) *tcell.EventKey {
	current, _ := app.panels.GetFrontPanel()
	if current == "compose" {
		return ev
	}
	app.focusManager.FocusNext()
	return nil

}

func (app *App) initBindings() {
	c := cbind.NewConfiguration()
	c.SetKey(tcell.ModNone, tcell.KeyTAB, app.handleToggle)
	app.ui.SetInputCapture(c.Capture)

}

func (app *App) initInputHandler() {
	app.focusManager = cview.NewFocusManager(app.ui.SetFocus)
	app.focusManager.SetWrapAround(true)
	app.focusManager.Add(app.menu, app.panels)

}

func (app *App) Start() {

	app.client = ipfs.NewClient("localhost:5001")

	// Initialize application
	app.ui = cview.NewApplication()

	app.initViews()
	app.initInputHandler()
	app.initBindings()

	app.ui.SetRoot(app.root, true)
	app.ui.SetFocus(app.menu)

	err := app.ui.Run()
	if err != nil {
		log.Fatal(err)
	}
}
