package main

import "database/sql"

// ItemCarrito es una línea del carrito (un producto con su cantidad).
type ItemCarrito struct {
	ID         int64   `json:"id"`
	ProductoID int64   `json:"producto_id"`
	Nombre     string  `json:"nombre"`
	Precio     float64 `json:"precio"`
	Cantidad   int     `json:"cantidad"`
	Subtotal   float64 `json:"subtotal"`
}

// AgregarAlCarrito mete un producto al carrito del usuario. Si ese producto ya
// estaba en el carrito, le suma la cantidad (no crea una línea repetida).
func (s *Store) AgregarAlCarrito(usuarioID, productoID int64, cantidad int) error {
	if cantidad < 1 {
		cantidad = 1
	}
	var id int64
	err := s.db.QueryRow(
		s.rb("SELECT id FROM carrito WHERE usuario_id = ? AND producto_id = ?"),
		usuarioID, productoID).Scan(&id)
	if err == nil {
		// ya estaba: sumamos la cantidad
		_, err = s.db.Exec(s.rb("UPDATE carrito SET cantidad = cantidad + ? WHERE id = ?"), cantidad, id)
		return err
	}
	if err != sql.ErrNoRows {
		return err
	}
	_, err = s.db.Exec(
		s.rb("INSERT INTO carrito (usuario_id, producto_id, cantidad) VALUES (?, ?, ?)"),
		usuarioID, productoID, cantidad)
	return err
}

// RestarDelCarrito baja en 1 la cantidad de una línea. Si llega a 0, la quita.
func (s *Store) RestarDelCarrito(usuarioID, itemID int64) error {
	var cant int
	err := s.db.QueryRow(
		s.rb("SELECT cantidad FROM carrito WHERE id = ? AND usuario_id = ?"),
		itemID, usuarioID).Scan(&cant)
	if err != nil {
		return err
	}
	if cant <= 1 {
		_, err = s.db.Exec(s.rb("DELETE FROM carrito WHERE id = ? AND usuario_id = ?"), itemID, usuarioID)
		return err
	}
	_, err = s.db.Exec(s.rb("UPDATE carrito SET cantidad = cantidad - 1 WHERE id = ? AND usuario_id = ?"), itemID, usuarioID)
	return err
}

// VerCarrito devuelve los productos del carrito del usuario (con su nombre y precio).
func (s *Store) VerCarrito(usuarioID int64) ([]ItemCarrito, error) {
	rows, err := s.db.Query(s.rb(`
		SELECT c.id, p.id, p.nombre, p.precio, c.cantidad
		FROM carrito c
		JOIN productos p ON p.id = c.producto_id
		WHERE c.usuario_id = ?
		ORDER BY c.id`), usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []ItemCarrito{}
	for rows.Next() {
		var it ItemCarrito
		if err := rows.Scan(&it.ID, &it.ProductoID, &it.Nombre, &it.Precio, &it.Cantidad); err != nil {
			return nil, err
		}
		it.Subtotal = it.Precio * float64(it.Cantidad)
		items = append(items, it)
	}
	return items, rows.Err()
}

// QuitarDelCarrito borra una línea del carrito (solo si es del usuario).
func (s *Store) QuitarDelCarrito(usuarioID, itemID int64) (bool, error) {
	res, err := s.db.Exec(s.rb("DELETE FROM carrito WHERE id = ? AND usuario_id = ?"), itemID, usuarioID)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}
