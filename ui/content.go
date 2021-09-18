package ui

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"

	"code.rocketnine.space/tslocum/cview"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/glamour"
	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/gdamore/tcell/v2"
)

type Content struct {
	*cview.TextView

	app *App
}

func NewContentView(app *App) *Content {
	m := &Content{
		TextView: cview.NewTextView(),
		app:      app,
	}

	m.SetBorder(true)
	m.SetPadding(1, 1, 1, 1)
	m.SetTitle("content")
	m.SetBackgroundColor(tcell.ColorDefault)

	return m
}

func (c *Content) Update() {
	current := c.app.state.currentFile
	path, entry := current.path, current.entry
	c.Clear()

	data, err := c.app.ipfs.ReadFile(path, entry)
	if err != nil {
		panic(err)
	}
	contentType, err := getFileContentType(data)
	if err != nil {
		c.SetText(err.Error())
	}

	c.SetTextAlign(cview.AlignLeft)

	switch contentType {
	case "image/png", "image/jpeg":

		c.SetTextAlign(cview.AlignCenter)
		c.SetDynamicColors(true)
		_, _, w, h := c.GetRect()
		r := bytes.NewReader(data)
		img := translateImage(r, w, h)
		c.SetText(img)
	case "text/html; charset=utf-8":
		c.SetDynamicColors(true)
		converter := md.NewConverter("", true, nil)
		markdown, err := converter.ConvertBytes(data)
		if err != nil {
			log.Fatal(err)
		}

		r, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
		)

		rendered, err := r.Render(string(markdown))
		if err != nil {
			c.SetText(err.Error())
		}
		trans := cview.TranslateANSI(rendered)
		c.SetText(trans)

	case "text/plain; charset=utf-8":
		c.Write(data)

	case "audio/wave", "audio/mp3", "audio/ogg":
		c.SetText(fmt.Sprintf("contentType: %s\nPress 'o' to play in browser", contentType))

	default:
		c.SetText(fmt.Sprintf("ContentType: %s\nPress 'o' to open in browser", contentType))
	}

	c.ScrollToBeginning()

}

func translateImage(reader io.Reader, x, y int) string {
	img, err := buildImage(reader, x, y)
	if err != nil {
		return ""
	}
	ansi := img.Render()
	return cview.TranslateANSI(ansi)

}

func buildImage(reader io.Reader, x, y int) (*ansimage.ANSImage, error) {
	pix, err := ansimage.NewScaledFromReader(reader, y, x, color.Transparent, ansimage.ScaleModeFit, ansimage.NoDithering)
	if err != nil {
		return nil, err
	}
	return pix, nil

}

func getFileContentType(data []byte) (string, error) {

	s := bytes.NewBuffer(data)

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := s.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
