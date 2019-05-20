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

	err = g.generateFunctionDefine("wrapping", name, wrap.Name, wrap.Type, wrap.Requests, nil)
	if err != nil {
		return err
	}

	g.buf.WriteString(`{
`)

	err = g.generateResponsesVar(wrap.Responses)
	if err != nil {
		return err
	}

	errName, err := g.generateResponsesErrorName(wrap.Responses)
	if err != nil {
		return err
	}

	err = g.generateCall(wrap.Name, nil, wrap.PkgPath, wrap.Type, wrap.Requests, wrap.Responses, errName)
	if err != nil {
		return err
	}

	err = g.generateResponses(wrap.Responses, "400", errName)
	if err != nil {
		return err
	}

	return
}

func (g *GenRoute) generateResponseErrorReturn(errName string, code string, noFmtErr bool) error {
	if !noFmtErr {
		for _, wrap := range g.api.Wrappings {
			if len(wrap.Responses) == 0 {
				continue
			}
			resp := wrap.Responses[0]
			if resp.Ref != "" {
				resp = g.api.Responses[resp.Ref]
			}

			// if resp.Name != errName {
			// 	continue
			// }

			name := g.getWrappingFunctionName(wrap)
			g.buf.WriteFormat(`
			%s(w, r, %s)
`, name, errName)

			return nil
		}
	}
	g.buf.AddImport("", "net/http")
	g.buf.WriteFormat(`
		http.Error(w, %s.Error(), %s)
`, errName, code)
	return nil
}

func (g *GenRoute) generateResponseError(errName string, code string, noFmtErr bool) error {
	g.buf.WriteFormat(`
	if %s != nil {`, errName)
	g.generateResponseErrorReturn(errName, code, noFmtErr)
	g.buf.WriteFormat(`
		return
	}
`)
	return nil
}
