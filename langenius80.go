package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	e := http.ListenAndServe(":80", nil)
	if e != nil {
		fmt.Println(e)
	}
}
