package segments

import "bytes"

type Reader interface {
	Read(*bytes.Buffer)
}
