package route

import (
	"github.com/wzshiming/gen/spec"
)

func (g *GenRoute) generateSecurityFunction(secu *spec.Security) (err error) {
	name := g.getSecurityFunctionName(secu)

	if g.only[name] {
		return nil
	}
	g.only[name] = true

	err = g.generateFunctionDefine("security", name, secu.Name, secu.Type, nil, secu.Responses)
	if err != nil {
		return err
	}

	g.buf.WriteString(`{
`)

	err = g.generateRequestsVar(secu.Requests)
	if err != nil {
		return err
	}

	err = g.generateCallExec(secu.Name, nil, secu.PkgPath, secu.Type, secu.Requests, secu.Responses, true)
	if err != nil {
		return err
	}
	g.buf.WriteString(`
	return
}
`)

	return
}
