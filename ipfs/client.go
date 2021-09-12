package ipfs

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	api "github.com/ipfs/go-ipfs-api"
	ipfs "github.com/ipfs/go-ipfs-http-client"
	ipld "github.com/ipfs/go-ipld-format"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"gopkg.in/yaml.v2"
)

type DagData struct {
	Data  string
	Links []ipld.Link
}

type Client struct {
	nodeURL string
	sh      *api.Shell
	api     *ipfs.HttpApi
}

func NewClient(url string) *Client {
	addr, err := ipfs.ApiAddr("~/.ipfs")
	if err != nil {
		log.Fatal("failed to read ipfs api file")
	}
	fmt.Println("connected to api at: ", addr.String())
	iapi, err := ipfs.NewApi(addr)
	if err != nil {
		log.Fatal(err)
	}

	c := &Client{
		nodeURL: url,
		sh:      api.NewLocalShell(),
		api:     iapi,
	}
	return c
}

// func (c *Client) Get(path string, entry *api.MfsLsEntry) ([]byte, error) {
// }

func (c *Client) ReadFile(path string, entry *api.MfsLsEntry) ([]byte, error) {

	if entry.Type == api.TDirectory {
		return []byte("directory"), nil
	}
	r, err := c.sh.FilesRead(context.Background(), path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

func (c *Client) GetDag(ref string) (dag *DagData, err error) {
	dag = &DagData{}

	err = c.sh.DagGet(ref, dag)
	if err != nil {
		return nil, err
	}
	return dag, nil
}

func (c *Client) ListFiles(path string) (entries []*api.MfsLsEntry, err error) {
	entries, err = c.sh.FilesLs(context.Background(), path, api.FilesLs.Stat(true))
	if err != nil {
		fmt.Println(err)
		return entries, err
	}
	return entries, nil
}

func (c *Client) StatFile(path string, entry *api.MfsLsEntry) (string, error) {

	f, err := c.sh.FilesStat(context.Background(), path)
	if err != nil {
		return "", err
	}

	statOut, err := yaml.Marshal(f)
	if err != nil {
		return "", err
	}

	out := string(statOut)

	if entry.Type == api.TDirectory {
		children, err := c.sh.FilesLs(context.Background(), fmt.Sprintf("/%s", entry.Name))
		if err != nil {
			return "", err
		}

		out = fmt.Sprintf("%s\nchildren: %d", out, len(children))
	}

	return out, nil

}

func (c *Client) GetPeers() ([]iface.ConnectionInfo, error) {
	return c.api.Swarm().Peers(context.TODO())
}

func (c *Client) GetPeer(p string) (*api.PeerInfo, error) {
	return c.sh.FindPeer(p)
}
