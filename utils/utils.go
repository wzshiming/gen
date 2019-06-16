package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"go/build"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
)

func PackageOmitted(pkgpath string) []string {
	dirs := build.Default.SrcDirs()
	pkgpath = strings.TrimSuffix(pkgpath, "/...")

	isLocal := strings.HasPrefix(pkgpath, "./")

	dir := ""
	for _, d := range dirs {
		fi, err := os.Stat(path.Join(d, pkgpath))
		if err == nil && fi.IsDir() {
			dir = d
			break
		}
	}

	pkgs := []string{pkgpath}
	outs := []string{}
	for i := 0; i != len(pkgs); i++ {
		pkg := pkgs[i]
		fis, err := ioutil.ReadDir(path.Join(dir, pkg))
		if err != nil {
			continue
		}

		curr := true
		for _, fi := range fis {
			name := fi.Name()
			if fi.IsDir() {
				pkgs = append(pkgs, path.Join(pkg, name))
			} else if curr && strings.HasSuffix(name, ".go") {
				if isLocal {
					pkg = "./" + pkg
				}
				outs = append(outs, pkg)
				curr = false
			}
		}
	}
	return outs
}

func GetPackagePath(path string) string {
	const (
		src    = "/src/"
		vendor = "/vendor/"
	)
	if index := strings.LastIndex(path, vendor); index != -1 {
		return path[index+len(vendor):]
	}
	if index := strings.Index(path, src); index != -1 {
		return path[index+len(src):]
	}
	return strings.TrimLeft(path, "/")
}

func MergeLine(t string) string {
	return strings.TrimSpace(strings.Replace(t, "\n", " ", -1))
}

func CommentLine(t string) string {
	t = strings.TrimSpace(t)
	if t == "" {
		return "\n"
	}
	return "// " + strings.Join(strings.Split(t, "\n"), "\n// ") + "\n"
}

func Hash(s ...string) string {
	h := md5.New()
	for _, v := range s {
		h.Write([]byte(v))
	}
	return hex.EncodeToString(h.Sum(nil)[:])
}

// GetTag [#[^#]+#]...
func GetTag(text string) (string, reflect.StructTag) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", ""
	}
	ss := []string{}
	other := bytes.NewBuffer(nil)
	for _, text := range strings.Split(text, "\n") {
		if text == "" {
			continue
		}
		prev := -1
		for i, v := range text {
			if v != '#' {
				if prev == -1 {
					other.WriteRune(v)
				}
				continue
			}
			if prev == -1 {
				prev = i
			} else {
				ss = append(ss, strings.TrimSpace(text[prev+1:i]))
				prev = -1
			}
		}
		other.WriteRune('\n')
	}

	tag := strings.Join(ss, "#\n#")
	if tag != "" {
		other.WriteString("\n#")
		other.WriteString(tag)
		other.WriteRune('#')
	}
	return other.String(), reflect.StructTag(strings.Join(ss, " "))
}
