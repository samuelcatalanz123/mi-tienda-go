package main

import "testing"

// usuarioDePrueba registra un usuario y devuelve su id (para tests del carrito).
func usuarioDePrueba(t *testing.T, s *Store) int64 {
	id, err := s.RegistrarUsuario("cliente@test.com", "1234")
	if err != nil {
		t.Fatalf("no se pudo registrar usuario: %v", err)
	}
	return id
}

func TestCarritoAgrupaCantidades(t *testing.T) {
	s := storeDePrueba(t)
	uid := usuarioDePrueba(t, s)
	prod, _ := s.Crear(Producto{Nombre: "Camiseta", Precio: 50, Stock: 10})

	// Agrego el mismo producto 3 veces: debe quedar UNA línea con cantidad 3.
	for i := 0; i < 3; i++ {
		if err := s.AgregarAlCarrito(uid, prod.ID, 1); err != nil {
			t.Fatalf("error agregando: %v", err)
		}
	}

	items, _ := s.VerCarrito(uid)
	if len(items) != 1 {
		t.Fatalf("esperaba 1 línea agrupada, obtuve %d", len(items))
	}
	if items[0].Cantidad != 3 {
		t.Fatalf("esperaba cantidad 3, obtuve %d", items[0].Cantidad)
	}
	if items[0].Subtotal != 150 {
		t.Fatalf("esperaba subtotal 150, obtuve %v", items[0].Subtotal)
	}
}

func TestRestarDelCarrito(t *testing.T) {
	s := storeDePrueba(t)
	uid := usuarioDePrueba(t, s)
	prod, _ := s.Crear(Producto{Nombre: "Gorra", Precio: 30, Stock: 10})
	s.AgregarAlCarrito(uid, prod.ID, 2)

	items, _ := s.VerCarrito(uid)
	lineaID := items[0].ID

	// Resto 1 → debe quedar cantidad 1.
	if err := s.RestarDelCarrito(uid, lineaID); err != nil {
		t.Fatalf("error restando: %v", err)
	}
	items, _ = s.VerCarrito(uid)
	if items[0].Cantidad != 1 {
		t.Fatalf("esperaba cantidad 1, obtuve %d", items[0].Cantidad)
	}

	// Resto otra vez → la línea debe desaparecer.
	s.RestarDelCarrito(uid, lineaID)
	items, _ = s.VerCarrito(uid)
	if len(items) != 0 {
		t.Fatalf("el carrito debería estar vacío, tiene %d", len(items))
	}
}

func TestQuitarDelCarritoAjeno(t *testing.T) {
	s := storeDePrueba(t)
	uid := usuarioDePrueba(t, s)
	prod, _ := s.Crear(Producto{Nombre: "Lentes", Precio: 90, Stock: 5})
	s.AgregarAlCarrito(uid, prod.ID, 1)
	items, _ := s.VerCarrito(uid)

	// Otro usuario no debería poder quitar la línea ajena.
	otro, _ := s.RegistrarUsuario("otro@test.com", "1234")
	ok, _ := s.QuitarDelCarrito(otro, items[0].ID)
	if ok {
		t.Fatal("un usuario no debería quitar el carrito de otro")
	}
}
