package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ahui2016/mima-web/util"
)

func main() {
	fs := http.FileServer(http.Dir("public/"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/index", checkLogin(indexPage))
	http.HandleFunc("/login", noMore(loginPage))
	http.HandleFunc("/api/login", noMore(loginHandler))
	http.HandleFunc("/add", checkLogin(addPage))

	addr := "0.0.0.0:9000"
	fmt.Println(addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fallthrough
	case "/home":
		http.Redirect(w, r, "/index", 302)
	default:
		http.NotFound(w, r)
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
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

func addPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, htmlFiles["add"])
}
