package main

import (
	"database/sql"
	"os"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL (para internet/Render)
	_ "modernc.org/sqlite"             // SQLite (para tu Mac)
)

// abrirDB decide qué base de datos usar:
//   - Si hay DATABASE_URL (en Render) → PostgreSQL (los datos NO se borran).
//   - Si no → SQLite en un archivo local (cómodo para probar en tu Mac).
//
// Devuelve la conexión y el nombre del driver ("postgres" o "sqlite").
func abrirDB(ruta string) (*sql.DB, string, error) {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		db, err := sql.Open("pgx", url)
		return db, "postgres", err
	}
	db, err := sql.Open("sqlite", ruta)
	return db, "sqlite", err
}

// rb ("rebind") adapta los signos de pregunta de las consultas. SQLite usa "?"
// pero PostgreSQL usa "$1, $2, $3...". Así escribimos las consultas una sola vez.
func (s *Store) rb(query string) string {
	if s.driver != "postgres" {
		return query
	}
	var sb strings.Builder
	n := 0
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			n++
			sb.WriteByte('$')
			sb.WriteString(strconv.Itoa(n))
		} else {
			sb.WriteByte(query[i])
		}
	}
	return sb.String()
}

// insertID ejecuta un INSERT y devuelve el id de la fila nueva. SQLite usa
// LastInsertId(); PostgreSQL necesita "RETURNING id". Esto oculta esa diferencia.
func (s *Store) insertID(query string, args ...any) (int64, error) {
	if s.driver == "postgres" {
		var id int64
		err := s.db.QueryRow(s.rb(query)+" RETURNING id", args...).Scan(&id)
		return id, err
	}
	res, err := s.db.Exec(s.rb(query), args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// tablas devuelve las sentencias CREATE TABLE adecuadas para cada base de datos.
func tablas(driver string) []string {
	if driver == "postgres" {
		return []string{
			`CREATE TABLE IF NOT EXISTS productos (
				id          BIGSERIAL PRIMARY KEY,
				nombre      TEXT NOT NULL,
				precio      DOUBLE PRECISION NOT NULL,
				stock       INTEGER NOT NULL DEFAULT 0,
				descripcion TEXT NOT NULL DEFAULT '',
				imagen      TEXT NOT NULL DEFAULT ''
			)`,
			`CREATE TABLE IF NOT EXISTS pedidos (
				id         BIGSERIAL PRIMARY KEY,
				usuario_id BIGINT NOT NULL,
				resumen    TEXT NOT NULL,
				total      DOUBLE PRECISION NOT NULL,
				fecha      TEXT NOT NULL
			)`,
			`CREATE TABLE IF NOT EXISTS usuarios (
				id            BIGSERIAL PRIMARY KEY,
				email         TEXT NOT NULL UNIQUE,
				password_hash TEXT NOT NULL
			)`,
			`CREATE TABLE IF NOT EXISTS carrito (
				id          BIGSERIAL PRIMARY KEY,
				usuario_id  BIGINT NOT NULL,
				producto_id BIGINT NOT NULL,
				cantidad    INTEGER NOT NULL DEFAULT 1
			)`,
		}
	}
	return []string{
		`CREATE TABLE IF NOT EXISTS productos (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			nombre      TEXT NOT NULL,
			precio      REAL NOT NULL,
			stock       INTEGER NOT NULL DEFAULT 0,
			descripcion TEXT NOT NULL DEFAULT '',
			imagen      TEXT NOT NULL DEFAULT ''
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
}
