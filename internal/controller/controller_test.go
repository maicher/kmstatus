package controller

import (
	"fmt"
	"strings"
	"testing"

	"github.io/maicher/kmstatus/internal/view"
	"github.io/maicher/kmstatus/pkg/parsers/cpu"
)

var (
	ch    = make(chan any)
	outCh = make(chan string)
	v     = &View{}
	c     = NewController(ch, v)
)

type View struct {
}

func (v *View) Render(d *view.Data) {
	outCh <- fmt.Sprintf("%dkHz", d.CPU.Freq)
}
func (v *View) RenderErr(e error) {
	outCh <- fmt.Sprintf("%s", e)
}

func Test_Controller_OnData(t *testing.T) {
	go c.AggregateDataAndRenderView()
	ch <- cpu.Freq(1000)
	ch <- view.RenderView{}

	if v := <-outCh; v != "1000kHz" {
		t.Fatalf("Rendered: %s, want: 1000kHz", v)
	}
}

func Test_Controller_OnError(t *testing.T) {
	go c.AggregateDataAndRenderView()
	ch <- fmt.Errorf("This is an error")

	if v := <-outCh; !strings.Contains(v, "This is an error") {
		t.Fatalf("Error message: %s does not contain phrase 'This is an error'", v)
	}
}
