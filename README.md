# 🛒 Mi Tienda — Backend (Go)

API REST de una tienda online, hecha en **Go** con base de datos **SQLite**.
Este es el **proyecto grande** de Samuel, construido por etapas. 🚀

## ✅ Etapa 1 (esta) — Productos
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

## 🛣️ Próximas etapas
- **Etapa 2:** login de usuarios + carrito de compras.
- **Etapa 3:** la página web (frontend) + desplegar en la nube.

## 🛠️ Tecnología
Go (librería estándar), SQLite, pruebas + CI (GitHub Actions), CORS.

Hecho por **Samuel Catalán** 🇬🇹
