package route

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) generateWrappingFunction(wrap *spec.Wrapping) (err error) {
	name := g.getWrappingFunctionName(wrap)

	if g.only[name] {
		return nil
	}
	g.only[name] = true

	err = g.generateFunctionDefine("wrapping", name, wrap.Name, wrap.Type)
	if err != nil {
		return err
	}

	g.buf.WriteFormat(`(`)
	for i, resp := range wrap.Responses {
		if i != 0 {
			g.buf.WriteByte(',')
		}
		if resp.Ref != "" {
			resp = g.api.Responses[resp.Ref]
		}
		g.buf.WriteFormat("%s ", g.getVarName(resp.Name, resp.Type))
		g.Types(resp.Type)
	}
	g.buf.WriteString(`) {
`)
	err = g.generateRequestsVar(wrap.Requests, false)
	if err != nil {
		return err
	}

	err = g.generateCallExec(wrap.Name, nil, wrap.PkgPath, wrap.Type, wrap.Requests, wrap.Responses, true)
	if err != nil {
		return err
	}
	g.buf.WriteString(`
	return
}
`)

	return
}
