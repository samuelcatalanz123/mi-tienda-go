package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Pedido es una compra finalizada.
type Pedido struct {
	ID      int64   `json:"id"`
	Resumen string  `json:"resumen"`
	Total   float64 `json:"total"`
	Fecha   string  `json:"fecha"`
}

// CrearPedido convierte el carrito del usuario en un pedido y vacía el carrito.
func (s *Store) CrearPedido(usuarioID int64) (Pedido, error) {
	items, err := s.VerCarrito(usuarioID)
	if err != nil {
		return Pedido{}, err
	}
	if len(items) == 0 {
		return Pedido{}, errors.New("el carrito está vacío")
	}

	var partes []string
	var total float64
	for _, it := range items {
		partes = append(partes, fmt.Sprintf("%s x%d", it.Nombre, it.Cantidad))
		total += it.Subtotal
		// bajar el stock del producto comprado (nunca por debajo de 0)
		s.db.Exec(
			s.rb("UPDATE productos SET stock = CASE WHEN stock < ? THEN 0 ELSE stock - ? END WHERE id = ?"),
			it.Cantidad, it.Cantidad, it.ProductoID)
	}
	resumen := strings.Join(partes, ", ")
	fecha := time.Now().Format("2006-01-02 15:04")

	id, err := s.insertID(
		"INSERT INTO pedidos (usuario_id, resumen, total, fecha) VALUES (?, ?, ?, ?)",
		usuarioID, resumen, total, fecha)
	if err != nil {
		return Pedido{}, err
	}

	// Vaciar el carrito.
	if _, err := s.db.Exec(s.rb("DELETE FROM carrito WHERE usuario_id = ?"), usuarioID); err != nil {
		return Pedido{}, err
	}
	return Pedido{ID: id, Resumen: resumen, Total: total, Fecha: fecha}, nil
}

// VerPedidos devuelve los pedidos del usuario (más recientes primero).
func (s *Store) VerPedidos(usuarioID int64) ([]Pedido, error) {
	rows, err := s.db.Query(
		s.rb("SELECT id, resumen, total, fecha FROM pedidos WHERE usuario_id = ? ORDER BY id DESC"), usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pedidos := []Pedido{}
	for rows.Next() {
		var p Pedido
		if err := rows.Scan(&p.ID, &p.Resumen, &p.Total, &p.Fecha); err != nil {
			return nil, err
		}
		pedidos = append(pedidos, p)
	}
	return pedidos, rows.Err()
}
