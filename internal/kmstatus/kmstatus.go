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
	textBuf  *bytes.Buffer
	buf      *bytes.Buffer
	view     *ui.View
	segments *segments.Segments

	msgQueue chan any
}

func New(view *ui.View, segs *segments.Segments) *KMStatus {
	k := KMStatus{
		textBuf:  &bytes.Buffer{},
		buf:      &bytes.Buffer{},
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

	k.view.Flush(&b)
}

func (k *KMStatus) loop() {
	for msg := range k.msgQueue {
		k.handleMessage(msg)
	}
}

func (k *KMStatus) handleMessage(msg any) {
	switch msg := msg.(type) {
	case refresh:
		k.segments.Refresh()
	case render:
		k.buf.WriteString(k.textBuf.String())
		k.segments.Read(k.buf)
		k.view.Flush(k.buf)
	case setText:
		k.textBuf.Reset()
		k.textBuf.WriteString(msg.text)
	}
}
