package ipfs

import (
	"context"
	"fmt"
	"io/ioutil"

	api "github.com/ipfs/go-ipfs-api"
	ipld "github.com/ipfs/go-ipld-format"
	"gopkg.in/yaml.v2"
)

type DagData struct {
	Data  string
	Links []ipld.Link
}

type Client struct {
	nodeURL string
	sh      *api.Shell
}

func NewClient(url string) *Client {

	c := &Client{
		nodeURL: url,
		sh:      api.NewLocalShell(),
	}
	return c
}

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
	for _, e := range entries {
		fmt.Println(e.Name, e.Hash, e.Size, e.Type)
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

func (c *Client) GetPeers() (*api.SwarmConnInfos, error) {
	return c.sh.SwarmPeers(context.TODO())
}
