package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Producto es un artículo de la tienda.
type Producto struct {
	ID          int64   `json:"id"`
	Nombre      string  `json:"nombre"`
	Precio      float64 `json:"precio"`
	Stock       int     `json:"stock"`
	Descripcion string  `json:"descripcion"`
	Imagen      string  `json:"imagen"`
}

// Store guarda la conexión a la base de datos.
type Store struct {
	db *sql.DB
}

// NewStore abre (o crea) la base de datos y prepara la tabla de productos.
func NewStore(ruta string) (*Store, error) {
	db, err := sql.Open("sqlite", ruta)
	if err != nil {
		return nil, err
	}
	tablas := []string{
		`CREATE TABLE IF NOT EXISTS productos (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			nombre      TEXT NOT NULL,
			precio      REAL NOT NULL,
			stock       INTEGER NOT NULL DEFAULT 0,
			descripcion TEXT,
			imagen      TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS pedidos (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			usuario_id INTEGER NOT NULL,
			resumen    TEXT NOT NULL,
			total      REAL NOT NULL,
			fecha      TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS usuarios (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			email         TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS carrito (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			usuario_id  INTEGER NOT NULL,
			producto_id INTEGER NOT NULL,
			cantidad    INTEGER NOT NULL DEFAULT 1
		)`,
	}
	for _, t := range tablas {
		if _, err := db.Exec(t); err != nil {
			return nil, err
		}
	}
	return &Store{db: db}, nil
}

// Crear guarda un producto nuevo y devuelve el producto con su id.
func (s *Store) Crear(p Producto) (Producto, error) {
	res, err := s.db.Exec(
		"INSERT INTO productos (nombre, precio, stock, descripcion, imagen) VALUES (?, ?, ?, ?, ?)",
		p.Nombre, p.Precio, p.Stock, p.Descripcion, p.Imagen)
	if err != nil {
		return Producto{}, err
	}
	p.ID, _ = res.LastInsertId()
	return p, nil
}

// Listar devuelve todos los productos.
func (s *Store) Listar() ([]Producto, error) {
	rows, err := s.db.Query("SELECT id, nombre, precio, stock, descripcion, imagen FROM productos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productos := []Producto{}
	for rows.Next() {
		var p Producto
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Precio, &p.Stock, &p.Descripcion, &p.Imagen); err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}
	return productos, rows.Err()
}

// Obtener devuelve un producto por id (y false si no existe).
func (s *Store) Obtener(id int64) (Producto, bool, error) {
	var p Producto
	err := s.db.QueryRow(
		"SELECT id, nombre, precio, stock, descripcion, imagen FROM productos WHERE id = ?", id).
		Scan(&p.ID, &p.Nombre, &p.Precio, &p.Stock, &p.Descripcion, &p.Imagen)
	if err == sql.ErrNoRows {
		return Producto{}, false, nil
	}
	if err != nil {
		return Producto{}, false, err
	}
	return p, true, nil
}

// Actualizar cambia un producto existente. Devuelve false si no existía.
func (s *Store) Actualizar(id int64, p Producto) (bool, error) {
	res, err := s.db.Exec(
		"UPDATE productos SET nombre=?, precio=?, stock=?, descripcion=?, imagen=? WHERE id=?",
		p.Nombre, p.Precio, p.Stock, p.Descripcion, p.Imagen, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// Borrar elimina un producto por id. Devuelve false si no existía.
func (s *Store) Borrar(id int64) (bool, error) {
	res, err := s.db.Exec("DELETE FROM productos WHERE id = ?", id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}
