package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"final/db"
	"final/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/smtp"
	"os"
)

func GenerateToken() (string, error) {
	bytes := make([]byte, 16) // 128-bit token
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func SendVerificationEmail(to, token string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	auth := smtp.PlainAuth("", from, password, smtpHost)

	verificationURL := "http://localhost:3000/verify?token=" + token
	message := []byte("To: " + to + "\r\n" +
		"Subject: Verify your email address\r\n\r\n" +
		"Click the link to verify your email address: " + verificationURL)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
}

func VerificationHandler(w http.ResponseWriter, r *http.Request) {
	dbClient, err := db.Connect()
	if err != nil {
		http.Error(w, "An error occurred", http.StatusInternalServerError)
	}
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}
	var result models.User
	err = dbClient.Database("go_restaurants").Collection("users").
		FindOneAndUpdate(
			context.TODO(),
			bson.M{"verification_token": token},
			bson.M{"$set": bson.M{"email_verified": true}},
		).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "Invalid token", http.StatusNotFound)
		} else {
			http.Error(w, "An error occurred", http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, "/auth", http.StatusSeeOther)
	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Verification successful"))
}
