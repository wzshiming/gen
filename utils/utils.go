package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

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
	return hex.EncodeToString(h.Sum(nil)[:2])
}
