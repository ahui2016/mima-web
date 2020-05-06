package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ahui2016/mima-web/util"
)

func main() {
	http.HandleFunc("/", checkLogin(homePage))
	http.HandleFunc("/home", checkLogin(homePage))
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/api/login", loginHandler)

	addr := "0.0.0.0:9000"
	fmt.Println(addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if db.IsEmpty() {
		http.Redirect(w, r, "/login", 303)
		return
	}
	allItems, err := json.Marshal(db.AllItems())
	if checkErr(w, err, 500) {
		return
	}
	fmt.Fprint(w, string(allItems))
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, htmlFiles["login"])
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	if db.IsReady() {
	}
	if checkErr(w, db.Init(password), 400) {
		return
	}
	sessionManager.Add(w, util.NewID())
}
