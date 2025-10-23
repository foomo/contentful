package contentful

import (
	"bytes"
	"encoding/json"
	"sync"
)

var buffers = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// Marshal returns the JSON encoding of v.
func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Buffer returns a new buffer and a function to return the buffer to the pool.
func Buffer() (*bytes.Buffer, func()) {
	buf := buffers.Get().(*bytes.Buffer)
	buf.Reset()
	return buf, func() {
		buffers.Put(buf)
	}
}

// Encode encodes v into JSON and writes the result to buf.
func Encode[S any](buf *bytes.Buffer, v S) error {
	return json.NewEncoder(buf).Encode(v)
}

// Decode decodes the JSON-encoded data from buf into v.
func Decode[T any](buf *bytes.Buffer, v *T) error {
	return json.NewDecoder(buf).Decode(v)
}

// DeepCopy performs a deep copy from src to dst by serializing and deserializing using a buffer.
// Returns an error if the encoding or decoding processes fail.
func DeepCopy[S, T any](dst *T, src S) error {
	buf, done := Buffer()
	defer done()

	if err := Encode[S](buf, src); err != nil {
		return err
	}

	return Decode[T](buf, dst)
}
