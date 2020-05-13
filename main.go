package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ahui2016/mima-web/mima"
	"github.com/ahui2016/mima-web/util"
)

func main() {
	fs := http.FileServer(http.Dir("public/"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/index", checkLogin(indexPage))
	http.HandleFunc("/api/all-items", checkLogin(getAllHandler))
	http.HandleFunc("/login", noMore(loginPage))
	http.HandleFunc("/api/login", noMore(loginHandler))
	http.HandleFunc("/add", checkLogin(addPage))
	http.HandleFunc("/api/add", checkLogin(addHandler))
	http.HandleFunc("/api/random-password", checkLogin(randomPassword))
	http.HandleFunc("/edit", checkLogin(editPage))
	http.HandleFunc("/api/edit", checkLogin(editHandler))
	http.HandleFunc("/api/item", checkLogin(getItemHandler))
	http.HandleFunc("/api/delete-history", checkLogin(deleteHistory))
	http.HandleFunc("/api/delete", checkLogin(deleteHandler))
	http.HandleFunc("/recyclebin", checkLogin(recyclePage))
	http.HandleFunc("/api/deleted-items", checkLogin(getDeletedHandler))
	http.HandleFunc("/api/recover", checkLogin(recoverHandler))
	http.HandleFunc("/api/delete-forever", checkLogin(deleteForever))
	http.HandleFunc("/search", checkLogin(searchPage))
	http.HandleFunc("/api/search", checkLogin(searchHandler))

	addr := "0.0.0.0:8080"
	fmt.Println(addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fallthrough
	case "/home":
		redirect(w, r, "/index", 302)
	default:
		http.NotFound(w, r)
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, htmlFiles["index"])
}

func recyclePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, htmlFiles["recyclebin"])
}

func getAllHandler(w http.ResponseWriter, r *http.Request) {
	allItems, err := json.Marshal(db.HomeCache)
	if checkErr(w, err, 500) {
		return
	}
	fmt.Fprint(w, string(allItems))
}

func getDeletedHandler(w http.ResponseWriter, r *http.Request) {
	deleted, err := json.Marshal(db.RecycleCache)
	if checkErr(w, err, 500) {
		return
	}
	fmt.Fprint(w, string(deleted))
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, htmlFiles["login"])
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")

	// If db is not ready, initialize it.
	if !db.IsReady() {
		if err := db.Init(password); err != nil {
			passwordTry++
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

func addHandler(w http.ResponseWriter, r *http.Request) {
	form := &mima.EditForm{
		Title:    strings.TrimSpace(r.FormValue("title")),
		Username: strings.TrimSpace(r.FormValue("username")),
		Password: r.FormValue("password"),
		Notes:    strings.TrimSpace(r.FormValue("notes")),
	}
	m := mima.NewFrom(form)
	if checkErr(w, db.Insert(m), 400) {
		return
	}
	fmt.Fprint(w, m.ID)
}

func randomPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, randomString16())
}

func editPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, htmlFiles["edit"])
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	form := &mima.EditForm{
		ID:       r.FormValue("id"),
		Title:    strings.TrimSpace(r.FormValue("title")),
		Alias:    strings.TrimSpace(r.FormValue("alias")),
		Username: strings.TrimSpace(r.FormValue("username")),
		Password: r.FormValue("password"),
		Notes:    strings.TrimSpace(r.FormValue("notes")),
	}
	checkErr(w, db.Update(form), 400)
}

func getItemHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	_, m, err := db.GetByID(id)
	if checkErr(w, err, 400) {
		return
	}
	mJSON, err := json.Marshal(m)
	if checkErr(w, err, 500) {
		return
	}
	fmt.Fprint(w, string(mJSON))
}

func deleteHistory(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	datetime := r.FormValue("datetime")
	err := db.DeleteHistory(id, datetime)
	checkErr(w, err, 400)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	checkErr(w, db.RecycleByID(r.FormValue("id")), 400)
}

func recoverHandler(w http.ResponseWriter, r *http.Request) {
	checkErr(w, db.RecoverByID(r.FormValue("id")), 400)
}

func deleteForever(w http.ResponseWriter, r *http.Request) {
	checkErr(w, db.DeleteForever(r.FormValue("id")), 400)
}

func searchPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, htmlFiles["search"])
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	alias := strings.TrimSpace(r.FormValue("alias"))
	if alias == "" {
		http.Error(w, "search text is empty", 400)
		return
	}
}
