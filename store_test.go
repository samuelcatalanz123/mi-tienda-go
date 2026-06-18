package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func storeDePrueba(t *testing.T) *Store {
	s, err := NewStore(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("no se pudo crear el store: %v", err)
	}
	return s
}

func TestCrearYListar(t *testing.T) {
	s := storeDePrueba(t)
	s.Crear(Producto{Nombre: "Camiseta", Precio: 50, Stock: 10})
	s.Crear(Producto{Nombre: "Gorra", Precio: 30, Stock: 5})

	ps, err := s.Listar()
	if err != nil {
		t.Fatalf("error listando: %v", err)
	}
	if len(ps) != 2 {
		t.Fatalf("esperaba 2 productos, obtuve %d", len(ps))
	}
	if ps[0].Nombre != "Camiseta" || ps[0].Precio != 50 {
		t.Fatalf("producto incorrecto: %+v", ps[0])
	}
}

func TestObtenerYBorrar(t *testing.T) {
	s := storeDePrueba(t)
	p, _ := s.Crear(Producto{Nombre: "Gorra", Precio: 30})

	got, existe, _ := s.Obtener(p.ID)
	if !existe || got.Nombre != "Gorra" {
		t.Fatalf("no se obtuvo el producto: %+v", got)
	}

	ok, _ := s.Borrar(p.ID)
	if !ok {
		t.Fatal("debería haber borrado")
	}
	_, existe, _ = s.Obtener(p.ID)
	if existe {
		t.Fatal("el producto debería estar borrado")
	}
}

func TestActualizar(t *testing.T) {
	s := storeDePrueba(t)
	p, _ := s.Crear(Producto{Nombre: "Gorra", Precio: 30})

	ok, _ := s.Actualizar(p.ID, Producto{Nombre: "Gorra Premium", Precio: 45, Stock: 3})
	if !ok {
		t.Fatal("debería haber actualizado")
	}
	got, _, _ := s.Obtener(p.ID)
	if got.Nombre != "Gorra Premium" || got.Precio != 45 {
		t.Fatalf("no se actualizó bien: %+v", got)
	}
}

func TestCrearHandler(t *testing.T) {
	s := storeDePrueba(t)
	body := `{"nombre":"Zapatos","precio":120.5,"stock":8}`
	req := httptest.NewRequest(http.MethodPost, "/productos", strings.NewReader(body))
	rec := httptest.NewRecorder()

	crearHandler(s)(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperaba 201, obtuve %d (%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Zapatos") {
		t.Fatalf("la respuesta no tiene el producto: %s", rec.Body.String())
	}
}

func TestCrearHandlerSinNombre(t *testing.T) {
	s := storeDePrueba(t)
	req := httptest.NewRequest(http.MethodPost, "/productos", strings.NewReader(`{"precio":10}`))
	rec := httptest.NewRecorder()

	crearHandler(s)(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400 sin nombre, obtuve %d", rec.Code)
	}
}
