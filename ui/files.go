package ui

import (
	"fmt"
	"os"

	api "github.com/ipfs/go-ipfs-api"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

type RepoTree struct {
	*cview.TreeView
	app *App
}

func NewRepoTree(app *App) *RepoTree {
	m := &RepoTree{
		TreeView: cview.NewTreeView(),
		app:      app,
	}
	m.SetBorder(true)
	m.SetPadding(1, 1, 1, 1)
	m.SetTitle("menu")
	m.SetBackgroundColor(tcell.ColorDefault)
	m.SetSelectedTextColor(tcell.ColorTeal)

	rootNode := cview.NewTreeNode("/")
	m.SetRoot(rootNode)
	m.SetCurrentNode(rootNode)

	entries, err := app.client.ListFiles("/")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, i := range entries {
		node := cview.NewTreeNode(i.Name)
		node.SetReference(i)
		rootNode.AddChild(node)
	}

	m.SetSelectedFunc(func(n *cview.TreeNode) {
		ref := n.GetReference()
		entry, ok := ref.(*api.MfsLsEntry)
		if !ok {
			return
		}

		m.app.info.SetItem(entry)

	})
	return m
}
