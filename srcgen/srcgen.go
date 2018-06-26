package srcgen

import (
	"bytes"
	"fmt"
	"unsafe"
)

type srcgen struct {
	bytes.Buffer
}

func (s *srcgen) WriteFormat(str string, a ...interface{}) error {
	_, err := s.Buffer.WriteString(fmt.Sprintf(str, a...))
	if err != nil {
		return err
	}
	return nil
}

func (s *srcgen) String() string {
	data := s.Buffer.Bytes()
	return *(*string)(unsafe.Pointer(&data))
}

func (s *srcgen) Bytes() []byte {
	return s.Buffer.Bytes()
}
