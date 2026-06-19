package main

import (
	"database/sql"
)

// Producto es un artículo de la tienda.
type Producto struct {
	ID          int64   `json:"id"`
	Nombre      string  `json:"nombre"`
	Precio      float64 `json:"precio"`
	Stock       int     `json:"stock"`
	Descripcion string  `json:"descripcion"`
	Imagen      string  `json:"imagen"`
	Categoria   string  `json:"categoria"`
}

// Store guarda la conexión a la base de datos y qué tipo es ("sqlite" o "postgres").
type Store struct {
	db     *sql.DB
	driver string
}

// NewStore abre (o crea) la base de datos y prepara las tablas. Usa SQLite en
// local y PostgreSQL en Render (según la variable DATABASE_URL).
func NewStore(ruta string) (*Store, error) {
	db, driver, err := abrirDB(ruta)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	for _, t := range tablas(driver) {
		if _, err := db.Exec(t); err != nil {
			return nil, err
		}
	}
	s := &Store{db: db, driver: driver}
	// Migración: agrega la columna "categoria" a tiendas que ya existían sin ella.
	// Si la columna ya está, el error se ignora a propósito.
	s.db.Exec("ALTER TABLE productos ADD COLUMN categoria TEXT NOT NULL DEFAULT 'General'")
	return s, nil
}

// Crear guarda un producto nuevo y devuelve el producto con su id.
func (s *Store) Crear(p Producto) (Producto, error) {
	if p.Categoria == "" {
		p.Categoria = "General"
	}
	id, err := s.insertID(
		"INSERT INTO productos (nombre, precio, stock, descripcion, imagen, categoria) VALUES (?, ?, ?, ?, ?, ?)",
		p.Nombre, p.Precio, p.Stock, p.Descripcion, p.Imagen, p.Categoria)
	if err != nil {
		return Producto{}, err
	}
	p.ID = id
	return p, nil
}

// Listar devuelve todos los productos.
func (s *Store) Listar() ([]Producto, error) {
	rows, err := s.db.Query("SELECT id, nombre, precio, stock, descripcion, imagen, categoria FROM productos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productos := []Producto{}
	for rows.Next() {
		var p Producto
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Precio, &p.Stock, &p.Descripcion, &p.Imagen, &p.Categoria); err != nil {
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
		s.rb("SELECT id, nombre, precio, stock, descripcion, imagen, categoria FROM productos WHERE id = ?"), id).
		Scan(&p.ID, &p.Nombre, &p.Precio, &p.Stock, &p.Descripcion, &p.Imagen, &p.Categoria)
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
	if p.Categoria == "" {
		p.Categoria = "General"
	}
	res, err := s.db.Exec(
		s.rb("UPDATE productos SET nombre=?, precio=?, stock=?, descripcion=?, imagen=?, categoria=? WHERE id=?"),
		p.Nombre, p.Precio, p.Stock, p.Descripcion, p.Imagen, p.Categoria, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// Borrar elimina un producto por id. Devuelve false si no existía.
func (s *Store) Borrar(id int64) (bool, error) {
	res, err := s.db.Exec(s.rb("DELETE FROM productos WHERE id = ?"), id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}
