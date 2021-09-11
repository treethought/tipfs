package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"code.rocketnine.space/tslocum/cbind"
	"github.com/atotto/clipboard"
	api "github.com/ipfs/go-ipfs-api"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

type RepoTree struct {
	*cview.TreeView
	app          *App
	inputHandler *cbind.Configuration
}

type TreeEntry struct {
	entry *api.MfsLsEntry
	path  string
}

func (r *RepoTree) buildNodes(basePath string, entries ...*api.MfsLsEntry) []*cview.TreeNode {
	nodes := []*cview.TreeNode{}
	for _, i := range entries {
		fmt.Println(i.Name)
		node := cview.NewTreeNode(i.Name)
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
			node.SetExpanded(false)
		}

		nodes = append(nodes, node)
	}
	return nodes
}

func NewRepoTree(app *App) *RepoTree {
	m := &RepoTree{
		TreeView: cview.NewTreeView(),
		app:      app,
	}
	m.SetBorder(true)
	m.SetPadding(1, 1, 1, 1)
	m.SetTitle("repo")
	m.SetBackgroundColor(tcell.ColorDefault)
	m.SetScrollBarVisibility(cview.ScrollBarNever)

	rootNode := cview.NewTreeNode("/")
	m.SetRoot(rootNode)
	m.SetCurrentNode(rootNode)

	m.inputHandler = cbind.NewConfiguration()
	m.initBindings()

	entries, err := app.ipfs.ListFiles("/")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nodes := m.buildNodes("/", entries...)
	rootNode.SetChildren(nodes)

	return m
}

func (r *RepoTree) Update() {}

func (r *RepoTree) handleOpen(ev *tcell.EventKey) *tcell.EventKey {
	ref := r.GetCurrentNode().GetReference()
	e, ok := ref.(TreeEntry)
	if !ok {
		return nil
	}

	url := fmt.Sprintf("ipfs://%s", e.entry.Hash)
	openbrowser(url)
	return nil

}

func (r *RepoTree) handleSelect(ev *tcell.EventKey) *tcell.EventKey {
	node := r.GetCurrentNode()
	ref := node.GetReference()
	e, ok := ref.(TreeEntry)
	if !ok {
		return nil
	}

	if len(node.GetChildren()) > 0 {
		node.SetExpanded(true)
	}

	r.app.state.SetItem(e)

	return nil
}

func (t *RepoTree) initBindings() {
	t.inputHandler.SetKey(tcell.ModNone, tcell.KeyEnter, t.handleSelect)
	t.inputHandler.SetRune(tcell.ModNone, 'o', t.handleOpen)
	t.inputHandler.SetRune(tcell.ModNone, 'y', func(ev *tcell.EventKey) *tcell.EventKey {
		node := t.GetCurrentNode()
		ref := node.GetReference()
		e, ok := ref.(TreeEntry)
		if !ok {
			return nil
		}
		clipboard.WriteAll(e.entry.Hash)

		return nil
	})
	t.SetInputCapture(t.inputHandler.Capture)

}

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
