package main

import (
	"net/http"

	"github.com/KengoWada/playmo"
)

func main() {
	http.HandleFunc("/", playmo.Spotify)
	http.ListenAndServe(":8000", nil)
}
