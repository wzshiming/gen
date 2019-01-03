package named

import (
	"fmt"
	"unsafe"
)

type Named struct {
	set map[string]map[unsafe.Pointer]string
}

func NewNamed() *Named {
	return &Named{
		set: map[string]map[unsafe.Pointer]string{},
	}
}

func (n *Named) GetName(name string, addr unsafe.Pointer) string {
	d, ok := n.set[name]
	if !ok {
		n.set[name] = map[unsafe.Pointer]string{
			addr: name,
		}
		return name
	}

	if name, ok := d[addr]; ok {
		return name
	}

	name = fmt.Sprintf("%s_%d", name, len(d))
	d[addr] = name
	return name
}
