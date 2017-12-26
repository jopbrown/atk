// Copyright 2017 visualfc. All rights reserved.

package tk

import (
	"fmt"
)

type WindowInfo struct {
	X      int
	Y      int
	Width  int
	Height int
}

var (
	globalWindowInfoMap = make(map[string]*WindowInfo)
)

type Window struct {
	BaseWidget
}

func (w *Window) Type() string {
	return "Window"
}

func (w *Window) SetTitle(title string) *Window {
	eval(fmt.Sprintf("wm title %v {%v}", w.id, title))
	return w
}

func (w *Window) Title() string {
	s, _ := evalAsString(fmt.Sprintf("wm title %v", w.id))
	return s
}

func (w *Window) SetAlpha(alpha float64) *Window {
	eval(fmt.Sprintf("wm attributes %v -alpha {%v}", w.id, alpha))
	return w
}

func (w *Window) Alpha() float64 {
	r, _ := evalAsFloat64(fmt.Sprintf("wm attributes %v -alpha", w.id))
	return r
}

func (w *Window) SetFullScreen(full bool) *Window {
	eval(fmt.Sprintf("wm attributes %v -fullscreen %v", w.id, boolToInt(full)))
	return w
}

func (w *Window) IsFullScreen() bool {
	r, _ := evalAsBool(fmt.Sprintf("wm attributes %v -fullscreen", w.id))
	return r
}

func (w *Window) SetTopmost(full bool) *Window {
	eval(fmt.Sprintf("wm attributes %v -topmost %v", w.id, boolToInt(full)))
	return w
}

func (w *Window) IsTopmost() bool {
	r, _ := evalAsBool(fmt.Sprintf("wm attributes %v -topmost", w.id))
	return r
}

func (w *Window) SetGeometryN(x int, y int, width int, height int) *Window {
	globalWindowInfoMap[w.id] = &WindowInfo{x, y, width, height}
	eval(fmt.Sprintf("wm geometry %v %vx%v+%v+%v", w.id, width, height, x, y))
	return w
}

func (w *Window) SetGeometry(v Geometry) *Window {
	return w.SetGeometryN(v.X, v.Y, v.Width, v.Height)
}

func (w *Window) GeometryN() (x int, y int, width int, height int) {
	if !w.IsVisible() {
		if info, ok := globalWindowInfoMap[w.id]; ok {
			return info.X, info.Y, info.Width, info.Height
		}
	}
	s, err := evalAsString(fmt.Sprintf("update\nwm geometry %v", w.id))
	if err != nil {
		return
	}
	var ar []*int = []*int{&width, &height, &x, &y}
	var n *int = ar[0]
	var index int
	for _, r := range s {
		if r == 'x' || r == '+' {
			index++
			n = ar[index]
		} else {
			*n = *n*10 + int(r-'0')
		}
	}
	return
}

func (w *Window) Geometry() Geometry {
	x, y, width, height := w.GeometryN()
	return Geometry{x, y, width, height}
}

func (w *Window) MoveN(x int, y int) *Window {
	return w.SetPosN(x, y)
}

func (w *Window) Move(pos Pos) *Window {
	return w.SetPosN(pos.X, pos.Y)
}

func (w *Window) SetPosN(x int, y int) *Window {
	globalWindowInfoMap[w.id].X = x
	globalWindowInfoMap[w.id].Y = y
	eval(fmt.Sprintf("wm geometry %v +%v+%v", w.id, x, y))
	return w
}

func (w *Window) SetPos(pos Pos) *Window {
	return w.SetPosN(pos.X, pos.Y)
}

func (w *Window) PosN() (x int, y int) {
	x, y, _, _ = w.GeometryN()
	return
}

func (w *Window) Pos() Pos {
	x, y, _, _ := w.GeometryN()
	return Pos{x, y}
}

func (w *Window) ResizeN(width int, height int) *Window {
	return w.SetSizeN(width, height)
}

func (w *Window) Resize(sz Size) *Window {
	return w.SetSizeN(sz.Width, sz.Height)
}

func (w *Window) SetSizeN(width int, height int) *Window {
	globalWindowInfoMap[w.id].Width = width
	globalWindowInfoMap[w.id].Height = height
	eval(fmt.Sprintf("wm geometry %v %vx%v", w.id, width, height))
	return w
}

func (w *Window) SetSize(sz Size) *Window {
	return w.SetSizeN(sz.Width, sz.Height)
}

func (w *Window) SizeN() (width int, height int) {
	_, _, width, height = w.GeometryN()
	return
}

func (w *Window) Size() Size {
	_, _, width, height := w.GeometryN()
	return Size{width, height}
}

func (w *Window) SetWidth(width int) *Window {
	_, _, _, height := w.GeometryN()
	w.SetSizeN(width, height)
	return w
}

func (w *Window) Width() (width int) {
	_, _, width, _ = w.GeometryN()
	return
}

func (w *Window) SetHeight(height int) *Window {
	_, _, width, _ := w.GeometryN()
	w.SetSizeN(width, height)
	return w
}

