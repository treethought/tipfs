package ipfs

import (
	"context"
	"fmt"

	api "github.com/ipfs/go-ipfs-api"
	"gopkg.in/yaml.v2"
)

type Client struct {
	nodeURL string
	sh      *api.Shell
}

func NewClient(url string) *Client {
	c := &Client{
		nodeURL: url,
		sh:      api.NewShell(url),
	}
	return c
}

func (c *Client) GetDag(ref string) (out map[string]interface{}, err error) {
	out = make(map[string]interface{})
	err = c.sh.DagGet(ref, &out)
	return out, err
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

func (c *Client) StatEntry(entry *api.MfsLsEntry) (string, error) {

	f, err := c.sh.FilesStat(context.Background(), fmt.Sprintf("/%s", entry.Name))
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
