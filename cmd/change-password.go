package main

import (
	"flag"
	"log"

	mimadb "github.com/ahui2016/mima-web/db"
	"github.com/ahui2016/mima-web/util"
)

var (
	oldPass = flag.String("oldPass", "", "The current password")
	newPass = flag.String("newPass", "", "A new password")
)

func main() {
	flag.Parse()
	if *newPass == "" {
		log.Fatal("the new password is empty")
	}
	db := mimadb.NewDB(util.DatabaseDefaultDir())
	if err := db.Init(*oldPass); err != nil {
		log.Fatal(err)
	}
	if err := db.ChangePassword(*newPass); err != nil {
		log.Fatal(err)
	}
	log.Println("The password has changed successfully.")
}
