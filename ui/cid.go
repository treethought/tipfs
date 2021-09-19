package ui

import (
	"encoding/hex"
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
	b := cview.NewBox()
	b.SetBackgroundColor(tcell.ColorDefault)

	m := &CidView{
		Frame: cview.NewFrame(b),
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
		return "base58btc"
	}
	enc, _, err := multibase.Decode(c.String())
	if err == nil {
		multiBaseEncoding = multibase.EncodingToStr[enc]
	}
	return multiBaseEncoding
}

func prettyVersion(c cid.Cid) string {
	return fmt.Sprintf("cidv%d", c.Prefix().Version)
}

func codecName(c cid.Cid) string {
	// codeToStr gives protbuf
	// while multicodec lib give dad-pb

	// return cid.CodecToStr[c.Prefix().Codec]
	return multicodec.Code(c.Prefix().Codec).String()

}

func codecCode(c cid.Cid) string {
	return fmt.Sprintf("%#x", c.Prefix().Codec)
}

func multiHashCode(c cid.Cid) string {
	return fmt.Sprintf("%#x", c.Prefix().MhType)
}

func humanMultiHash(c cid.Cid) string {
	return fmt.Sprintf("%s: %d: %s",
		multihash.Codes[c.Prefix().MhType],
		c.Prefix().MhLength*8,
		hashDigest(c),
	)
}

func toV1(c cid.Cid) cid.Cid {
	if c.Version() == 1 {
		return c
	}
	return cid.NewCidV1(cid.DagProtobuf, c.Hash())
}

func hashDigest(c cid.Cid) string {
	data := []byte(c.Hash())

	dh, err := multihash.Decode(data)
	if err != nil {
		return fmt.Sprintf("error decoding multihash: %v", err)
	}

	hexout := hex.EncodeToString(dh.Digest)
	return strings.ToUpper(hexout)
}

func (c *CidView) updateText(fcid cid.Cid) *cview.TextView {
	text := cview.NewTextView()
	text.SetBackgroundColor(tcell.ColorDefault)

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

	out = append(out, "# Multibase")
	out = append(out, fmt.Sprintf("prefix: %s\n", string([]rune(fcid.String())[0])))
	out = append(out, fmt.Sprintf("name: %s\n", getMultiBaseEncoding(fcid)))

	out = append(out, "# Multicodec")
	out = append(out, fmt.Sprintf("code: %s\n", codecCode(fcid)))
	out = append(out, fmt.Sprintf("name: %s\n", codecName(fcid)))

	out = append(out, "# Multihash")
	out = append(out, fmt.Sprintf("code: %s\n", multiHashCode(fcid)))
	out = append(out, fmt.Sprintf("name: %s\n", multihash.Codes[fcid.Prefix().MhType]))
	out = append(out, fmt.Sprintf("bits: %d\n", fcid.Prefix().MhLength*8))
	out = append(out, fmt.Sprintf("digest (hex): %s\n", hashDigest(fcid)))

	out = append(out, "# CIDV1 (Base32)")
	v1 := toV1(fcid)
	out = append(out, v1.String())

	info := strings.Join(out, "\n\n")

	rendered, err := renderMarkdown(info)
	if err != nil {
		text.SetText(info)
	}
	text.SetDynamicColors(true)

	text.SetTextAlign(cview.AlignCenter)
	text.SetText(rendered)
	return text
}

func (c *CidView) Update() {
	current := c.app.state.currentHash
	fcid, err := cid.Parse(current)
	if err != nil {
		log.Fatal(err)
	}

	text := c.updateText(fcid)
	frame := cview.NewFrame(text)
	frame.SetBackgroundColor(tcell.ColorDefault)
	frame.SetBorder(true)
	// top header
	frame.AddText(fmt.Sprintf("Version: v%d", fcid.Version()), true, cview.AlignLeft, tcell.ColorWhite)
	frame.AddText(fcid.String(), true, cview.AlignCenter, tcell.ColorWhite)

	c.Frame = frame

}
