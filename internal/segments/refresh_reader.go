package segments

import "bytes"

type RefreshReader interface {
	Refresh()
	Read(*bytes.Buffer)
}
