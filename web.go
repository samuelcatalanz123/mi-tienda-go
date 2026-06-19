package main

import (
	_ "embed"
	"net/http"
)

//go:embed index.html
var paginaWeb []byte

//go:embed manifest.json
var manifestJSON []byte

//go:embed sw.js
var swJS []byte

//go:embed icon-192.png
var icon192 []byte

//go:embed icon-512.png
var icon512 []byte

// homeHandler sirve la página web de la tienda.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(paginaWeb)
}

// manifestHandler sirve la "cédula" de la app instalable.
func manifestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/manifest+json")
	_, _ = w.Write(manifestJSON)
}

// swHandler sirve el service worker (modo sin internet).
func swHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	_, _ = w.Write(swJS)
}

// iconoHandler devuelve un handler que sirve un ícono PNG.
func iconoHandler(datos []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		_, _ = w.Write(datos)
	}
}
