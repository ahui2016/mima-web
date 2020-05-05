package main

import mimadb "github.com/ahui2016/mima-web/db"

const password = "abb"

func main() {
	db := mimadb.NewDB()
	db.Create(password, "")
}
