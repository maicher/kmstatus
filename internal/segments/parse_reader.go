package segments

import "bytes"

type ParseReader interface {
	Parse()
	Read(*bytes.Buffer)
}
