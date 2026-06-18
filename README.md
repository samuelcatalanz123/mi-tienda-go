# 🛒 Mi Tienda — Backend (Go)

API REST de una tienda online, hecha en **Go** con base de datos **SQLite**.
Este es el **proyecto grande** de Samuel, construido por etapas. 🚀

## ✅ Etapa 1 — Productos
CRUD completo de productos: crear, listar, ver, actualizar y borrar.

### Endpoints
| Método | Ruta | Qué hace |
|--------|------|----------|
| GET | `/productos` | lista todos |
| POST | `/productos` | crea uno |
| GET | `/productos/{id}` | trae uno |
| PUT | `/productos/{id}` | actualiza |
| DELETE | `/productos/{id}` | borra |
| GET | `/health` | comprueba que está vivo |

```bash
go run .        # arranca en :8080
go test ./...   # pruebas
```

## ✅ Etapa 2 (esta) — Usuarios y carrito
Registro y login con **JWT** + contraseñas cifradas (**bcrypt**). Carrito por usuario con total.

### Endpoints nuevos
- `POST /registro` · `POST /login` (devuelven token)
- `POST /carrito` · `GET /carrito` · `DELETE /carrito/{id}` (requieren token)

## ✅ Etapa 3 (esta) — Página web + despliegue
Página web (frontend) servida por el propio servidor Go: ver productos, registrarse, iniciar sesión, agregar al carrito y ver el total. Lista para desplegar en Render (incluye render.yaml).

## 🛠️ Tecnología
Go, SQLite, **JWT**, **bcrypt**, pruebas + CI (GitHub Actions), CORS.

Hecho por **Samuel Catalán** 🇬🇹
