package main

import "testing"

func TestCrearPedidoVaciaCarritoYBajaStock(t *testing.T) {
	s := storeDePrueba(t)
	uid := usuarioDePrueba(t, s)
	prod, _ := s.Crear(Producto{Nombre: "Audífonos", Precio: 100, Stock: 5})
	s.AgregarAlCarrito(uid, prod.ID, 2)

	ped, err := s.CrearPedido(uid)
	if err != nil {
		t.Fatalf("error creando pedido: %v", err)
	}
	if ped.Total != 200 {
		t.Fatalf("esperaba total 200, obtuve %v", ped.Total)
	}
	if ped.Resumen != "Audífonos x2" {
		t.Fatalf("resumen incorrecto: %q", ped.Resumen)
	}

	// El carrito debe quedar vacío.
	items, _ := s.VerCarrito(uid)
	if len(items) != 0 {
		t.Fatalf("el carrito debería vaciarse, tiene %d", len(items))
	}

	// El stock debe haber bajado de 5 a 3.
	got, _, _ := s.Obtener(prod.ID)
	if got.Stock != 3 {
		t.Fatalf("esperaba stock 3 tras vender 2, obtuve %d", got.Stock)
	}
}

func TestCrearPedidoCarritoVacio(t *testing.T) {
	s := storeDePrueba(t)
	uid := usuarioDePrueba(t, s)

	_, err := s.CrearPedido(uid)
	if err == nil {
		t.Fatal("comprar con carrito vacío debería dar error")
	}
}

func TestVerPedidos(t *testing.T) {
	s := storeDePrueba(t)
	uid := usuarioDePrueba(t, s)
	prod, _ := s.Crear(Producto{Nombre: "Reloj", Precio: 300, Stock: 10})

	// Dos compras.
	s.AgregarAlCarrito(uid, prod.ID, 1)
	s.CrearPedido(uid)
	s.AgregarAlCarrito(uid, prod.ID, 1)
	s.CrearPedido(uid)

	pedidos, err := s.VerPedidos(uid)
	if err != nil {
		t.Fatalf("error viendo pedidos: %v", err)
	}
	if len(pedidos) != 2 {
		t.Fatalf("esperaba 2 pedidos, obtuve %d", len(pedidos))
	}
}
