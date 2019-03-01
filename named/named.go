package named

import (
	"fmt"
)

type Named struct {
	set map[string]map[string]string
	seg string
	sub map[string]*Named
}

func NewNamed(seg string) *Named {
	return &Named{
		set: map[string]map[string]string{},
		seg: seg,
		sub: map[string]*Named{},
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

func (n *Named) GetSubNamed(addr string) *Named {
	name, ok := n.sub[addr]
	if ok {
		return name
	}
	name = NewNamed(n.seg)
	n.sub[addr] = name
	return name
}
