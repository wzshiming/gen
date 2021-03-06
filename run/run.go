package run

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wzshiming/gen/openapi"
	"github.com/wzshiming/gen/parser"
	"github.com/wzshiming/gen/route"
	oaspec "github.com/wzshiming/openapi/spec"
)

func Run(pkgs []string, port string, way string, explode bool) error {
	for _, pkg := range pkgs {
		get(pkg)
	}

	f, err := file(pkgs, port, way, explode)
	if err != nil {
		return err
	}

	dir := os.TempDir()
	pkg := filepath.Join(dir, "gen-run")
	err = os.MkdirAll(pkg, 0755)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	const modFile = `module gen-run`

	err = ioutil.WriteFile(filepath.Join(pkg, "go.mod"), []byte(modFile), 0666)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(pkg, "main.go"), f, 0666)
	if err != nil {
		return err
	}

	get(pkg)

	cmd := exec.Command("go", "run", "./main.go")
	cmd.Dir = pkg
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func get(pkg string) {
	exec.Command("go", "get", pkg).Run()
}

func file(pkgs []string, port string, way string, explode bool) ([]byte, error) {
	def := parser.NewParser(nil)

	for _, pkg := range pkgs {
		err := def.Import(pkg, way)
		if err != nil {
			return nil, err
		}
	}

	server := "http://127.0.0.1" + port

	api, err := openapi.NewGenOpenAPI(def.API()).
		SetInfo(&oaspec.Info{
			Title:       "OpenAPI Demo",
			Description: "The current environment is only used for testing documents, and some interfaces do not work properly. Generated by gen run " + strings.Join(pkgs, " "),
			Version:     "test",
			Contact: &oaspec.Contact{
				Name: "wzshiming",
				URL:  "https://github.com/wzshiming/gen",
			},
		}).
		WithServices(server).
		SetExplode(explode).
		Generate()
	if err != nil {
		return nil, err
	}

	router, err := route.NewGenRoute(def.API()).
		WithOpenAPI(api).
		SetExplode(explode).
		SetBuildIgnore(true).
		Generate("main", ".", "Router")
	if err != nil {
		return nil, err
	}

	router.AddImport("", "net/http")
	router.AddImport("", "fmt")
	router.AddImport("", "os")
	router.AddImport("", "github.com/gorilla/handlers")

	router.WriteFormat(`

	func main() {
		mux := Router()
		mux0 := handlers.RecoveryHandler()(mux)
		mux0 = handlers.CombinedLoggingHandler(os.Stdout, mux0)
	
		fmt.Printf("Open %s/swagger/#\n")
		fmt.Printf("  or %s/swaggerui/#\n")
		fmt.Printf("  or %s/swaggereditor/#\n")
		fmt.Printf("  or %s/redoc/#\n")
		fmt.Printf("  with your browser.\n")
	
		err := http.ListenAndServe("%s", mux0)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
`, server, server, server, server, port)

	return router.Bytes(), nil
}
