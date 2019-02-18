package utils

import (
	"crypto/md5"
	"encoding/hex"
	"go/build"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func PackageOmitted(pkgpath string) []string {
	dirs := build.Default.SrcDirs()
	pkgpath = strings.TrimSuffix(pkgpath, "/...")

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

func GetName(name string) string {
	i := strings.Index(name, ".")
	if i == -1 {
		return name
	}
	return name[:i]
}

func MergeLine(t string) string {
	return strings.TrimSpace(strings.Replace(t, "\n", " ", -1))
}

func CommentLine(t string) string {
	return "// " + strings.Join(strings.Split(strings.TrimSpace(t), "\n"), "\n// ") + "\n"
}

func Hash(s ...string) string {
	h := md5.New()
	for _, v := range s {
		h.Write([]byte(v))
	}
	return hex.EncodeToString(h.Sum(nil)[:])
}
