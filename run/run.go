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
	"github.com/wzshiming/gen/route"
	"github.com/wzshiming/gotype"
	"github.com/wzshiming/openapi/util"
)

func Run(pkgs []string, port string, format string) error {
	f, err := file(pkgs, port, format)
	if err != nil {
		return err
	}
	path := filepath.Join(os.TempDir(), "main.go")
	err = ioutil.WriteFile(path, f, 0666)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "run", path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func file(pkgs []string, port string, format string) ([]byte, error) {
	imp := gotype.NewImporter(gotype.WithCommentLocator())
	def := parser.NewParser(imp)

	for _, pkg := range pkgs {
		err := def.Import(pkg)
		if err != nil {
			return nil, err
		}
	}

	server := "http://127.0.0.1" + port

	router, err := route.NewGenRoute(def.API()).Generate("main", ".", "Router")
	if err != nil {
		return nil, err
	}
	router.AddImport("", "net/http")
	router.AddImport("", "fmt")
	router.AddImport("", "github.com/urfave/negroni")
	router.AddImport("", "github.com/wzshiming/gen/ui/swaggerui")

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
		"Openapi": "`" + string(d) + "`",
		"Format":  format,
		"Server":  server,
		"Port":    port,
		"Router":  router,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var tpl = template.Must(template.New("").Parse(temp))

const temp = `//+build ignore

{{ .Router }}

func main() {
	mux := &http.ServeMux{}
	mux.Handle("/", Router())
	mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.HandleWithFile("openapi.{{ .Format }}", openapi)))
	fmt.Printf("Open {{ .Server }}/swagger/?url=openapi.{{ .Format }}# with your browser.\n")
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
