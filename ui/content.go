package ui

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Content struct {
	*tview.TextView
	app *App
}

func NewContentView(app *App) *Content {
	m := &Content{
		TextView: tview.NewTextView(),
		app:      app,
	}
	m.SetBorder(true)
	m.SetTitle("content")
	m.SetBackgroundColor(tcell.ColorDefault)
	return m
}

func (c *Content) Update() {
	current := c.app.state.currentFile
	if current.entry == nil {
		return
	}
	c.Clear()

	go func() {
		data, err := c.app.ipfs.ReadFile(current.path, current.entry)
		c.app.ui.QueueUpdateDraw(func() {
			if err != nil {
				c.SetText(fmt.Sprintf("error: %v", err))
				return
			}
			contentType, err := getFileContentType(data)
			if err != nil {
				c.SetText(err.Error())
				return
			}
			c.SetTextAlign(tview.AlignLeft)
			switch contentType {
			case "image/png", "image/jpeg":
				c.SetTextAlign(tview.AlignCenter)
				c.SetDynamicColors(true)
				_, _, w, h := c.GetRect()
				img := translateImage(bytes.NewReader(data), w, h)
				c.SetText(img)
			case "text/html; charset=utf-8":
				converter := md.NewConverter("", true, nil)
				markdown, err := converter.ConvertBytes(data)
				if err != nil {
					log.Fatal(err)
				}
				c.SetText(string(markdown))
			case "text/plain; charset=utf-8":
				c.Write(data)
			case "audio/wave", "audio/mp3", "audio/ogg":
				c.SetText(fmt.Sprintf("contentType: %s\nPress 'o' to play in browser", contentType))
			default:
				c.SetText(fmt.Sprintf("ContentType: %s\nPress 'o' to open in browser", contentType))
			}
			c.ScrollToBeginning()
		})
	}()
}

func translateImage(reader io.Reader, x, y int) string {
	img, err := buildImage(reader, x, y)
	if err != nil {
		return ""
	}
	return tview.TranslateANSI(img.Render())
}

func buildImage(reader io.Reader, x, y int) (*ansimage.ANSImage, error) {
	return ansimage.NewScaledFromReader(reader, y, x, color.Transparent, ansimage.ScaleModeFit, ansimage.NoDithering)
}

func getFileContentType(data []byte) (string, error) {
	buf := make([]byte, 512)
	n, err := bytes.NewBuffer(data).Read(buf)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buf[:n]), nil
}
