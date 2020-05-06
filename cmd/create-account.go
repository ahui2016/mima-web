package main

import (
	mimadb "github.com/ahui2016/mima-web/db"
	"github.com/ahui2016/mima-web/util"
)

const password = "abb"

func main() {
	db := mimadb.NewDB(util.DatabaseDefaultDir())
	db.Create(password)
}
