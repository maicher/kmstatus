//go:build X

package out

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
// #include <unistd.h>
// #include <stdio.h>
// #include <stdlib.h>
// #include <stdarg.h>
// #include <string.h>
// #include <strings.h>
// #include <sys/time.h>
// #include <time.h>
// #include <sys/types.h>
// #include <sys/wait.h>
// #include <sys/statvfs.h>
// #include <X11/Xlib.h>
// Window getDefaultRootWindow(Display *dpy){return DefaultRootWindow(dpy);}
import "C"
import (
	"bytes"
	"fmt"
)

type Window struct {
	display     *C.Display
	defaultroot C.Window
}

func NewWindow() (*Window, error) {
	display := C.XOpenDisplay(nil)

	if display == nil {
		return nil, fmt.Errorf("Can't open XWindow display")
	}

	return &Window{
		display,
		(C.Window)(C.getDefaultRootWindow(display)),
	}, nil
}

func (w *Window) SetStatus(buffer *bytes.Buffer) {
	C.XStoreName(w.display, w.defaultroot, C.CString(buffer.String()))
	C.XSync(w.display, 0)
}
