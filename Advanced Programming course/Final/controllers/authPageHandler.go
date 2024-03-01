package controllers

import (
	"context"
	"encoding/json"
	"final/db"
	"final/initializers"
	"final/middlewares"
	"final/models"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Credentials struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var User models.User
var DBClient *mongo.Client

func AuthPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := initTemplates()
	headerData := models.HeaderData{CurrentSite: "Auth"}
	headData := models.HeadData{HeadTitle: "Authorization Page", StyleName: "Auth"}
	data := models.PageData{
		HeaderData: headerData,
		HeadData:   headData,
	}

	err := tmpl.ExecuteTemplate(w, "Auth.html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		initializers.LogError("decoding json while signing up", err, nil)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user already exists
	usersCollection := db.GetUsersCollection()
	var existingUser models.User
	err = usersCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		// User exists
		http.Error(w, "A user with this email already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		initializers.LogError("generating hash password when signing up", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	user.IsAdmin = false

	token, err := GenerateToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.VerificationToken = token
	user.EmailVerified = false

	_, err = usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		initializers.LogError("creating new user while signing up", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = SendVerificationEmail(user.Email, token)
	if err != nil {
		initializers.LogError("sending verification email", err, nil)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		initializers.LogError("decoding json while signing in", err, nil)
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}
	usersCollection := db.GetUsersCollection()
	err = usersCollection.FindOne(context.TODO(), bson.M{"email": creds.Email}).Decode(&User)
	if err != nil {
		initializers.LogError("finding the user", err, nil)
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(creds.Password))
	if err != nil {
		initializers.LogError("comparing hash password", err, nil)
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}
	tokenString, err := middlewares.GenerateJWT(creds.Email, User.ObjectId)
	if err != nil {
		initializers.LogError("comparing jwt token", err, nil)
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}
	response := map[string]interface{}{
		"token":   tokenString,
		"isAdmin": User.IsAdmin,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
