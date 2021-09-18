package ui

import (
	"fmt"
	"log"
	"strings"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"
)

type CIDInfo struct {
	*cview.Flex
	app *App
}

type CidView struct {
	*cview.Frame
	app *App
}

func NewCIDView(app *App) *CidView {

	m := &CidView{
		Frame: cview.NewFrame(cview.NewBox()),
		app:   app,
	}

	m.SetBorder(true)
	m.SetPadding(1, 1, 1, 1)
	m.SetTitle("cid")
	m.SetBackgroundColor(tcell.ColorDefault)
	return m
}

func getMultiBaseEncoding(c cid.Cid) string {
	multiBaseEncoding := "unknown"
	if c.Version() == 0 {
		multiBaseEncoding = "base58btc"
	} else {
		enc, _, err := multibase.Decode(c.String())
		if err == nil {
			multiBaseEncoding = multibase.EncodingToStr[enc]
		}
	}
	return multiBaseEncoding
}

func prettyVersion(c cid.Cid) string {
	return fmt.Sprintf("cidv%d", c.Prefix().Version)
}

func codecName(c cid.Cid) string {
	// codeToStr gives protbuf
	// while multicodec lib give dad-ob

	// return cid.CodecToStr[c.Prefix().Codec]
	return multicodec.Code(c.Prefix().Codec).String()

}

// TODO figure out how to get actual code
func codecCode(c cid.Cid) string {
	s := fmt.Sprintf("%v", c.Prefix().Codec)
	code := multicodec.Code(c.Prefix().Codec)
	_ = code.Set(s)
	// return string(code)

	return fmt.Sprintf("%v", uint64(code))

	// s := strconv.FormatUint(c.Prefix().Codec, 64)
}

func humanMultiHash(c cid.Cid) string {
	h := c.Hash()
	hash := &h
	return fmt.Sprintf("%s: %d: %s",
		multihash.Codes[c.Prefix().MhType],
		c.Prefix().MhLength*8,
		strings.ToUpper(string(hash.HexString())),
	)
}

func (c *CidView) Update() {
	current := c.app.state.currentFile
	_, entry := current.path, current.entry

	fcid, err := cid.Parse(entry.Hash)
	if err != nil {
		log.Fatal(err)
	}

	text := cview.NewTextView()
	text.SetBackgroundColor(tcell.ColorDefault)

	frame := cview.NewFrame(text)
	frame.SetBackgroundColor(tcell.ColorDefault)
	// top header
	frame.AddText(fmt.Sprintf("Version: v%d", fcid.Version()), true, cview.AlignLeft, tcell.ColorDefault)
	frame.AddText(fcid.String(), true, cview.AlignCenter, tcell.ColorDefault)

	c.Frame = frame

	out := []string{}
	out = append(out, "# Human Readable CID")

	human := fmt.Sprintf("%v - %v - %v - %s\n",
		getMultiBaseEncoding(fcid),
		prettyVersion(fcid),
		codecName(fcid),
		humanMultiHash(fcid),
	)
	out = append(out, human)

	out = append(out, "multibase - version - multicodec - multihash\n")
	out = append(out, "---")

	out = append(out, "# Multibase")
	out = append(out, fmt.Sprintf("code: %s", codecCode(fcid)))
	out = append(out, fmt.Sprintf("name: %s", codecName(fcid)))

	info := strings.Join(out, "\n\n")

	rendered, err := renderMarkdown(info)
	if err != nil {
		text.SetText(info)
	}
	text.SetDynamicColors(true)

	text.SetText(rendered)

}
