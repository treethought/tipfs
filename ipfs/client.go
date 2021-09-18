package ipfs

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ipfs/go-cid"
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

// files api using go-ipfs-api
// ideally we'd like to use the newer go-ipfs-http-client
// because of it's nicer api and better interfaces
// however, MFS support for files api is not yet supported

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

// use go-ipfs-http-client for non files api calls

func (c *Client) GetDag(ref string) (node ipld.Node, err error) {

	refCid, err := cid.Parse(ref)
	if err != nil {
		panic(err)
	}

	node, err = c.api.Dag().Get(context.TODO(), refCid)
	if err != nil {
		return nil, err
	}
	return node, nil

}

func (c *Client) GetPeers() ([]iface.ConnectionInfo, error) {
	return c.api.Swarm().Peers(context.TODO())
}
