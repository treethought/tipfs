package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/atotto/clipboard"
	api "github.com/ipfs/go-ipfs-api"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type RepoTree struct {
	*tview.TreeView
	app *App
}

type TreeEntry struct {
	entry *api.MfsLsEntry
	path  string
}

func (r *RepoTree) buildNodes(basePath string, entries ...*api.MfsLsEntry) []*tview.TreeNode {
	nodes := []*tview.TreeNode{}
	for _, i := range entries {
		node := tview.NewTreeNode(i.Name)
		ref := TreeEntry{
			entry: i,
			path:  filepath.Join(basePath, i.Name),
		}
		node.SetReference(ref)

		if i.Type == api.TDirectory {
			children, err := r.app.ipfs.ListFiles(ref.path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			childrenNodes := r.buildNodes(ref.path, children...)
			node.SetChildren(childrenNodes)
			node.Collapse()
		}

		nodes = append(nodes, node)
	}
	return nodes
}

func NewRepoTree(app *App) *RepoTree {
	m := &RepoTree{
		TreeView: tview.NewTreeView(),
		app:      app,
	}
	m.SetBorder(true)
	m.SetTitle("repo")
	m.SetBackgroundColor(tcell.ColorDefault)

	rootNode := tview.NewTreeNode("/")
	m.SetRoot(rootNode)
	m.SetCurrentNode(rootNode)

	entries, err := app.ipfs.ListFiles("/")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootNode.SetChildren(m.buildNodes("/", entries...))

	m.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		e, ok := ref.(TreeEntry)
		if !ok {
			return
		}
		if len(node.GetChildren()) > 0 {
			node.Expand()
		}
		m.app.state.SetItem(e)
	})

	m.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		node := m.GetCurrentNode()
		if node == nil {
			return ev
		}
		ref := node.GetReference()
		e, ok := ref.(TreeEntry)
		if !ok {
			return ev
		}
		switch ev.Rune() {
		case 'o':
			openbrowser(fmt.Sprintf("ipfs://%s", e.entry.Hash))
			return nil
		case 'y':
			clipboard.WriteAll(e.entry.Hash)
			return nil
		}
		return ev
	})

	return m
}

func (r *RepoTree) Update() {}

func openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
