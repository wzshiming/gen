package route

import (
	"bytes"
	"strings"

	"github.com/wzshiming/gen/model"
	"github.com/wzshiming/gen/spec"
)

// GenClient is the generating generating
type GenRoute struct {
	api *spec.API
	buf *bytes.Buffer
	model.GenModel
}

func NewGenRoute(api *spec.API) *GenRoute {
	buf := bytes.NewBuffer(nil)
	return &GenRoute{
		api:      api,
		buf:      buf,
		GenModel: *model.NewGenModel(api, buf),
	}
}

func (g *GenRoute) Generate() ([]byte, error) {
	g.buf.WriteString(`// Code generated; DO NOT EDIT.
package ` + g.api.Package + `

import (
	"encoding/json"
	"github.com/gorilla/mux"
)

`)
	err := g.GenerateRoutes()
	if err != nil {
		return nil, err
	}

	return g.buf.Bytes(), nil
}

func (g *GenRoute) GenerateRoutes() (err error) {
	// route := mux.NewRouter()
	g.buf.WriteString(`
func Router() *mux.Router {
	router := mux.NewRouter()
`)
	for _, v := range g.api.Operations {
		g.buf.WriteString(`
	router.Path("` + v.Path + `").
	Methods("` + strings.ToUpper(v.Method) + `").
`)

		for _, req := range v.Requests {
			if req.Ref != "" {
				req = g.api.Requests[req.Ref]
			}
			switch req.In {
			case "query":
				g.buf.WriteString(`Queries("` + req.Name + `").`)
			case "header":
				g.buf.WriteString(`Headers("` + req.Name + `").`)
			}
		}

		g.buf.WriteString(`
HandlerFunc(func(_w http.ResponseWriter, _r *http.Request) {
	_vars := mux.Vars(_r)
`)

		for i, resp := range v.Responses {
			if resp.Ref != "" {
				resp = g.api.Responses[resp.Ref]
			}
			if i != 0 {
				g.buf.WriteByte(',')
			}
			g.buf.WriteString(resp.Name)
		}

		g.buf.WriteString(":=")
		g.buf.WriteString(v.Name)

		g.buf.WriteString("(")
		for i, req := range v.Requests {
			if req.Ref != "" {
				req = g.api.Requests[req.Ref]
			}
			if i != 0 {
				g.buf.WriteByte(',')
			}
			switch req.In {
			case "body":
				err := g.UnmarshalContentBody(req)
				if err != nil {
					return err
				}

			case "cookie":
				// TODO
				g.buf.WriteString("nil")
			default:
				g.buf.WriteString(`_vars["` + req.Name + `"]`)
			}
		}
		g.buf.WriteString(")")

		for _, resp := range v.Responses {
			if resp.Ref != "" {
				resp = g.api.Responses[resp.Ref]
			}

			g.MarshalContentBody(resp)

			//g.buf.WriteString(resp.Name)
		}

		g.buf.WriteString(`
})
`)

	}
	g.buf.WriteString(`
	return router
}
`)
	return
}

func (g *GenRoute) UnmarshalContentBody(req *spec.Request) error {
	g.buf.WriteString(`func()(_b `)
	g.Types(req.Type)
	g.buf.WriteString(`){
	data, err := ioutil.ReadAll(_body)
	if err != nil {
		return
	}
	_body.Close()
	`)
	switch req.Content {
	case "json":
		g.buf.WriteString(`json.Unmarshal(data,&_b)`)
	case "xml":
		g.buf.WriteString(`xml.Unmarshal(data,&_b)`)
	}
	g.buf.WriteString(`
	return
}()`)
	return nil
}

func (g *GenRoute) MarshalContentBody(resp *spec.Response) error {
	g.buf.WriteString(`
	if ` + resp.Name + ` != nil {`)

	contentType := ""
	switch resp.Content {
	case "json":
		contentType = "application/json; charset=utf-8"
		g.buf.WriteString(`
	data, err:=json.Marshal(_b)
	if err!=nil {
		_w.WriteHeader(500)
		_w.Write([]byte(err.Error()))
		return
	}`)

	case "xml":
		contentType = "application/xml; charset=utf-8"
		g.buf.WriteString(`
	data,err:=xml.Marshal(_b)
	if err!=nil {
		_w.WriteHeader(500)
		_w.Write([]byte(err.Error()))
		return
	}`)
	default:
		contentType = "text/plain; charset=utf-8"

	}
	g.buf.WriteString(`
	_w.Header().Set("Content-Type","` + contentType + `")
	_w.WriteHeader(` + resp.Code + `)
	_w.Write([]byte(err.Error()))
	return
`)
	g.buf.WriteString(`}`)
	return nil
}
