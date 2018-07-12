package run

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/wzshiming/gen/openapi"
	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gotype"
	"github.com/wzshiming/openapi/util"
)

func Run(pkg string, port string, format string) error {
	f, err := file(pkg, port, format)
	if err != nil {
		return err
	}
	path := filepath.Join(os.TempDir(), "main.go")
	err = ioutil.WriteFile(path, f, 0666)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "generate", pkg)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("go", "run", path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func file(pkg string, port string, format string) ([]byte, error) {
	imp := gotype.NewImporter(gotype.WithCommentLocator())
	def := parser.NewParser(imp)
	err := def.Import(pkg)
	if err != nil {
		return nil, err
	}

	server := "http://127.0.0.1" + port

	api, err := openapi.NewGenOpenAPI(def.API()).WithServices(server).Generate()
	if err != nil {
		return nil, err
	}

	d, err := json.MarshalIndent(api, "", " ")
	if err != nil {
		return nil, err
	}
	switch format {
	case "json":

	case "yaml":
		d, err = util.JSON2YAML(d)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("undefined format %s", format)
	}

	buf := bytes.NewBuffer(nil)
	err = tpl.Execute(buf, map[string]interface{}{
		"Package": pkg,
		"Openapi": "`" + string(d) + "`",
		"Format":  format,
		"Server":  server,
		"Port":    port,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var tpl = template.Must(template.New("").Parse(temp))

const temp = `//+build ignore

package main

import (
	"net/http"
	"github.com/wzshiming/gen/ui/swaggerui"
	"github.com/urfave/negroni"
	"bytes"
	"time"
	"fmt"
	o "{{ .Package }}"
)

func main() {
	mux := &http.ServeMux{}

	mux.Handle("/", o.Router())	
	mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handle))
	mux.HandleFunc("/swagger/openapi.{{ .Format }}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "openapi.{{ .Format }}", time.Time{}, bytes.NewReader(openapi))
	})
	fmt.Printf("Open {{ .Server }}/swagger/?url=./openapi.{{ .Format }}# with your browser.\n")
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
  	n.UseHandler(mux)
	err := http.ListenAndServe("{{ .Port }}", n)
	if err != nil {
		fmt.Println(err)	
	}
	return
}

var openapi = []byte({{ .Openapi }})

`
