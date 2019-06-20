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

	pname := secu.Name
	if typ := secu.Type; typ != nil {
		if typ.Ref != "" {
			typ = g.api.Types[typ.Ref]
		}
		pname = typ.Name + "." + pname
	}
	err = g.generateRequestsVar(pname, secu.Requests)
	if err != nil {
		return err
	}

	errName, err := g.generateResponsesErrorName(secu.Responses)
	if err != nil {
		return err
	}

	err = g.generateCallExec(secu.Name, nil, secu.PkgPath, secu.Type, secu.Requests, secu.Responses, errName, true)
	if err != nil {
		return err
	}
	g.buf.WriteString(`
	return
}
`)

	return
}
