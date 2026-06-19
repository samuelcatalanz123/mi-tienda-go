package main

import "testing"

func TestRebindPostgres(t *testing.T) {
	s := &Store{driver: "postgres"}
	got := s.rb("INSERT INTO x (a, b, c) VALUES (?, ?, ?)")
	want := "INSERT INTO x (a, b, c) VALUES ($1, $2, $3)"
	if got != want {
		t.Fatalf("rebind postgres incorrecto:\n got: %s\nwant: %s", got, want)
	}
}

func TestRebindSQLiteNoCambia(t *testing.T) {
	s := &Store{driver: "sqlite"}
	q := "SELECT * FROM x WHERE a = ? AND b = ?"
	if got := s.rb(q); got != q {
		t.Fatalf("en sqlite no debería cambiar la consulta:\n got: %s\nwant: %s", got, q)
	}
}
