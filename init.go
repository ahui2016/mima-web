package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	mimadb "github.com/ahui2016/mima-web/db"
	"github.com/ahui2016/mima-web/mima"
	"github.com/ahui2016/mima-web/session"
	"github.com/ahui2016/mima-web/util"
)

const (
	staticFolder   = "static"
	maxAge         = 60 * 30 // seconds
	passwordMaxTry = 5
	backupFilePath = "hidden/mimaweb.db.bak"
	exportFilePath = "hidden/mimaweb.json"
)

var (
	db             *mimadb.DB
	sessionManager *session.Manager
)

var (
	passwordTry = 0
	htmlFiles   = make(map[string]string)
)

type (
	Mima = mima.Mima
)

func init() {
	dbDir := util.DatabaseDefaultDir()
	db = mimadb.NewDB(dbDir)
	fmt.Println(dbDir)
	sessionManager = session.NewManager(maxAge)
	fillHtmlFiles()
}

func randomString16() string {
	return mimadb.RandomString64()[0:16]
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

func checkErr(w http.ResponseWriter, err error, code int) bool {
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), code)
		return true
	}
	return false
}
