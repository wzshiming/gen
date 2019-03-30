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

	g.buf.AddImport("", "net/http")
	g.buf.WriteFormat(`
// %s Is the security of %s
func %s(`, name, secu.Name, name)
	if secu.Type != nil {
		g.buf.WriteString("s ")
		g.PtrTypes(secu.Type)
		g.buf.WriteString(", ")
	}
	g.buf.WriteFormat(`w http.ResponseWriter, r *http.Request) (`)

	for i, resp := range secu.Responses {
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

	err = g.generateRequestsVar(secu.Requests, false)
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
