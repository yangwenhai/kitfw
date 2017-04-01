package protocol

import (
	"bytes"
	"io"
)

type Serializable interface {
	Save(w io.Writer) error
	Load(r io.Reader) error
}

func Encode(obj Serializable) ([]byte, error) {
	var o bytes.Buffer
	err := obj.Save(&o)
	return o.Bytes(), err
}

func Decode(obj Serializable, bin []byte) error {
	buffer := bytes.NewBuffer(bin)
	err := obj.Load(buffer)
	return err
}
