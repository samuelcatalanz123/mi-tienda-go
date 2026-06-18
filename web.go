package main

import (
	_ "embed"
	"net/http"
)

//go:embed index.html
var paginaWeb []byte

// homeHandler sirve la página web de la tienda.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(paginaWeb)
}
