package main

import (
	"fmt"
	"net/http"
)

func checkLogin(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isLoggedOut(r) {
			if r.Method != http.MethodPost {
				fmt.Fprint(w, htmlFiles["login"])
				return
			}
			http.Error(w, "Logged Out", 400)
			return
		}
		fn(w, r)
	}
}

func isLoggedOut(r *http.Request) bool {
	return !db.IsReady() || !sessionManager.Check(r)
}
