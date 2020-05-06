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
	http.HandleFunc("/login", noMore(loginPage))
	http.HandleFunc("/api/login", noMore(loginHandler))

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

	// If db is not ready, initialize it.
	if !db.IsReady() {
		if err := db.Init(password); err != nil {
			passwordTry += 1
			if checkPasswordTry(w) {
				return
			}
			http.Error(w, err.Error(), 400)
			return
		}
	}

	// if db is ready but it's session has not set.
	if !db.CheckPassword(password) {
		passwordTry += 1
		if checkPasswordTry(w) {
			return
		}
		http.Error(w, "wrong password", 400)
		return
	}

	// db is ready and the password is correct.
	passwordTry = 0
	sessionManager.Add(w, util.NewID())
}
