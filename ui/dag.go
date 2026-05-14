package ui

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/rivo/tview"
)

type DagInfo struct {
	*tview.TreeView
	app          *App
	currentEntry TreeEntry
	currentHash  string
}

func NewDagInfo(app *App) *DagInfo {
	m := &DagInfo{
		app:      app,
		TreeView: tview.NewTreeView(),
	}
	m.SetBorder(true)
	m.SetTitle("dag")
	m.SetBackgroundColor(tcell.ColorDefault)

	root := tview.NewTreeNode("")
	m.SetRoot(root)
	m.SetCurrentNode(root)

	m.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		link, ok := ref.(ipld.Link)
		if !ok {
			return
		}
		if len(node.GetChildren()) > 0 {
			node.Expand()
		}
		m.currentHash = link.Cid.String()
		m.Update()
	})

	m.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		if ev.Rune() == 'y' {
			node := m.GetCurrentNode()
			if node == nil {
				return nil
			}
			if link, ok := node.GetReference().(ipld.Link); ok {
				clipboard.WriteAll(link.Cid.String())
			}
			return nil
		}
		return ev
	})

	return m
}

func (i *DagInfo) Update() {
	fileNode := i.app.state.currentFile
	if fileNode.entry == nil {
		return
	}

	if i.currentEntry.path != fileNode.path {
		i.currentEntry = fileNode
		i.currentHash = fileNode.entry.Hash
	}

	root := tview.NewTreeNode(i.currentHash)
	i.SetRoot(root)
	i.SetCurrentNode(root)

	go func() {
		dag, err := i.app.ipfs.GetDag(i.currentHash)
		i.app.ui.QueueUpdateDraw(func() {
			if err != nil {
				root.SetText(fmt.Sprintf("error: %v", err))
				return
			}
			root.AddChild(tview.NewTreeNode("data: " + truncateMiddle(string(dag.Data), 12)))
			linksNode := tview.NewTreeNode(fmt.Sprintf("links (%d)", len(dag.Links)))
			root.AddChild(linksNode)
			for _, l := range dag.Links {
				node := tview.NewTreeNode(l.Cid.String())
				node.SetReference(l)
				node.AddChild(tview.NewTreeNode("name: " + l.Name))
				node.AddChild(tview.NewTreeNode("size: " + byteCount(l.Size)))
				linksNode.AddChild(node)
			}
		})
	}()
}
