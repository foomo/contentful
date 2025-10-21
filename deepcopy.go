package contentful

import (
	"bytes"
	"encoding/json"
	"sync"
)

var buffers = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func DeepCopy[S, T any](dst *T, src S) error {
	buf := buffers.Get().(*bytes.Buffer)
	buf.Reset()
	defer buffers.Put(buf)

	enc := json.NewEncoder(buf)
	if err := enc.Encode(src); err != nil {
		return err
	}

	dec := json.NewDecoder(buf)
	return dec.Decode(dst)
}
