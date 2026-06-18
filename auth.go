package main

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// secret es la clave para firmar los tokens. En producción debe venir de la
// variable de entorno JWT_SECRET (larga y secreta).
func secret() []byte {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return []byte(s)
	}
	return []byte("clave-de-desarrollo-cambiar-en-produccion")
}

// RegistrarUsuario crea un usuario con la contraseña cifrada (bcrypt).
func (s *Store) RegistrarUsuario(email, password string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	res, err := s.db.Exec("INSERT INTO usuarios (email, password_hash) VALUES (?, ?)", email, string(hash))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// AutenticarUsuario comprueba email + contraseña. Devuelve el id si es correcto.
func (s *Store) AutenticarUsuario(email, password string) (int64, bool, error) {
	var id int64
	var hash string
	err := s.db.QueryRow("SELECT id, password_hash FROM usuarios WHERE email = ?", email).Scan(&id, &hash)
	if err != nil {
		return 0, false, nil // no existe (no revelamos detalles)
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return 0, false, nil // contraseña incorrecta
	}
	return id, true, nil
}

// GenerarToken crea un token JWT para el usuario (caduca en 24h).
func GenerarToken(userID int64) (string, error) {
	claims := jwt.MapClaims{"user_id": userID, "exp": time.Now().Add(24 * time.Hour).Unix()}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret())
}

// ParseToken valida un token y devuelve el id del usuario.
func ParseToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}
		return secret(), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("token inválido o caducado")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("token inválido")
	}
	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("token sin user_id")
	}
	return int64(id), nil
}

// protegido envuelve un handler que necesita usuario logueado. Lee el token del
// header "Authorization: Bearer <token>" y le pasa el id del usuario al handler.
func protegido(h func(http.ResponseWriter, *http.Request, int64)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			escribirJSON(w, http.StatusUnauthorized, map[string]string{"error": "necesitas iniciar sesión"})
			return
		}
		id, err := ParseToken(strings.TrimPrefix(auth, "Bearer "))
		if err != nil {
			escribirJSON(w, http.StatusUnauthorized, map[string]string{"error": "token inválido"})
			return
		}
		h(w, r, id)
	}
}
