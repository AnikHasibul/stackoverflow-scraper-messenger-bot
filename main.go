package main

import "net/http"

func main() {
	defer func() {
		err := http.ListenAndServe(":8080", nil)
		CheckErr("Error on listener:", err)
	}()
	http.HandleFunc("/webhook", botHandler)
}
