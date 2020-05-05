package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/ahui2016/mima-web/util"
)

const (
	staticFolder = "static"
)

var htmlFiles = make(map[string]string)

func init() {
	fillHtmlFiles()
}

func fillHtmlFiles() {
	filePaths, err := util.GetPathsByExt(staticFolder, ".html")
	if err != nil {
		panic(err)
	}
	for _, path := range filePaths {
		base := filepath.Base(path)
		name := strings.TrimSuffix(base, ".html")
		html, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		htmlFiles[name] = string(html)
	}
}
