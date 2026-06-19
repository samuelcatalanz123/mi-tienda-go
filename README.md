<div align="center">

<img src="icon-192.png" width="96" alt="Mi Tienda" />

# 🛒 Mi Tienda

**Tienda online full-stack hecha en Go** — API REST + frontend + app instalable (PWA).

[![En vivo](https://img.shields.io/badge/🌐_Demo_en_vivo-mi--tienda--go.onrender.com-2563eb)](https://mi-tienda-go.onrender.com)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Base de datos](https://img.shields.io/badge/DB-SQLite_%2B_PostgreSQL-336791?logo=postgresql&logoColor=white)](#-base-de-datos-doble-sqlite--postgresql)
[![PWA](https://img.shields.io/badge/PWA-instalable-7c3aed)](#)

</div>

---

## 📖 ¿Qué es esto?

Una **tienda online completa** construida desde cero en **Go**, sin frameworks pesados: solo
la librería estándar (`net/http`) para la API y **HTML/CSS/JavaScript puro** para la web.
El mismo servidor de Go sirve la API **y** la página, todo en un solo binario.

Proyecto construido **por etapas**, desde un CRUD básico hasta una app instalable con
base de datos en la nube. 🚀

🌐 **Pruébala en vivo:** https://mi-tienda-go.onrender.com

<div align="center">

<img src="capturas/tienda.png" width="640" alt="Mi Tienda — vista de escritorio" />

<br/>

<img src="capturas/tienda-movil.png" width="260" alt="Mi Tienda — vista móvil (responsive)" />

</div>

---

## ✨ Funciones

**🛍️ Para el cliente**
- Catálogo de productos con imágenes, precios y stock
- Buscador en vivo y ordenar por precio / nombre
- Carrito de compras deslizante con cantidades (➕ / ➖)
- Finalizar compra → genera un pedido
- Historial de "Mis pedidos" + 📲 enviar el pedido por **WhatsApp**
- 🌙 Modo claro/oscuro (recuerda tu preferencia)
- 📱 **Instalable** como app (PWA) y funciona **sin internet**

**🔐 Para el dueño (admin)**
- Registro / login con **JWT** y contraseñas cifradas con **bcrypt**
- CRUD de productos protegido: crear, editar y borrar (solo con sesión)
- 📉 El stock baja solo al vender; los agotados se bloquean
- 📊 Panel con métricas: productos, inventario, valor, pedidos y total vendido

---

## 🛠️ Tecnología

| Área | Herramientas |
|------|--------------|
| **Backend** | Go (librería estándar `net/http`, enrutado por método+ruta de Go 1.22) |
| **Base de datos** | SQLite (`modernc.org/sqlite`, sin CGo) y PostgreSQL (`pgx`) |
| **Auth** | JWT (`golang-jwt`) + bcrypt (`golang.org/x/crypto`) |
| **Frontend** | HTML, CSS y JavaScript puro (sin frameworks) |
| **PWA** | Web App Manifest + Service Worker (offline) |
| **Despliegue** | Render (`render.yaml`) · CI con GitHub Actions |

---

## 🗄️ Base de datos doble (SQLite ↔ PostgreSQL)

El mismo código funciona con **dos bases de datos** y elige sola según el entorno:

- **En local** → SQLite en un archivo (cero instalación, ideal para desarrollar).
- **En producción** → **PostgreSQL** si existe la variable `DATABASE_URL` (datos persistentes).

Una capa fina traduce las diferencias (los `?` ↔ `$1, $2…` y el `RETURNING id`),
así las consultas se escriben **una sola vez**. Ver [`db.go`](db.go).

---

## 🧭 Estructura del proyecto

```
main.go             # rutas y arranque del servidor
db.go               # elige SQLite/PostgreSQL + traductor de consultas
store.go            # productos (CRUD)
carrito.go          # carrito por usuario (agrupa cantidades)
pedidos.go          # pedidos / finalizar compra + baja de stock
auth.go             # JWT, bcrypt y middleware de protección
handlers_auth.go    # handlers de registro, login, carrito y pedidos
web.go              # sirve la web (index.html) y archivos PWA
index.html          # frontend completo (una sola página)
manifest.json·sw.js # archivos de la app instalable (PWA)
*_test.go           # pruebas automáticas
```

---

## 🔌 API

| Método | Ruta | Protegida | Qué hace |
|--------|------|:---------:|----------|
| `GET` | `/productos` | | Lista productos |
| `GET` | `/productos/{id}` | | Trae un producto |
| `POST` | `/productos` | 🔒 | Crea un producto |
| `PUT` | `/productos/{id}` | 🔒 | Actualiza un producto |
| `DELETE` | `/productos/{id}` | 🔒 | Borra un producto |
| `POST` | `/registro` | | Crea cuenta → devuelve token |
| `POST` | `/login` | | Inicia sesión → devuelve token |
| `POST` | `/carrito` | 🔒 | Agrega al carrito (agrupa cantidad) |
| `GET` | `/carrito` | 🔒 | Ve el carrito y el total |
| `POST` | `/carrito/{id}/restar` | 🔒 | Baja 1 a una línea |
| `DELETE` | `/carrito/{id}` | 🔒 | Quita una línea |
| `POST` | `/pedidos` | 🔒 | Finaliza la compra |
| `GET` | `/pedidos` | 🔒 | Historial de pedidos |

🔒 = requiere `Authorization: Bearer <token>`

---

## ▶️ Cómo correrlo

```bash
# 1) Clonar
git clone https://github.com/samuelcatalanz123/mi-tienda-go.git
cd mi-tienda-go

# 2) Correr (usa SQLite automáticamente)
go run .            # abre http://localhost:8080

# 3) Pruebas
go test ./...
```

Para usar PostgreSQL, define la variable `DATABASE_URL`:

```bash
DATABASE_URL="postgres://usuario:clave@host:5432/basededatos" go run .
```

---

## 💡 Puntos destacados (para entrevista)

- **Sin frameworks:** API y frontend sobre la librería estándar de Go → demuestra entender lo básico a fondo.
- **Una base de datos, dos motores:** abstracción limpia entre SQLite y PostgreSQL sin duplicar consultas.
- **Seguridad real:** JWT + bcrypt, rutas de escritura protegidas por middleware.
- **Un solo binario:** el servidor de Go embebe (`//go:embed`) el frontend y los archivos PWA.
- **Pensado para el mundo real:** pedidos por WhatsApp, app instalable, modo offline.

---

<div align="center">

Hecho con 💙 por **Samuel Catalán** · Guatemala 🇬🇹

</div>
