package named

import (
	"fmt"
)

type Named struct {
	set map[string]map[string]string
	seg string
}

func NewNamed(seg string) *Named {
	return &Named{
		set: map[string]map[string]string{},
		seg: seg,
	}
}

func (n *Named) GetName(name string, addr string) string {
	d, ok := n.set[name]
	if !ok {
		n.set[name] = map[string]string{
			addr: name,
		}
		return name
	}

	if name, ok := d[addr]; ok {
		return name
	}

	name = fmt.Sprintf("%s%s%d", name, n.seg, len(d))
	d[addr] = name
	return name
}
