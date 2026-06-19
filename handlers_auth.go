package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// POST /registro → crea una cuenta y devuelve un token.
func registroHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil || u.Email == "" || len(u.Password) < 4 {
			escribirJSON(w, 400, map[string]string{"error": "falta correo o contraseña (mín. 4)"})
			return
		}
		id, err := s.RegistrarUsuario(u.Email, u.Password)
		if err != nil {
			escribirJSON(w, 400, map[string]string{"error": "ese correo ya está registrado"})
			return
		}
		token, _ := GenerarToken(id)
		escribirJSON(w, 201, map[string]string{"token": token, "mensaje": "¡cuenta creada!"})
	}
}

// POST /login → comprueba la cuenta y devuelve un token.
func loginHandler(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			escribirJSON(w, 400, map[string]string{"error": "json inválido"})
			return
		}
		id, ok, _ := s.AutenticarUsuario(u.Email, u.Password)
		if !ok {
			escribirJSON(w, 401, map[string]string{"error": "correo o contraseña incorrectos"})
			return
		}
		token, _ := GenerarToken(id)
		escribirJSON(w, 200, map[string]string{"token": token})
	}
}

// POST /carrito → agrega un producto al carrito (requiere login).
func agregarCarritoHandler(s *Store) func(http.ResponseWriter, *http.Request, int64) {
	return func(w http.ResponseWriter, r *http.Request, userID int64) {
		var b struct {
			ProductoID int64 `json:"producto_id"`
			Cantidad   int   `json:"cantidad"`
		}
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil || b.ProductoID == 0 {
			escribirJSON(w, 400, map[string]string{"error": "falta el producto_id"})
			return
		}
		if err := s.AgregarAlCarrito(userID, b.ProductoID, b.Cantidad); err != nil {
			escribirJSON(w, 500, map[string]string{"error": "no se pudo agregar"})
			return
		}
		escribirJSON(w, 201, map[string]string{"mensaje": "producto agregado al carrito"})
	}
}

// GET /carrito → muestra el carrito y el total (requiere login).
func verCarritoHandler(s *Store) func(http.ResponseWriter, *http.Request, int64) {
	return func(w http.ResponseWriter, r *http.Request, userID int64) {
		items, err := s.VerCarrito(userID)
		if err != nil {
			escribirJSON(w, 500, map[string]string{"error": "no se pudo leer el carrito"})
			return
		}
		var total float64
		for _, it := range items {
			total += it.Subtotal
		}
		escribirJSON(w, 200, map[string]any{"items": items, "total": total})
	}
}

// POST /pedidos → finaliza la compra (convierte el carrito en pedido).
func crearPedidoHandler(s *Store) func(http.ResponseWriter, *http.Request, int64) {
	return func(w http.ResponseWriter, r *http.Request, userID int64) {
		ped, err := s.CrearPedido(userID)
		if err != nil {
			escribirJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		escribirJSON(w, 201, map[string]any{"mensaje": "¡Pedido realizado! 🎉", "pedido": ped})
	}
}

// GET /pedidos → lista los pedidos del usuario.
func verPedidosHandler(s *Store) func(http.ResponseWriter, *http.Request, int64) {
	return func(w http.ResponseWriter, r *http.Request, userID int64) {
		pedidos, err := s.VerPedidos(userID)
		if err != nil {
			escribirJSON(w, 500, map[string]string{"error": "no se pudieron leer los pedidos"})
			return
		}
		escribirJSON(w, 200, pedidos)
	}
}

// POST /carrito/{id}/restar → baja en 1 la cantidad (requiere login).
func restarCarritoHandler(s *Store) func(http.ResponseWriter, *http.Request, int64) {
	return func(w http.ResponseWriter, r *http.Request, userID int64) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			escribirJSON(w, 400, map[string]string{"error": "id inválido"})
			return
		}
		if err := s.RestarDelCarrito(userID, id); err != nil {
			escribirJSON(w, 500, map[string]string{"error": "error del servidor"})
			return
		}
		escribirJSON(w, 200, map[string]string{"mensaje": "actualizado"})
	}
}

// DELETE /carrito/{id} → quita una línea del carrito (requiere login).
func quitarCarritoHandler(s *Store) func(http.ResponseWriter, *http.Request, int64) {
	return func(w http.ResponseWriter, r *http.Request, userID int64) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			escribirJSON(w, 400, map[string]string{"error": "id inválido"})
			return
		}
		ok, err := s.QuitarDelCarrito(userID, id)
		if err != nil {
			escribirJSON(w, 500, map[string]string{"error": "error del servidor"})
			return
		}
		if !ok {
			escribirJSON(w, 404, map[string]string{"error": "no está en tu carrito"})
			return
		}
		escribirJSON(w, 200, map[string]string{"mensaje": "quitado del carrito"})
	}
}