func (w *Window) Height() (height int) {
	_, _, _, height = w.GeometryN()
	return
}

func (w *Window) SetNaturalSize() *Window {
	eval(fmt.Sprintf("wm geometry %v {}", w.id))
	return w
}

func (w *Window) SetResizable(enableWidth bool, enableHeight bool) *Window {
	eval(fmt.Sprintf("wm resizable %v %v %v", w.id, boolToInt(enableWidth), boolToInt(enableHeight)))
	return w
}

func (w *Window) IsResizable() (enableWidth bool, enableHeight bool) {
	s, err := evalAsString(fmt.Sprintf("wm resizable %v", w.id))
	if err == nil {
		n1, n2 := parserTwoInt(s)
		enableWidth = n1 != 0
		enableHeight = n2 != 0
	}
	return
}

func (w *Window) Iconify() *Window {
	eval(fmt.Sprintf("wm iconify %v", w.id))
	return w
}

func (w *Window) IsIconify() bool {
	r, _ := evalAsString(fmt.Sprintf("wm state %v", w.id))
	return r == "iconic"
}

func (w *Window) ShowNormal() *Window {
	if w.IsFullScreen() {
		w.SetFullScreen(false)
	}
	eval(fmt.Sprintf("wm state %v normal", w.id))
	return w
}

func (w *Window) ShowFullScreen() *Window {
	return w.SetFullScreen(true)
}

func (w *Window) ShowMinimized() *Window {
	return w.Iconify()
}

func (w *Window) IsMinimized() bool {
	return w.IsIconify()
}

func (w *Window) Hide() *Window {
	eval(fmt.Sprintf("wm state %v withdrawn", w.id))
	return w
}

func (w *Window) IsVisible() bool {
	s, _ := evalAsString(fmt.Sprintf("wm state %v", w.id))
	return s != "withdrawn"
}

func (w *Window) SetVisible(b bool) *Window {
	if w.IsVisible() != b {
		if b {
			w.ShowNormal()
		} else {
			w.Hide()
		}
	}
	return w
}

func (w *Window) Deiconify() *Window {
	eval(fmt.Sprintf("wm deiconify %v", w.id))
	return w
}

func (w *Window) SetMaximumSizeN(width int, height int) *Window {
	eval(fmt.Sprintf("wm maxsize %v %v %v", w.id, width, height))
	return w
}

func (w *Window) SetMaximumSize(sz Size) *Window {
	return w.SetMaximumSizeN(sz.Width, sz.Height)
}

func (w *Window) MaximumSizeN() (int, int) {
	s, _ := evalAsString(fmt.Sprintf("wm maxsize %v", w.id))
	return parserTwoInt(s)
}

func (w *Window) MaximumSize() Size {
	width, height := w.MaximumSizeN()
	return Size{width, height}
}

func (w *Window) SetMinimumSizeN(width int, height int) *Window {
	eval(fmt.Sprintf("wm minsize %v %v %v", w.id, width, height))
	return w
}

func (w *Window) SetMinimumSize(sz Size) *Window {
	return w.SetMinimumSizeN(sz.Width, sz.Height)
}

func (w *Window) MinimumSizeN() (int, int) {
	s, _ := evalAsString(fmt.Sprintf("wm minsize %v", w.id))
	return parserTwoInt(s)
}

func (w *Window) MinimumSize() Size {
	width, height := w.MinimumSizeN()
	return Size{width, height}
}

func (w *Window) ScreenSizeN() (width int, height int) {
	width, _ = evalAsInt(fmt.Sprintf("winfo screenwidth %v", w.id))
	height, _ = evalAsInt(fmt.Sprintf("winfo screenheight %v", w.id))
	return
}

func (w *Window) ScreenSize() Size {
	width, height := w.ScreenSizeN()
	return Size{width, height}
}

func (w *Window) Center() *Window {
	sw, sh := w.ScreenSizeN()
	width, height := w.SizeN()
	x := (sw - width) / 2
	y := (sh - height) / 2
	return w.MoveN(x, y)
}

func (w *Window) OnClose(fn func()) error {
	fnName, err := mainInterp.CreateActionByType("window_close", func() {
		if fn != nil {
			fn()
		} else {
			w.Destroy()
		}
	})
	if err != nil {
		return err
	}
	return eval(fmt.Sprintf("wm protocol %v WM_DELETE_WINDOW %v", w.id, fnName))
}

func (w *Window) registerWindowInfo() {
	globalWindowInfoMap[w.id] = &WindowInfo{0, 0, 200, 200}
}

func MainWindow() *Window {
	return mainWindow
}

func NewWindow(id string) *Window {
	id = MakeWidgetId(nil, id)
	err := eval(fmt.Sprintf("tk::toplevel %v", id))
	if err != nil {
		return nil
	}
	w := &Window{}
	w.SetInternalId(id)
	w.Hide()
	w.registerWindowInfo()
	RegisterWidget(w)
	return w
}