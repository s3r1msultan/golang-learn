package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
)

var jwtKey = []byte(GetJWTKey())

func GenerateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   email,
		ExpiresAt: expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

func JWTAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwtToken") // "jwtToken" is the name of the cookie
		if err != nil {
			http.Redirect(w, r, "/auth", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Redirect(w, r, "/auth", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetJWTKey() string {
	jwtKey := os.Getenv("JWT_KEY")
	return jwtKey
}
