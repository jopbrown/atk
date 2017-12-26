// Copyright 2017 visualfc. All rights reserved.

package tk

import "fmt"

// NOTE: update must
func (w *Window) ShowMaximized() *Window {
	eval(fmt.Sprintf("update\nwm state %v zoomed", w.id))
	return w
}

func (w *Window) IsMaximized() bool {
	r, _ := evalAsString(fmt.Sprintf("wm state %v", w.id))
	return r == "zoomed"
}