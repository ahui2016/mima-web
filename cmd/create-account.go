package main

import (
	"flag"
	"log"

	mimadb "github.com/ahui2016/mima-web/db"
	"github.com/ahui2016/mima-web/util"
)

var password = flag.String("password", "", "The master password, required.")

func main() {
	flag.Parse()
	if *password == "" {
		log.Fatal("password is empty")
	}
	dbPath := util.DatabaseDefaultDir()
	db := mimadb.NewDB(dbPath)
	db.Create(*password)
	log.Printf("Database created: %s", dbPath)
}
