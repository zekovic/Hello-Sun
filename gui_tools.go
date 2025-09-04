package main

import (
	"image"
	"image/color"
	"strconv"
	"strings"

	"github.com/gen2brain/iup-go/iup"
)

func getScreenSize() (int, int) {
	w_h_arr := strings.Split(iup.GetGlobal("SCREENSIZE"), "x")
	w := 0
	h := 0
	if len(w_h_arr) > 1 {
		w, _ = strconv.Atoi(w_h_arr[0])
		h, _ = strconv.Atoi(w_h_arr[1])
	}
	return w, h
}

func getMousePosition() (int, int) {
	x_y_arr := strings.Split(iup.GetGlobal("CURSORPOS"), "x")
	x := 0
	y := 0
	if len(x_y_arr) > 1 {
		x, _ = strconv.Atoi(x_y_arr[0])
		y, _ = strconv.Atoi(x_y_arr[1])
	}
	return x, y
}

func getWindowXY(ih iup.Ihandle) (int, int) {
	x_y_arr := strings.Split(ih.GetAttribute("SCREENPOSITION"), ",")
	x := 0
	y := 0
	if len(x_y_arr) > 1 {
		x, _ = strconv.Atoi(x_y_arr[0])
		y, _ = strconv.Atoi(x_y_arr[1])
	}
	return x, y
}

type WndInput struct {
	form iup.Ihandle
	label iup.Ihandle
	txtValue iup.Ihandle
	btnOk iup.Ihandle
	btnCancel iup.Ihandle
}
func (w *WndInput) Create() {
	btnAttrs := `USERSIZE="75x24", PADDING="12x8"`
	
	w.label = iup.Label("Value:")
	w.txtValue = iup.Text().SetAttributes(`USERSIZE="150x24"`)
	
	w.btnOk = iup.Button("&OK").SetCallback("ACTION", iup.ActionFunc(func(ih iup.Ihandle) int {
		w.form.SetAttribute("SIMULATEMODAL", "NO")
		iup.Hide(w.form)
		return iup.CLOSE
	})).SetAttributes(btnAttrs)
	
	w.btnCancel = iup.Button("&Cancel").SetCallback("ACTION", iup.ActionFunc(func(ih iup.Ihandle) int {
		w.form.SetAttribute("SIMULATEMODAL", "NO")
		iup.Hide(w.form)
		w.txtValue.SetAttribute("VALUE", "")
		return iup.CLOSE
	})).SetAttributes(btnAttrs)
	
	w.form = iup.Dialog(
		iup.Vbox(
			iup.Hbox(w.label, iup.Fill(), w.txtValue).SetAttributes(`CGAP=10`),
			iup.Hbox(iup.Fill(), w.btnOk, w.btnCancel).SetAttributes(`CGAP=10`),
		).SetAttributes(`MARGIN="5x5"`),
	)
}

func (w *WndInput) GetInput(parent iup.Ihandle, title, label string) string {
	w.form.SetAttribute("TITLE", title)
	w.label.SetAttribute("TITLE", label)
	w.txtValue.SetAttribute("VALUE", "")
	
	w.form.SetAttributes("SIMULATEMODAL=YES, BRINGFRONT=YES")
	w.form.SetAttribute("TOPMOST", parent.GetAttribute("TOPMOST"))
	mainX, mainY := getWindowXY(parent)
	//iup.ShowXY(w.form, max(3, main_x - 70), max(3, main_y + 50))
	iup.Popup(w.form, max(3, mainX - 70), max(3, mainY + 50))
	return w.txtValue.GetAttribute("VALUE")
}




func createRoundedImage(width, height, r int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < height * 4; i++ {
		for j := 0; j < width; j++ {
			img.Pix[i * width + j] = 255
		}
	}
	
	r_sq := (r+1) * (r+1);
	
	for i := 0; i < r; i++ {
		for j := 0; j < r; j++ {
			if (r-i)*(r-i) + (r-j)*(r-j) >= r_sq {
				img.Set(j, i, color.Transparent)
				img.Set(width - j-1, i, color.Transparent)
				img.Set(j, height - i, color.Transparent)
				img.Set(width - j-1, height - i-1, color.Transparent)
			}
		}
	}
	
	
	return img
}












