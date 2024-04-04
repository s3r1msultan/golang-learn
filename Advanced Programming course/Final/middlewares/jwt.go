package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"time"
)

var jwtKey = []byte(GetJWTKey())

type Claims struct {
	ObjectId primitive.ObjectID
	jwt.StandardClaims
}

var claims Claims

func GenerateJWT(email string, ObjectID primitive.ObjectID) (string, error) {
	expirationTime := time.Now().Add(30 * time.Minute)
	claims = Claims{
		ObjectID,
		jwt.StandardClaims{
			Subject:   email,
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func JWTAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		tokenString := cookie.Value

		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ParseObjectIdFromJWT(r *http.Request) (primitive.ObjectID, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("cookie 'jwtToken' not found: %w", err)
	}
	tokenString := cookie.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("error parsing token: %w", err)
	}
	if !token.Valid {
		return primitive.NilObjectID, fmt.Errorf("token is invalid")
	}
	return claims.ObjectId, nil
}

func GetJWTKey() string {
	jwtKey := os.Getenv("JWT_KEY")
	return jwtKey
}
