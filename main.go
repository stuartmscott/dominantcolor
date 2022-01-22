package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/cenkalti/dominantcolor"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
)

type DominantColor struct {
	Window fyne.Window
	Image  *canvas.Image
	List   *widget.List
	Colors []color.RGBA
}

func main() {
	dc := &DominantColor{
		Window: app.New().NewWindow("Dominant Color"),
		Image:  &canvas.Image{},
		List:   &widget.List{},
	}
	dc.Image.FillMode = canvas.ImageFillContain
	dc.List.Length = func() int {
		return len(dc.Colors)
	}
	dc.List.CreateItem = func() fyne.CanvasObject {
		r := &canvas.Rectangle{}
		r.SetMinSize(fyne.NewSize(64, 64))
		t := &canvas.Text{}
		t.Alignment = fyne.TextAlignCenter
		t.Text = "#FFFFFF"
		b := &widget.Button{}
		b.Icon = theme.ContentCopyIcon()
		b.Importance = widget.LowImportance
		b.OnTapped = func() {
			dc.Window.Clipboard().SetContent(t.Text)
		}
		return container.NewBorder(nil, nil, r, b, t)
	}
	dc.List.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
		t := obj.(*fyne.Container).Objects[0].(*canvas.Text)
		r := obj.(*fyne.Container).Objects[1].(*canvas.Rectangle)
		c := dc.Colors[id]
		t.Text = dominantcolor.Hex(c)
		r.FillColor = c
		t.Refresh()
		r.Refresh()
	}
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			dialog.ShowError(err, dc.Window)
			return
		}
		dc.Open(f.Name(), f)
	}
	dc.Window.SetContent(container.NewBorder(widget.NewToolbar(widget.NewToolbarAction(theme.FileIcon(), func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, dc.Window)
				return
			}
			if reader != nil {
				dc.Open(reader.URI().Name(), reader)
			}
		}, dc.Window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".jpeg", ".png"}))
		fd.Show()
	})), nil, nil, nil, container.NewHSplit(dc.Image, dc.List)))
	dc.Window.CenterOnScreen()
	dc.Window.Resize(fyne.NewSize(800, 600))
	dc.Window.ShowAndRun()
}

func (dc *DominantColor) Open(name string, reader io.ReadCloser) {
	dc.Window.SetTitle("Dominant Color - " + name)
	defer reader.Close()
	i, _, err := image.Decode(reader)
	if err != nil {
		dialog.ShowError(err, dc.Window)
		return
	}
	dc.Image.Image = i
	dc.Image.Refresh()
	dc.Colors = dominantcolor.FindN(i, 6)
	dc.List.Refresh()
}
