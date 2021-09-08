package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"code.rocketnine.space/tslocum/cbind"
	api "github.com/ipfs/go-ipfs-api"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

type RepoTree struct {
	*cview.TreeView
	app          *App
	inputHandler *cbind.Configuration
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
	m.SetSelectedTextColor(tcell.ColorTeal)

	rootNode := cview.NewTreeNode("/")
	m.SetRoot(rootNode)
	m.SetCurrentNode(rootNode)

	m.inputHandler = cbind.NewConfiguration()
	m.initBindings()

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
		m.app.dag.SetItem(entry)

	})
	return m
}
func (r *RepoTree) handleOpen(ev *tcell.EventKey) *tcell.EventKey {
	ref := r.GetCurrentNode().GetReference()
	entry, ok := ref.(*api.MfsLsEntry)
	if !ok {
		return nil
	}

	url := fmt.Sprintf("ipfs://%s", entry.Hash)
	openbrowser(url)
	return nil

}

func (t *RepoTree) initBindings() {
	t.inputHandler.SetRune(tcell.ModNone, 'o', t.handleOpen)
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
