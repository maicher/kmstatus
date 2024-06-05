package ipc

const (
	Refresh   = "r"
	SetText   = "s"
	UnsetText = "u"
)

type Cmd struct {
	Name    string `json:"n"`
	Payload string `json:"p"`
}

func NewRefreshCmd() *Cmd {
	return &Cmd{Name: Refresh}
}

func NewSetTextCmd(text string) *Cmd {
	return &Cmd{Name: SetText, Payload: text}
}

func NewUnsetTextCmd() *Cmd {
	return &Cmd{Name: UnsetText}
}
