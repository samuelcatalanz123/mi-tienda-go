package main

import "testing"

func TestTokenIdaYVuelta(t *testing.T) {
	token, err := GenerarToken(42)
	if err != nil {
		t.Fatalf("error generando: %v", err)
	}
	id, err := ParseToken(token)
	if err != nil {
		t.Fatalf("error leyendo: %v", err)
	}
	if id != 42 {
		t.Fatalf("esperaba 42, obtuve %d", id)
	}
}

func TestTokenManipuladoSeRechaza(t *testing.T) {
	token, _ := GenerarToken(7)
	if _, err := ParseToken(token + "x"); err == nil {
		t.Fatal("un token manipulado debería rechazarse")
	}
}

func TestRegistroYLogin(t *testing.T) {
	s := storeDePrueba(t)
	if _, err := s.RegistrarUsuario("samuel@correo.com", "secreta"); err != nil {
		t.Fatalf("error registrando: %v", err)
	}
	// Contraseña correcta → ok.
	if _, ok, _ := s.AutenticarUsuario("samuel@correo.com", "secreta"); !ok {
		t.Fatal("la contraseña correcta debería autenticar")
	}
	// Contraseña incorrecta → no.
	if _, ok, _ := s.AutenticarUsuario("samuel@correo.com", "mala"); ok {
		t.Fatal("una contraseña incorrecta NO debería autenticar")
	}
}

func TestCarrito(t *testing.T) {
	s := storeDePrueba(t)
	uid, _ := s.RegistrarUsuario("ana@correo.com", "1234")
	p, _ := s.Crear(Producto{Nombre: "Camiseta", Precio: 100})

	if err := s.AgregarAlCarrito(uid, p.ID, 2); err != nil {
		t.Fatalf("error agregando: %v", err)
	}
	items, _ := s.VerCarrito(uid)
	if len(items) != 1 {
		t.Fatalf("esperaba 1 item, obtuve %d", len(items))
	}
	if items[0].Nombre != "Camiseta" || items[0].Subtotal != 200 {
		t.Fatalf("subtotal mal calculado: %+v", items[0])
	}

	ok, _ := s.QuitarDelCarrito(uid, items[0].ID)
	if !ok {
		t.Fatal("debería haber quitado el item")
	}
	items, _ = s.VerCarrito(uid)
	if len(items) != 0 {
		t.Fatal("el carrito debería estar vacío")
	}
}
