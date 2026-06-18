package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	store, err := NewStore("tienda.db")
	if err != nil {
		log.Fatal("no se pudo abrir la base de datos: ", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /productos", listarHandler(store))
	mux.HandleFunc("POST /productos", crearHandler(store))
	mux.HandleFunc("GET /productos/{id}", obtenerHandler(store))
	mux.HandleFunc("PUT /productos/{id}", actualizarHandler(store))
	mux.HandleFunc("DELETE /productos/{id}", borrarHandler(store))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🛒 Tienda escuchando en :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, conCORS(mux)))
}

// conCORS permite que una página web (el frontend, más adelante) llame a la API.
func conCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func escribirJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	escribirJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// GET /productos → lista todos.
func listarHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productos, err := s.Listar()
		if err != nil {
			escribirJSON(w, 500, map[string]string{"error": "no se pudo listar"})
			return
		}
		escribirJSON(w, 200, productos)
	}
}

// POST /productos → crea un producto.
func crearHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p Producto
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil || p.Nombre == "" {
			escribirJSON(w, 400, map[string]string{"error": "falta el nombre del producto"})
			return
		}
		nuevo, err := s.Crear(p)
		if err != nil {
			escribirJSON(w, 500, map[string]string{"error": "no se pudo crear"})
			return
		}
		escribirJSON(w, 201, nuevo)
	}
}

// GET /productos/{id} → un producto.
func obtenerHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			escribirJSON(w, 400, map[string]string{"error": "id inválido"})
			return
		}
		p, existe, err := s.Obtener(id)
		if err != nil {
			escribirJSON(w, 500, map[string]string{"error": "error del servidor"})
			return
		}
		if !existe {
			escribirJSON(w, 404, map[string]string{"error": "producto no encontrado"})
			return
		}
		escribirJSON(w, 200, p)
	}
}

// PUT /productos/{id} → actualiza.
func actualizarHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			escribirJSON(w, 400, map[string]string{"error": "id inválido"})
			return
		}
		var p Producto
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil || p.Nombre == "" {
			escribirJSON(w, 400, map[string]string{"error": "falta el nombre"})
			return
		}
		existe, err := s.Actualizar(id, p)
		if err != nil {
			escribirJSON(w, 500, map[string]string{"error": "no se pudo actualizar"})
			return
		}
		if !existe {
			escribirJSON(w, 404, map[string]string{"error": "producto no encontrado"})
			return
		}
		escribirJSON(w, 200, map[string]string{"mensaje": "producto actualizado"})
	}
}

// DELETE /productos/{id} → borra.
func borrarHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			escribirJSON(w, 400, map[string]string{"error": "id inválido"})
			return
		}
		existe, err := s.Borrar(id)
		if err != nil {
			escribirJSON(w, 500, map[string]string{"error": "no se pudo borrar"})
			return
		}
		if !existe {
			escribirJSON(w, 404, map[string]string{"error": "producto no encontrado"})
			return
		}
		escribirJSON(w, 200, map[string]string{"mensaje": "producto borrado"})
	}
}
