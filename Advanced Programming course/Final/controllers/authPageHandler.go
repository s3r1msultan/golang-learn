package controllers

import (
	"context"
	"encoding/json"
	"final/db"
	"final/middlewares"
	"final/models"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	var err error
	DBClient, err = db.Connect()
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//err = r.ParseForm()
	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//User := models.User{
	//	Email:     r.FormValue("email"),
	//	Password:  r.FormValue("password"),
	//	FirstName: r.FormValue("first_name"),
	//	LastName:  r.FormValue("last_name"),
	//}

	// Hash the User's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	usersCollection := DBClient.Database("go_restaurants").Collection("users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := GenerateToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.VerificationToken = token
	user.EmailVerified = false
	_, err = usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = SendVerificationEmail(user.Email, token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	DBClient, err = db.Connect()
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var creds Credentials
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	usersCollection := DBClient.Database("go_restaurants").Collection("users")
	err = usersCollection.FindOne(context.TODO(), bson.M{"email": creds.Email}).Decode(&User)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(creds.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenString, err := middlewares.GenerateJWT(creds.Email)
	if err != nil {
		log.Fatal("Problems with token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString, "id": User.ObjectId.Hex()})
}
