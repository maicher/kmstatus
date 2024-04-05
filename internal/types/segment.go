package types

import (
	"bytes"
)

type Segment interface {
	Refresh()
	Read(*bytes.Buffer)
}
