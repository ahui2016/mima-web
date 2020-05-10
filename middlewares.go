package main

import (
	"fmt"
	"net/http"
)

func checkLogin(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if checkPasswordTry(w) {
			return
		}
		if isLoggedOut(r) {
			if r.Method != http.MethodPost {
				fmt.Fprint(w, htmlFiles["login"])
				return
			}

			// http.MethodPost
			http.Error(w, "Logged Out", 400)
			return
		}
		fn(w, r)
	}
}

// noMore is a middleware only apply to login.
func noMore(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if checkPasswordTry(w) {
			return
		}
		if !isLoggedOut(r) {
			http.Error(w, "You've logged in.", 400)
			return
		}
		fn(w, r)
	}
}

func isLoggedOut(r *http.Request) bool {
	return !db.IsReady() || !sessionManager.Check(r)
}

func checkPasswordTry(w http.ResponseWriter) bool {
	if passwordTry >= passwordMaxTry {
		db.Reset()
		http.Error(w, "No more try. Input wrong password too many times.", 403)
		return true
	}
	return false
}
