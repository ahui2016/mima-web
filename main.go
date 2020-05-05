package main

func main() {
	fs := http.FileServer(http.Dir("static/"))
}
