package namer

import (
	"path/filepath"
	"strings"
)

type FileExt map[string]struct{}

func NewFileExt(extList []string) FileExt {
	fe := make(FileExt)
	for _, v := range extList {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if v[0] != '.' {
			v = "." + v
		}
		fe[strings.ToLower(v)] = struct{}{}
	}
	return fe
}

func NewFileExtFromString(s, sep string) FileExt {
	return NewFileExt(strings.Split(s, sep))
}

func (fe FileExt) Contains(ext string) bool {
	_, exists := fe[strings.ToLower(ext)]
	return exists
}

func (fe FileExt) Is(filename string) bool {
	return fe.Contains(filepath.Ext(filename))
}
