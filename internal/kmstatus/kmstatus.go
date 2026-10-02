package kmstatus

import (
	"bytes"

	"github.com/maicher/kmstatus/internal/segments"
	"github.com/maicher/kmstatus/internal/ui"
)

type refresh struct{}
type render struct{}
type setText struct{ text string }

type KMStatus struct {
	view     *ui.View
	segments *segments.Segments

	msgQueue chan any
	done     chan struct{}
}

func New(view *ui.View, segs *segments.Segments) *KMStatus {
	k := KMStatus{
		view:     view,
		segments: segs,
		msgQueue: make(chan any),
		done:     make(chan struct{}),
	}

	go k.loop()

	return &k
}

func (k *KMStatus) Refresh() {
	k.send(refresh{})
}

func (k *KMStatus) Render() {
	k.send(render{})
}

func (k *KMStatus) SetText(text string) {
	k.send(setText{text: text})
}

// Terminate stops the loop. Messages sent afterwards are dropped,
// so a command handled during shutdown does not panic.
func (k *KMStatus) Terminate() {
	close(k.done)
}

func (k *KMStatus) send(msg any) {
	select {
	case k.msgQueue <- msg:
	case <-k.done:
	}
}

func (k *KMStatus) SetGreeting(text string) {
	b := bytes.Buffer{}
	b.WriteString(text)

	k.view.Render(&b)
}

func (k *KMStatus) loop() {
	statusBuf := &bytes.Buffer{}
	textBuf := &bytes.Buffer{}

	for {
		var msg any

		select {
		case msg = <-k.msgQueue:
		case <-k.done:
			return
		}

		switch msg := msg.(type) {
		case refresh:
			k.segments.Refresh()
		case render:
			statusBuf.Write(textBuf.Bytes())
			k.segments.Read(statusBuf)
			k.view.Render(statusBuf)
			statusBuf.Reset()
		case setText:
			textBuf.Reset()
			textBuf.WriteString(msg.text)
		}
	}
}
