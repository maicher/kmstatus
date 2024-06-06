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
}

func New(view *ui.View, segs *segments.Segments) *KMStatus {
	k := KMStatus{
		view:     view,
		segments: segs,
		msgQueue: make(chan any),
	}

	go k.loop()

	return &k
}

func (k *KMStatus) Refresh() {
	k.msgQueue <- refresh{}
}

func (k *KMStatus) Render() {
	k.msgQueue <- render{}
}

func (k *KMStatus) SetText(text string) {
	k.msgQueue <- setText{text: text}
}

func (k *KMStatus) Terminate() {
	close(k.msgQueue)
}

func (k *KMStatus) SetGreeting(text string) {
	b := bytes.Buffer{}
	b.WriteString(text)

	k.view.Render(&b)
}

func (k *KMStatus) loop() {
	statusBuf := &bytes.Buffer{}
	textBuf := &bytes.Buffer{}

	for msg := range k.msgQueue {
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
