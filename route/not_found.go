package route

func (g *GenRoute) generateNotFound(name string) error {

	err := g.generateFunctionDefine("not found", name, "handler", nil, nil, nil)
	if err != nil {
		return err
	}

	g.buf.WriteFormat(`{
`)
	g.buf.WriteString(`
	err := fmt.Errorf("Not found '%s %s'", r.Method, r.URL.Path)
`)
	g.generateResponseErrorReturn("err", "404", false)

	g.buf.WriteFormat(`
}
`)
	return nil
}
