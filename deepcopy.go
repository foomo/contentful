package contentful

import (
	"bytes"
	"encoding/gob"
	"errors"
	"sync"
)

var buffers = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func DeepCopy[S, T any](dst *S, src T) error {
	buf, ok := buffers.Get().(*bytes.Buffer)
	if !ok {
		return errors.New("could not get buffer from pool")
	}
	defer func() {
		buf.Reset()
		buffers.Put(buf)
	}()

	if err := gob.NewEncoder(buf).Encode(src); err != nil {
		return err
	}

	return gob.NewDecoder(buf).Decode(dst)
}

// func DeepCopy[S, T any](dst *T, src S) error {
// 	buf := buffers.Get().(*bytes.Buffer)
// 	buf.Reset()
// 	defer buffers.Put(buf)
//
// 	enc := json.NewEncoder(buf)
// 	if err := enc.Encode(src); err != nil {
// 		return err
// 	}
//
// 	dec := json.NewDecoder(buf)
// 	return dec.Decode(dst)
// }
