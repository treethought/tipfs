package ui

import (
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"

	"code.rocketnine.space/tslocum/cbind"
	"code.rocketnine.space/tslocum/cview"
	ipld "github.com/ipfs/go-ipld-format"
)

type DagInfo struct {
	*cview.TreeView
	app          *App
	inputHandler *cbind.Configuration
	currentEntry TreeEntry
	currentHash  string
}

func NewDagInfo(app *App) *DagInfo {

	m := &DagInfo{
		app:      app,
		TreeView: cview.NewTreeView(),
	}

	m.SetBackgroundColor(tcell.ColorDefault)

	root := cview.NewTreeNode("")
	m.SetRoot(root)
	m.SetCurrentNode(root)

	m.SetBorder(true)
	m.SetPadding(1, 1, 1, 1)
	m.SetTitle("dag")
	m.SetBackgroundColor(tcell.ColorDefault)

	m.inputHandler = cbind.NewConfiguration()
	m.initBindings()

	return m
}

func (i *DagInfo) handleSelect(ev *tcell.EventKey) *tcell.EventKey {
	node := i.GetCurrentNode()
	ref := node.GetReference()
	link, ok := ref.(ipld.Link)
	if !ok {
		return nil
	}

	if len(node.GetChildren()) > 0 {
		node.SetExpanded(true)
	}

	i.currentHash = link.Cid.String()
	i.Update()

	return nil
}

func (d *DagInfo) initBindings() {
	d.inputHandler.SetKey(tcell.ModNone, tcell.KeyEnter, d.handleSelect)
	d.inputHandler.SetRune(tcell.ModNone, 'y', func(ev *tcell.EventKey) *tcell.EventKey {
		node := d.GetCurrentNode()
		ref := node.GetReference()
		e, ok := ref.(ipld.Link)
		if !ok {
			return nil
		}
		clipboard.WriteAll(e.Cid.String())

		return nil
	})
	d.SetInputCapture(d.inputHandler.Capture)
}

func (i *DagInfo) Update() {
	i.GetRoot().ClearChildren()
	fileNode := i.app.state.currentFile

	// new file was selected, show dag for that
	if i.currentEntry.path != fileNode.path {
		i.currentEntry = fileNode
		i.currentHash = fileNode.entry.Hash
	}

	current := i.GetCurrentNode()
	if current == nil {
		current = cview.NewTreeNode(fileNode.entry.Hash)
	}

	i.GetRoot().ClearChildren()
	i.SetRoot(i.GetCurrentNode())

	go i.app.ui.QueueUpdateDraw(func() {

		dag, err := i.app.ipfs.GetDag(i.currentHash)
		if err != nil {
			panic(err)
		}

		lines := []string{}
		truncData := truncateMiddle(dag.Data, 12)
		lines = append(lines, "data:", truncData)
		lines = append(lines, "links:")

		for _, l := range dag.Links {
			node := cview.NewTreeNode(l.Cid.String())
			node.SetReference(l)
			size := cview.NewTreeNode(byteCount(l.Size))
			name := cview.NewTreeNode(l.Name)
			node.AddChild(name)
			node.AddChild(size)
			i.GetRoot().AddChild(node)
		}

	})
}
